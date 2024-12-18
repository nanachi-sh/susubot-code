package room

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"

	"github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/db"
	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/protos/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/twoonone/player"
	"github.com/twmb/murmur3"
)

type Room struct {
	hash           string
	players        []*player.Player
	operatorNow    *player.Player
	landowner      *player.Player
	farmers        [2]*player.Player
	stage          twoonone_pb.RoomStage
	basicCoin      float64
	multiple       int
	cardPool       []*SendCard
	landownerCards [3]twoonone_pb.Card
}

type SendCard struct {
	SenderInfo        *player.Player
	SendCards         []twoonone_pb.Card
	SendCardType      twoonone_pb.CardType
	SendCardSize      int
	SendCardContinous int
}

func hash() string {
	buf := make([]byte, 100)
	for n := 0; n != len(buf); n++ {
		buf[n] = byte(rand.Intn(256))
	}
	h1, h2 := murmur3.SeedSum128(rand.Uint64(), rand.Uint64(), buf)
	return fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}
func New(basicCoin float64, multiple int) *Room {
	if basicCoin <= 0 {
		basicCoin = 200
	}
	if multiple <= 0 {
		multiple = 1
	}
	return &Room{
		hash:           hash(),
		players:        []*player.Player{},
		operatorNow:    nil,
		landowner:      nil,
		farmers:        [2]*player.Player{},
		stage:          twoonone_pb.RoomStage_WaitingStart,
		basicCoin:      basicCoin,
		multiple:       multiple,
		cardPool:       []*SendCard{},
		landownerCards: [3]twoonone_pb.Card{},
	}
}

func (r *Room) Join(ai *twoonone_pb.PlayerAccountInfo) *twoonone_pb.Errors {
	if r.GetStage() != twoonone_pb.RoomStage_WaitingStart {
		return twoonone_pb.Errors_RoomStarted.Enum()
	}
	if _, ok := r.GetPlayer(ai.Id); ok {
		return twoonone_pb.Errors_RoomExistPlayer.Enum()
	}
	if ai.Coin < r.basicCoin {
		return twoonone_pb.Errors_PlayerCoinLTRoomBasicCoin.Enum()
	}
	p := player.New(&twoonone_pb.PlayerInfo{
		AccountInfo: ai,
		TableInfo: &twoonone_pb.PlayerTableInfo{
			RoomHash:           r.GetHash(),
			Cards:              []twoonone_pb.Card{},
			RobLandownerAction: nil,
		},
	})
	if p == nil {
		return twoonone_pb.Errors_Unexpected.Enum()
	}
	r.players = append(r.players, p)
	return nil
}

func (r *Room) Exit(playerId string) *twoonone_pb.Errors {
	if r.GetStage() != twoonone_pb.RoomStage_WaitingStart {
		return twoonone_pb.Errors_RoomStarted.Enum()
	}
	if _, ok := r.GetPlayer(playerId); !ok {
		return twoonone_pb.Errors_PlayerNoExist.Enum()
	}
	if !r.delete(playerId) {
		return twoonone_pb.Errors_PlayerNoExist.Enum()
	}
	return nil
}

func (r *Room) Start() (*player.Player, *twoonone_pb.Errors) {
	if r.GetStage() != twoonone_pb.RoomStage_WaitingStart {
		return nil, twoonone_pb.Errors_RoomStarted.Enum()
	}
	if len(r.players) != 3 {
		return nil, twoonone_pb.Errors_RoomPlayerNoFull.Enum()
	}
	// 开一副牌
	cardsArr := r.generateCards()
	cards := cardsArr[:]
	r.resequence(&cards)
	// 为玩家发牌
	for n, N := 0, 0; N != len(r.players); N++ {
		playerCards := cards[n : n+17]
		n += 17
		r.resequence(&playerCards)
		r.players[N].AddCards(playerCards)
	}
	// 设置地主牌
	r.setLandownerCards([3]twoonone_pb.Card{
		cards[len(cards)-1],
		cards[len(cards)-2],
		cards[len(cards)-3],
	})
	// 进入下一阶段
	r.startRobLandowner()
	return r.GetOperatorNow(), nil
}

func (r *Room) startRobLandowner() {
	r.stage = twoonone_pb.RoomStage_RobLandownering
	r.operatorNow = r.players[rand.Intn(len(r.players))]
}

func (r *Room) setLandownerCards(cards [3]twoonone_pb.Card) {
	copy(r.landownerCards[:], cards[:])
}

func (r *Room) RobLandownerAction(p *player.Player, action twoonone_pb.RobLandownerActions) (*player.Player, *twoonone_pb.Errors) {
	if p.GetRoomHash() != r.GetHash() {
		return nil, twoonone_pb.Errors_Unexpected.Enum()
	}
	if r.GetStage() != twoonone_pb.RoomStage_RobLandownering {
		return nil, twoonone_pb.Errors_RoomNoRobLandownering.Enum()
	}
	if r.GetOperatorNow().GetId() != p.GetId() {
		return nil, twoonone_pb.Errors_PlayerNoOperatorNow.Enum()
	}
	p.SetRobLandownerAction(&action)
	// 尝试找到下一个未参与玩家
	next := r.nextRobLandownerOperator()
	if next != nil {
		r.operatorNow = next
		return next, nil
	}
	// 全部已参与
	robs := r.getRobLandowners()
	if len(robs) == 0 { //无人抢地主
		return nil, twoonone_pb.Errors_RoomNoRobLandownering.Enum()
	} else if len(robs) == 1 { //仅一人抢地主
		r.landowner = robs[0]
	} else { //一人以上抢地主
		// 按操作时间升序
		sort.Slice(robs, func(i, j int) bool {
			return robs[i].GetRobLandownerActionTime().UnixNano() < robs[j].GetRobLandownerActionTime().UnixNano()
		})
		r.landowner = robs[0]
	}
	r.startSendCard()
	return r.GetOperatorNow(), nil
}

type SendCardEvents struct {
	SenderCardNumberNotice bool
	GameFinish             bool
	SenderCardTypeNotice   bool

	GameFinishE *twoonone_pb.SendCardResponse_GameFinishEvent
	CardType    *twoonone_pb.CardType
	CardNumber  *int
}

func (r *Room) SendCardAction(p *player.Player, sendcards []twoonone_pb.Card, action twoonone_pb.SendCardActions) (*player.Player, SendCardEvents, *twoonone_pb.Errors, error) {
	if p.GetRoomHash() != r.GetHash() {
		return nil, SendCardEvents{}, twoonone_pb.Errors_Unexpected.Enum(), nil
	}
	if r.GetStage() != twoonone_pb.RoomStage_SendingCards {
		return nil, SendCardEvents{}, twoonone_pb.Errors_RoomNoSendingCards.Enum(), nil
	}
	if r.GetOperatorNow().GetId() != p.GetId() {
		return nil, SendCardEvents{}, twoonone_pb.Errors_PlayerNoOperatorNow.Enum(), nil
	}
	var next *player.Player
	event := new(SendCardEvents)
	switch action {
	case twoonone_pb.SendCardActions_Send:
		if len(sendcards) == 0 {
			return nil, SendCardEvents{}, twoonone_pb.Errors_Unexpected.Enum(), nil
		}
		x, y, err := r.sendCard(p, sendcards)
		if err != nil {
			return nil, SendCardEvents{}, err, nil
		}
		next = x
		if y != nil {
			*event = *y
		}
	case twoonone_pb.SendCardActions_NoSend:
		x, err := r.unSendCard(p)
		if err != nil {
			return nil, SendCardEvents{}, err, nil
		}
		next = x
	default:
		return nil, SendCardEvents{}, twoonone_pb.Errors_Unexpected.Enum(), nil
	}
	// 判断出牌玩家手牌剩余数
	switch l := len(p.GetCards()); l {
	case 0:
		event.GameFinish = true
		event.GameFinishE = r.gameFinish()
		if err := r.playerUpdateToDatabase(); err != nil {
			return nil, SendCardEvents{}, nil, err
		}
	case 1, 2:
		event.SenderCardNumberNotice = true
		event.CardNumber = &l
	}
	return next, *event, nil, nil
}

func (r *Room) sendCard(p *player.Player, sendcards []twoonone_pb.Card) (*player.Player, *SendCardEvents, *twoonone_pb.Errors) {
	lastcard := r.GetLastCard()
	// 正常出牌
	cardtype := r.matchCardType(sendcards)
	cardcontious := r.calcCardContinous(sendcards, cardtype)
	cardsize := r.calcCardSize(sendcards, cardtype)
	// 特殊情况判断
	if cardtype == twoonone_pb.CardType_Unknown { //未知牌型
		return nil, nil, twoonone_pb.Errors_SendCardUnknown.Enum()
	} else if lastcard == nil { //本局第一次出牌
		if err := r.playerSendCard(p, sendcards, cardtype, cardsize, cardcontious); err != nil {
			return nil, nil, err
		}
		next := r.nextSendCardOperator()
		r.operatorNow = next
		switch cardtype {
		case twoonone_pb.CardType_KingBomb, twoonone_pb.CardType_Bomb:
			return next, &SendCardEvents{
				SenderCardTypeNotice: true,
				CardType:             &cardtype,
			}, nil
		}
		return next, nil, nil
	} else if lastcard.SenderInfo.GetId() == p.GetId() { //同一人操作
		if err := r.playerSendCard(p, sendcards, cardtype, cardsize, cardcontious); err != nil {
			return nil, nil, err
		}
		next := r.nextSendCardOperator()
		r.operatorNow = next
		switch cardtype {
		case twoonone_pb.CardType_KingBomb, twoonone_pb.CardType_Bomb:
			return next, &SendCardEvents{
				SenderCardTypeNotice: true,
				CardType:             &cardtype,
			}, nil
		}
		return next, nil, nil
	} else if cardtype == twoonone_pb.CardType_KingBomb { // 王炸
		if err := r.playerSendCard(p, sendcards, cardtype, cardsize, cardcontious); err != nil {
			return nil, nil, err
		}
		next := r.nextSendCardOperator()
		r.operatorNow = next
		return next, &SendCardEvents{
			SenderCardTypeNotice: true,
			CardType:             &cardtype,
		}, nil
	} else if cardtype == twoonone_pb.CardType_Bomb && lastcard.SendCardType == twoonone_pb.CardType_Bomb { //上一副与当前都为炸弹
		if cardsize <= lastcard.SendCardSize {
			return nil, nil, twoonone_pb.Errors_SendCardSizeLELastCard.Enum()
		}
		if err := r.playerSendCard(p, sendcards, cardtype, cardsize, cardcontious); err != nil {
			return nil, nil, err
		}
		next := r.nextSendCardOperator()
		r.operatorNow = next
		return next, &SendCardEvents{
			SenderCardTypeNotice: true,
			CardType:             &cardtype,
		}, nil
	} else if cardtype == twoonone_pb.CardType_Bomb && lastcard.SendCardType != twoonone_pb.CardType_KingBomb { //上一副不为王炸且也不为炸弹
		if err := r.playerSendCard(p, sendcards, cardtype, cardsize, cardcontious); err != nil {
			return nil, nil, err
		}
		next := r.nextSendCardOperator()
		r.operatorNow = next
		return next, &SendCardEvents{
			SenderCardTypeNotice: true,
			CardType:             &cardtype,
		}, nil
	}
	//完全正常出牌
	if cardtype != lastcard.SendCardType {
		return nil, nil, twoonone_pb.Errors_SendCardTypeNELastCard.Enum()
	}
	if cardcontious != lastcard.SendCardContinous {
		return nil, nil, twoonone_pb.Errors_SendCardContinousNELastCard.Enum()
	}
	if cardsize <= lastcard.SendCardSize {
		return nil, nil, twoonone_pb.Errors_SendCardSizeLELastCard.Enum()
	}
	if err := r.playerSendCard(p, sendcards, cardtype, cardsize, cardcontious); err != nil {
		return nil, nil, err
	}
	next := r.nextSendCardOperator()
	r.operatorNow = next
	return next, nil, nil
}

func (r *Room) unSendCard(p *player.Player) (*player.Player, *twoonone_pb.Errors) {
	lastcard := r.GetLastCard()
	// 特殊情况判断
	if lastcard == nil { //第一次出牌
		return nil, twoonone_pb.Errors_PlayerIsOnlySendCarder.Enum()
	} else if lastcard.SenderInfo.GetId() == p.GetId() { //上一次出牌为同一人
		return nil, twoonone_pb.Errors_PlayerIsOnlySendCarder.Enum()
	}
	next := r.nextSendCardOperator()
	r.operatorNow = next
	return next, nil
}

func (r *Room) gameFinish() *twoonone_pb.SendCardResponse_GameFinishEvent {
	r.operatorNow = nil
	var winis twoonone_pb.Role
	if len(r.landowner.GetCards()) == 0 {
		winis = twoonone_pb.Role_Landowner
	} else {
		winis = twoonone_pb.Role_Farmer
	}
	spring := false
	if winis == twoonone_pb.Role_Landowner {
		if len(r.farmers[0].GetCards()) == 17 && len(r.farmers[1].GetCards()) == 17 {
			r.multiple *= 2
			spring = true
		}
	}
	changeCoin := r.basicCoin * float64(r.multiple)
	if winis == twoonone_pb.Role_Landowner { //地主获胜
		r.landowner.IncCoin(changeCoin)
		r.landowner.IncWinCount()
		for _, v := range r.farmers {
			v.DecCoin(changeCoin / 2)
			v.IncLoseCount()
		}
	} else { //农民获胜
		for _, v := range r.farmers {
			v.IncCoin(changeCoin / 2)
			v.IncWinCount()
		}
		r.landowner.DecCoin(changeCoin)
		r.landowner.IncLoseCount()
	}
	return &twoonone_pb.SendCardResponse_GameFinishEvent{
		Landowner: insidePlayerToPlayerInfo(r.landowner),
		Farmer1:   insidePlayerToPlayerInfo(r.farmers[0]),
		Farmer2:   insidePlayerToPlayerInfo(r.farmers[1]),
		Winner:    winis,
		Spring:    spring,
	}
}

func insidePlayerToPlayerInfo(p *player.Player) *twoonone_pb.PlayerInfo {
	if p == nil {
		return nil
	}
	pi := &twoonone_pb.PlayerInfo{
		AccountInfo: &twoonone_pb.PlayerAccountInfo{
			Id:                    p.GetId(),
			Name:                  p.GetName(),
			WinCount:              int32(p.GetWinCount()),
			LoseCount:             int32(p.GetLoseCount()),
			Coin:                  p.GetCoin(),
			LastGetDailyTimestamp: p.GetLastGetDailyTimestamp(),
		},
	}
	if p.GetRoomHash() != "" {
		pi.TableInfo = &twoonone_pb.PlayerTableInfo{
			RoomHash:           p.GetRoomHash(),
			Cards:              p.GetCards(),
			RobLandownerAction: p.GetRobLandownerAction(),
		}
	}
	return pi
}
func (r *Room) playerUpdateToDatabase() error {
	for _, p := range r.players {
		origin := p.GetOriginInfo()
		coin := p.GetCoin() - origin.Coin
		wincount := p.GetWinCount() - int(origin.WinCount)
		losecount := p.GetLoseCount() - int(origin.LoseCount)
		if coin != 0 {
			if err := db.IncCoin(origin.Id, coin); err != nil {
				return err
			}
		}
		if wincount > 0 {
			if err := db.IncWinCount(origin.Id, wincount); err != nil {
				return err
			}
		}
		if losecount > 0 {
			if err := db.IncLoseCount(origin.Id, losecount); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Room) playerSendCard(p *player.Player, sendcards []twoonone_pb.Card, ct twoonone_pb.CardType, cs int, cc int) *twoonone_pb.Errors {
	if !p.DeleteCards(sendcards) {
		return twoonone_pb.Errors_PlayerCardNoExist.Enum()
	}
	// 特殊牌型倍率翻倍
	switch ct {
	case twoonone_pb.CardType_KingBomb, twoonone_pb.CardType_Bomb:
		r.multiple *= 2
	}
	r.cardPool = append(r.cardPool, &SendCard{
		SenderInfo:        p,
		SendCards:         sendcards,
		SendCardType:      ct,
		SendCardSize:      cs,
		SendCardContinous: cc,
	})
	return nil
}

func (r *Room) nextSendCardOperator() *player.Player {
	ps := r.GetPlayers()
	for i, v := range ps {
		if r.GetOperatorNow().GetId() == v.GetId() {
			if len(ps)-1 == i {
				return ps[0]
			}
			return ps[i+1]
		}
	}
	return nil
}

// 匹配牌类型
func (r *Room) matchCardType(cards []twoonone_pb.Card) twoonone_pb.CardType {
	//结果经降序后选maths[0]，abcd...都是递增关系，abcd...代表card
	var (
		matchs []twoonone_pb.CardType
	)
	//升序
	sort.Slice(cards, func(i, j int) bool {
		return cards[i] < cards[j]
	})
	//计数
	counts := r.countCard(cards)
	//匹配简单牌型
	switch len(counts) {
	case 1: //可能为单牌，双牌，三牌，炸弹
		switch counts[0].count {
		case 1: //单牌
			matchs = append(matchs, twoonone_pb.CardType_Single)
		case 2: //双牌
			matchs = append(matchs, twoonone_pb.CardType_Double)
		case 3: //三牌
			matchs = append(matchs, twoonone_pb.CardType_ThreeCard)
		case 4: //炸弹
			matchs = append(matchs, twoonone_pb.CardType_Bomb)
		}
	case 2: //可能为三带单牌，三带双牌，王炸，纯飞机(aaabbb)
		//匹配三带双牌，前正向匹配(aaabb)，后反向匹配(bbbaa)
		if counts[0].count == 3 && counts[1].count == 2 || counts[0].count == 2 && counts[1].count == 3 {
			matchs = append(matchs, twoonone_pb.CardType_ThreeWithDouble)
		}
		//匹配三带单牌
		if counts[0].card == twoonone_pb.Card_Joker || counts[1].card == twoonone_pb.Card_King {
		} else if counts[0].count == 3 && counts[1].count == 1 || counts[0].count == 1 && counts[1].count == 3 {
			matchs = append(matchs, twoonone_pb.CardType_ThreeWithSingle)
		}
		//匹配王炸
		if counts[0].card == twoonone_pb.Card_Joker && counts[1].card == twoonone_pb.Card_King {
			matchs = append(matchs, twoonone_pb.CardType_KingBomb)
		}
		//匹配纯飞机(aaabbb)
		if counts[0].card == twoonone_pb.Card_Two || counts[1].card != twoonone_pb.Card_Two { //确定ab不为2
		} else if counts[0].card+1 == counts[1].card { //确定a+1 == b
			if counts[0].count == 3 && counts[1].count == 3 { //确定a、b数量都为3
				matchs = append(matchs, twoonone_pb.CardType_AirSequence)
			}
		}
	case 3: //可能为四带两单牌，四带两双牌，连对(aabbcc)，纯飞机(aaabbbccc)
		//匹配四带两单牌，前正向匹配(aaaabc)，中复杂匹配(abbbbc)，后反向匹配(abcccc)
		if counts[0].count == 4 && counts[1].count == 1 && counts[2].count == 1 || counts[0].count == 1 && counts[1].count == 4 && counts[2].count == 1 || counts[0].count == 1 && counts[1].count == 1 && counts[2].count == 4 {
			matchs = append(matchs, twoonone_pb.CardType_FourWithTwoSingle)
		}
		//匹配四带两双牌，前正向匹配(aaaabbcc)，中复杂匹配(aabbbbcc)，后反向匹配(aabbcccc)
		if counts[0].count == 4 && counts[1].count == 2 && counts[2].count == 2 || counts[0].count == 2 && counts[1].count == 4 && counts[2].count == 2 || counts[0].count == 2 && counts[1].count == 2 && counts[2].count == 4 {
			matchs = append(matchs, twoonone_pb.CardType_FourWithTwoDouble)
		}
		//匹配连对(aabbcc)
		switch {
		case counts[0].card == twoonone_pb.Card_Two || counts[1].card == twoonone_pb.Card_Two || counts[2].card == twoonone_pb.Card_Two: //确定abc都不为2
		case counts[0].card+1 != counts[1].card && counts[0].card+2 != counts[2].card: //确定a+1 != b && a+2 != c
		case counts[1].card+1 != counts[2].card: //确定b+1 != c
		default:
			if counts[0].count == 2 && counts[1].count == 2 && counts[2].count == 2 { //确定a、b、c值都为3
				matchs = append(matchs, twoonone_pb.CardType_DoubleSequence)
			}
		}
		//匹配纯飞机(aaabbbccc)
		switch {
		case counts[0].card == twoonone_pb.Card_Two || counts[1].card == twoonone_pb.Card_Two || counts[2].card == twoonone_pb.Card_Two: //确定abc都不为2
		case counts[0].card+1 != counts[1].card && counts[0].card+2 != counts[2].card: //确定a+1 != b && a+2 != c
		case counts[1].card+1 != counts[2].card: //确定b+1 != c
		default:
			if counts[0].count == 3 && counts[1].count == 3 && counts[2].count == 3 { //确定a、b、c值都为3
				matchs = append(matchs, twoonone_pb.CardType_AirSequence)
			}
		}
	default:
		//复杂匹配
		//可能为顺子，连对，纯飞机，带两单牌飞机，带两双牌飞机
		//匹配顺子
		match := func(counts []*cardCount) bool {
			//确保counts >= 5, 至少: abcde
			if len(counts) < 5 {
				return false
			}
			for n := 0; n != len(counts); n++ {
				if n != len(counts)-1 {
					if counts[n].card+1 != counts[n+1].card { //确定k(n)+1 == k(n+1)
						return false
					}
				}
				if counts[n].count != 1 { //确定v(n) == 1
					return false
				}
				switch counts[n].card {
				case twoonone_pb.Card_Joker, twoonone_pb.Card_King, twoonone_pb.Card_Two: //确定k[n]不为小王/大王/2
					return false
				}
			}
			return true
		}
		if match(counts) {
			matchs = append(matchs, twoonone_pb.CardType_SingleSequence)
		}
		//匹配连对
		match = func(counts []*cardCount) bool {
			//确保counts >= 3, 至少: aabbcc
			if len(counts) < 3 {
				return false
			}
			for n := 0; n != len(counts); n++ {
				if counts[n].card == twoonone_pb.Card_Two { //确定k(n) != 2
					return false
				}
				if n != len(counts)-1 {
					if counts[n].card+1 != counts[n+1].card { //确定k(n)+1 == k(n+1)
						return false
					}
				}
				if counts[n].count != 2 { //确定v(n) == 2
					return false
				}
			}
			return true
		}
		if match(counts) {
			matchs = append(matchs, twoonone_pb.CardType_DoubleSequence)
		}
		//匹配纯飞机
		match = func(counts []*cardCount) bool {
			//确保counts >= 2, 至少: aaabbb
			if len(counts) < 2 {
				return false
			}
			for n := 0; n != len(counts); n++ {
				if n != len(counts)-1 {
					if counts[n].card+1 != counts[n+1].card { //确定k(n)+1 == k(n+1)
						return false
					}
				}
				if counts[n].count != 3 { //确定v(n) == 3
					return false
				}
			}
			return true
		}
		if match(counts) {
			matchs = append(matchs, twoonone_pb.CardType_AirSequence)
		}
		//匹配带两单牌飞机
		//正向匹配：x,y = n+1, n+2
		//复杂匹配：x,y = n-1, n+1
		//反向匹配：x,y = n-1, n-2
		//正向匹配(nnn...xy)，复杂匹配(xnnn...y)，反向匹配(yxnnn...)
		match_type := func(counts []*cardCount) int {
			const (
				fd = iota
				complex
				rd
				no_match
			)
			switch {
			case counts[len(counts)-1].count == 1 && counts[(len(counts)-1)-1].count == 1: //正向
				switch counts[len(counts)-1].card { //确定y不为小王/大王
				case twoonone_pb.Card_Joker, twoonone_pb.Card_King:
					return no_match
				default:
					switch counts[(len(counts)-1)-1].card { //确定x不为小王/大王
					case twoonone_pb.Card_Joker, twoonone_pb.Card_King:
						return no_match
					default:
						return fd
					}
				}
			case counts[0].count == 1 && counts[len(counts)-1].count == 1: //复杂
				switch counts[0].card { //确定x不为小王/大王
				case twoonone_pb.Card_Joker, twoonone_pb.Card_King:
					return no_match
				default:
					switch counts[len(counts)-1].card { //确定y不为小王/大王
					case twoonone_pb.Card_Joker, twoonone_pb.Card_King:
						return no_match
					default:
						return complex
					}
				}
			case counts[0].count == 1 && counts[1].count == 1: //反向
				switch counts[0].card { //确定y不为小王/大王
				case twoonone_pb.Card_Joker, twoonone_pb.Card_King:
					return no_match
				default:
					switch counts[1].card { //确定x不为小王/大王
					case twoonone_pb.Card_Joker, twoonone_pb.Card_King:
						return no_match
					default:
						return rd
					}
				}
			}
			return no_match //无匹配
		}
		match_1 := func(counts []*cardCount, t int) bool {
			const (
				fd = iota
				complex
				rd
				no_match
			)
			//确保len(counts)-2 >= 2, 至少: aaabbbcd
			if len(counts)-2 < 2 {
				return false
			}
			switch t {
			case fd:
				switch {
				case counts[(len(counts)-1)].count == 1 && counts[(len(counts)-1)-1].count == 1: //确定x和y的v为1
				default:
					return false
				}
				counts = counts[:(len(counts)-1)-1] //(nnn...xy) -> (nnn...)
			case complex:
				switch {
				case counts[(len(counts)-1)].count == 1 && counts[0].count == 1: //确定x和y的v为1
				default:
					return false
				}
				counts = counts[:(len(counts) - 1)] //(xnnn...y) -> (xnnn...)
				counts = counts[1:]                 //(xnnn...) -> (nnn...)
			case rd:
				switch {
				case counts[0].count == 1 && counts[1].count == 1: //确定x和y的v为1
				default:
					return false
				}
				counts = counts[2:] //(yxnnn...) -> (nnn...)
			case no_match:
				return false
			}
			for n := 0; n != len(counts); n++ {
				if counts[n].card == twoonone_pb.Card_Two { //确定k(n) != 2
					return false
				}
				if n != len(counts)-1 {
					if counts[n].card+1 != counts[n+1].card { //确定k(n)+1 == k(n+1)
						return false
					}
				}
				if counts[n].count != 3 { //确定v(n) == 3
					return false
				}
			}
			return true
		}
		if match_1(counts, match_type(counts)) {
			matchs = append(matchs, twoonone_pb.CardType_AirSequenceWithTwoSingle)
		}
		//匹配带两双牌飞机
		//正向匹配：x,y = n+1, n+2
		//复杂匹配：x,y = n-1, n+1
		//反向匹配：x,y = n-1, n-2
		//正向匹配(nnn...xxyy)，复杂匹配(xxnnn...yy)，反向匹配(yyxxnnn...)
		match_type = func(counts []*cardCount) int {
			const (
				fd = iota
				complex
				rd
				no_match
			)
			switch {
			case counts[(len(counts)-1)].count == 2 && counts[(len(counts)-1)-1].count == 2: //正向
				return fd
			case counts[0].count == 2 && counts[(len(counts)-1)].count == 2: //复杂
				return complex
			case counts[0].count == 2 && counts[1].count == 2: //反向
				return rd
			}
			return no_match //无匹配
		}
		match_1 = func(counts []*cardCount, t int) bool {
			const (
				fd = iota
				complex
				rd
				no_match
			)
			//确保len(counts)-2 >= 2, 至少: aaabbbccdd
			if len(counts)-2 < 2 {
				return false
			}
			switch t {
			case fd:
				switch {
				case counts[(len(counts)-1)].count == 2 && counts[(len(counts)-1)-1].count == 2: //确定x和y的v为1
				default:
					return false
				}
				counts = counts[:(len(counts)-1)-1] //(nnn...xxyy) -> (nnn...)
			case complex:
				switch {
				case counts[(len(counts)-1)].count == 2 && counts[0].count == 2: //确定x和y的v为1
				default:
					return false
				}
				counts = counts[:(len(counts) - 1)] //(xxnnn...yy) -> (xxnnn...)
				counts = counts[1:]                 //(xxnnn...) -> (nnn...)
			case rd:
				switch {
				case counts[0].count == 2 && counts[1].count == 2: //确定x和y的v为1
				default:
					return false
				}
				counts = counts[2:] //(yyxxnnn...) -> (nnn...)
			case no_match:
				return false
			}
			for n := 0; n != len(counts); n++ {
				if counts[n].card == twoonone_pb.Card_Two { //确定k(n) != 2
					return false
				}
				if n != len(counts)-1 {
					if counts[n].card+1 != counts[n+1].card { //确定k(n)+1 == k(n+1)
						return false
					}
				}
				if counts[n].count != 3 { //确定v(n) == 3
					return false
				}
			}
			return true
		}
		if match_1(counts, match_type(counts)) {
			matchs = append(matchs, twoonone_pb.CardType_AirSequenceWithTwoDouble)
		}
	}
	if len(matchs) != 0 {
		if len(matchs) > 1 {
			sort.Slice(matchs, func(i, j int) bool {
				return matchs[i] > matchs[j]
			})
		}
		return matchs[0]
	} else {
		//无匹配
		return twoonone_pb.CardType_Unknown
	}
}

// 计算指定牌类型的连续次数
func (r *Room) calcCardContinous(cards []twoonone_pb.Card, cardtype twoonone_pb.CardType) int {
	switch cardtype {
	default:
		return 0
	case twoonone_pb.CardType_SingleSequence:
		return len(cards)
	case twoonone_pb.CardType_DoubleSequence:
		return len(cards) / 2
	case twoonone_pb.CardType_AirSequence:
		return len(cards) / 3
	case twoonone_pb.CardType_AirSequenceWithTwoSingle:
		return (len(cards) - 2) / 3
	case twoonone_pb.CardType_AirSequenceWithTwoDouble:
		return (len(cards) - 2*2) / 3
	}
}

// 计算指定牌类型的大小
func (r *Room) calcCardSize(cards []twoonone_pb.Card, cardtype twoonone_pb.CardType) int {
	//升序
	sort.Slice(cards, func(i, j int) bool {
		return cards[i] < cards[j]
	})
	//计数
	counts := r.countCard(cards)
	var cardsize int
	switch cardtype {
	case twoonone_pb.CardType_ThreeWithSingle:
		//上正向匹配(aaab)，下反向匹配(bbba)
		switch {
		case counts[0].count == 3 && counts[1].count == 1:
			cardsize += int(counts[0].card) * counts[0].count
		case counts[0].count == 1 && counts[1].count == 3:
			cardsize += int(counts[1].card) * counts[1].count
		}
	case twoonone_pb.CardType_ThreeWithDouble:
		//上正向匹配(aaabb)，下反向匹配(bbbaa)
		switch {
		case counts[0].count == 3 && counts[1].count == 2:
			cardsize += counts[0].count * int(counts[0].card)
		case counts[0].count == 2 && counts[1].count == 3:
			cardsize += counts[1].count * int(counts[1].card)
		}
	case twoonone_pb.CardType_FourWithTwoSingle:
		//上正向匹配(aaaabc)，中复杂匹配(abbbbc)，下反向匹配(abcccc)
		switch {
		case counts[0].count == 4 && counts[1].count == 1 && counts[2].count == 1:
			cardsize += counts[0].count * int(counts[0].card)
		case counts[0].count == 1 && counts[1].count == 4 && counts[2].count == 1:
			cardsize += counts[1].count * int(counts[1].card)
		case counts[0].count == 1 && counts[1].count == 1 && counts[2].count == 4:
			cardsize += counts[2].count * int(counts[2].card)
		}
	case twoonone_pb.CardType_FourWithTwoDouble:
		//前正向匹配(aaaabbcc)，中复杂匹配(aabbbbcc)，后反向匹配(aabbcccc)
		switch {
		case counts[0].count == 4 && counts[1].count == 2 && counts[2].count == 2:
			cardsize += counts[0].count * int(counts[0].card)
		case counts[0].count == 2 && counts[1].count == 4 && counts[2].count == 2:
			cardsize += counts[1].count * int(counts[1].card)
		case counts[0].count == 2 && counts[1].count == 2 && counts[2].count == 4:
			cardsize += counts[2].count * int(counts[2].card)
		}
	case twoonone_pb.CardType_AirSequenceWithTwoSingle:
		match_type := func(counts []*cardCount) int {
			const (
				fd = iota
				complex
				rd
				no_match
			)
			switch {
			case counts[(len(counts)-1)].count == 1 && counts[(len(counts)-1)-1].count == 1: //正向
				return fd
			case counts[0].count == 1 && counts[(len(counts)-1)].count == 1: //复杂
				return complex
			case counts[0].count == 1 && counts[1].count == 1: //反向
				return rd
			}
			return no_match //无匹配
		}
		calculate := func(counts []*cardCount, t int) int {
			var size int
			const (
				fd = iota
				complex
				rd
				no_match
			)
			switch t {
			case fd:
				counts = counts[:(len(counts)-1)-1] //(nnn...xxyy) -> (nnn...)
			case complex:
				counts = counts[:(len(counts) - 1)] //(xxnnn...yy) -> (xxnnn...)
				counts = counts[1:]                 //(xxnnn...) -> (nnn...)
			case rd:
				counts = counts[2:] //(yyxxnnn...) -> (nnn...)
			case no_match:
				return -1
			}
			for _, v := range counts {
				size += v.count * int(v.card)
			}
			return size
		}
		size := calculate(counts, match_type(counts))
		if size != -1 {
			cardsize += size
		}
	case twoonone_pb.CardType_AirSequenceWithTwoDouble:
		match_type := func(counts []*cardCount) int {
			const (
				fd = iota
				complex
				rd
				no_match
			)
			switch {
			case counts[(len(counts)-1)].count == 2 && counts[(len(counts)-1)-1].count == 2: //正向
				return fd
			case counts[0].count == 2 && counts[(len(counts)-1)].count == 2: //复杂
				return complex
			case counts[0].count == 2 && counts[1].count == 2: //反向
				return rd
			}
			return no_match //无匹配
		}
		calculate := func(counts []*cardCount, t int) int {
			var size int
			const (
				fd = iota
				complex
				rd
				no_match
			)
			switch t {
			case fd:
				counts = counts[:(len(counts)-1)-1] //(nnn...xxyy) -> (nnn...)
			case complex:
				counts = counts[:(len(counts) - 1)] //(xxnnn...yy) -> (xxnnn...)
				counts = counts[1:]                 //(xxnnn...) -> (nnn...)
			case rd:
				counts = counts[2:] //(yyxxnnn...) -> (nnn...)
			case no_match:
				return -1
			}
			for _, v := range counts {
				size += v.count * int(v.card)
			}
			return size
		}
		size := calculate(counts, match_type(counts))
		if size != -1 {
			cardsize += size
		}
	default:
		for _, v := range cards {
			cardsize += int(v)
		}
	}
	return cardsize
}

type cardCount struct {
	card  twoonone_pb.Card
	count int
}

// 牌计数器
func (r *Room) countCard(cards []twoonone_pb.Card) []*cardCount {
	var (
		counter [15]int //0-14代表对应的card
		counts  []*cardCount
	)
	for _, card := range cards {
		counter[card]++
	}
	for c, count := range counter {
		//跳过不存在的牌
		if count == 0 {
			continue
		}
		counts = append(counts, &cardCount{
			card:  twoonone_pb.Card(c),
			count: count,
		})
	}
	return counts
}

func (r *Room) nextRobLandownerOperator() *player.Player {
	for _, v := range r.players {
		if v.GetRobLandownerAction() == nil {
			return v
		}
	}
	return nil
}

func (r *Room) startSendCard() {
	// 为地主发牌
	lo := r.GetLandowner()
	lo.AddCards(r.landownerCards[:])
	// 设置农民
	for _, v := range r.GetPlayers() {
		if v.GetId() != lo.GetId() {
			switch {
			case r.farmers[0] == nil:
				r.farmers[0] = v
			case r.farmers[1] == nil:
				r.farmers[1] = v
			}
		}
	}
	// 设置阶段
	r.stage = twoonone_pb.RoomStage_SendingCards
	// 重排玩家
	r.resequence(&r.players)
	r.operatorNow = lo
}

func (r *Room) getRobLandowners() []*player.Player {
	pis := []*player.Player{}
	for _, v := range r.players {
		action := v.GetRobLandownerAction()
		if action == nil {
			continue
		}
		if *action == twoonone_pb.RobLandownerActions_Rob {
			pis = append(pis, v)
		}
	}
	return pis
}

func (r *Room) generateCards() [54]twoonone_pb.Card {
	return [54]twoonone_pb.Card{
		0, 0, 0, 0,
		1, 1, 1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3,
		4, 4, 4, 4,
		5, 5, 5, 5,
		6, 6, 6, 6,
		7, 7, 7, 7,
		8, 8, 8, 8,
		9, 9, 9, 9,
		10, 10, 10, 10,
		11, 11, 11, 11,
		12, 12, 12, 12,
		13,
		14,
	}
}

func (r *Room) delete(id string) bool {
	for n := 0; n != len(r.players); n++ {
		p := r.players[n]
		if p.IsEmpty() {
			continue
		}
		if p.GetId() == id {
			if len(r.players) == 1 {
				r.players = r.players[:0]
			} else {
				r.players = append(r.players[:n], r.players[n+1:]...)
			}
			return true
		}
	}
	return false
}

func (r *Room) resequence(x any) {
	switch xPoint := x.(type) {
	case *[]*player.Player: //重排玩家
		x := *xPoint
		lo := r.GetLandowner()
		if lo != nil { //游戏已开始，重排顺序
			x[0] = lo
			x[1] = r.farmers[0]
			x[2] = r.farmers[1]
		}
	case *[]twoonone_pb.Card: //重排牌
		x := *xPoint
		rand.Shuffle(len(x), func(i, j int) {
			x[i], x[j] = x[j], x[i]
		})
	}
}

func (r *Room) GetPlayer(playerId string) (*player.Player, bool) {
	for _, v := range r.GetPlayers() {
		if v.IsEmpty() {
			continue
		}
		if v.GetId() == playerId {
			return v, true
		}
	}
	return nil, false
}

func (r *Room) GetHash() string {
	return r.hash
}

func (r *Room) GetPlayers() []*player.Player {
	return r.players
}

func (r *Room) GetOperatorNow() *player.Player {
	return r.operatorNow
}

func (r *Room) GetStage() twoonone_pb.RoomStage {
	return r.stage
}

func (r *Room) GetLandowner() *player.Player {
	return r.landowner
}

func (r *Room) GetFarmers() [2]*player.Player {
	return r.farmers
}

func (r *Room) GetLastCard() *SendCard {
	if len(r.cardPool) == 0 {
		return nil
	}
	return r.cardPool[len(r.cardPool)-1]
}

func (r *Room) GetCardPool() []*SendCard {
	return r.cardPool
}

func (r *Room) GetBasicCoin() float64 {
	return r.basicCoin
}

func (r *Room) GetMultiple() int {
	return r.multiple
}

func (r *Room) GetLandownerCard() [3]twoonone_pb.Card {
	return r.landownerCards
}
