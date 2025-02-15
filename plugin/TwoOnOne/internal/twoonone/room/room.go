package room

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"sync"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone/card"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone/event"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone/player"
	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/twmb/murmur3"
	"github.com/zeromicro/go-zero/core/logx"
)

type Room struct {
	lock           sync.Mutex
	hash           string
	players        []*player.Player
	operatorNow    *player.Player
	landowner      *player.Player
	farmers        [2]*player.Player
	stage          twoonone_pb.RoomStage
	basicCoin      float64
	multiple       int
	sendCards      []*SendCard
	landownerCards [3]card.Card
	event          *event.EventStream
}

type SendCard struct {
	SenderInfo        *player.Player
	SendCards         []card.Card
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
	hash := hash()
	return &Room{
		hash:           hash,
		players:        []*player.Player{},
		operatorNow:    nil,
		landowner:      nil,
		farmers:        [2]*player.Player{},
		stage:          twoonone_pb.RoomStage_ROOM_STAGE_WAITTING_START,
		basicCoin:      basicCoin,
		multiple:       multiple,
		sendCards:      []*SendCard{},
		landownerCards: [3]card.Card{},
		event:          event.NewEventStream(hash),
	}
}

func (r *Room) Join(logger logx.Logger, id, name string, coin float64) *types.AppError {
	r.lock.Lock()
	r.lock.Unlock()
	if r.stage != twoonone_pb.RoomStage_ROOM_STAGE_WAITTING_START {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_STARTED, "")
	}
	if len(r.players) == 3 {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_FULL, "")

	}
	if _, ok := r.GetPlayer(id); ok {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_EXIST_PLAYER, "")
	}
	if coin < r.basicCoin {
		return types.NewError(twoonone_pb.Error_ERROR_PLAYER_COIN_LT_ROOM_COIN, "")
	}
	p := player.New(&twoonone_pb.PlayerInfo{
		User: &twoonone_pb.PlayerInfo_UserInfo{
			Id:   id,
			Name: name,
		},
		Table: &twoonone_pb.PlayerInfo_TableInfo{
			RoomHash:         r.hash,
			RoblandownerInfo: &twoonone_pb.RobLandownerInfo{},
		},
	})
	r.players = append(r.players, p)
	r.event.Emit(&twoonone_pb.EventRoomResponse{
		Body: &twoonone_pb.EventRoomResponse_RoomJoinPlayer_{
			RoomJoinPlayer: &twoonone_pb.EventRoomResponse_RoomJoinPlayer{
				JoinerInfo:  player.FormatInternalPlayer2Protobuf(p),
				PlayerInfos: player.FormatInternalPlayers2Protobuf(r.players),
			},
		},
	})
	return nil
}

func (r *Room) Exit(logger logx.Logger, playerId string) *types.AppError {
	r.lock.Lock()
	r.lock.Unlock()
	if r.stage != twoonone_pb.RoomStage_ROOM_STAGE_WAITTING_START {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_STARTED, "")
	}
	p, ok := r.GetPlayer(playerId)
	if !ok {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST_PLAYER, "")
	}
	r.delete(playerId)
	r.event.Emit(&twoonone_pb.EventRoomResponse{
		Body: &twoonone_pb.EventRoomResponse_RoomExitPlayer_{
			RoomExitPlayer: &twoonone_pb.EventRoomResponse_RoomExitPlayer{
				LeaverInfo:  player.FormatInternalPlayer2Protobuf(p),
				PlayerInfos: player.FormatInternalPlayers2Protobuf(r.players),
			},
		},
	})
	return nil
}

func (r *Room) Start(logger logx.Logger) *types.AppError {
	r.lock.Lock()
	r.lock.Unlock()
	if r.GetStage() != twoonone_pb.RoomStage_ROOM_STAGE_WAITTING_START {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_STARTED, "")
	}
	if len(r.players) != 3 {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_FULL, "")
	}
	// 开一副牌
	cardsArr := generateCards()
	cards := cardsArr[:]
	r.resequenceCards(&cards)
	// 为玩家发牌
	for n, N := 0, 0; N != len(r.players); N++ {
		playerCards := cards[n : n+17]
		n += 17
		r.resequenceCards(&playerCards)
		r.players[N].AddCards(playerCards)
	}
	// 设置地主牌
	r.setLandownerCards([3]card.Card{
		cards[len(cards)-1],
		cards[len(cards)-2],
		cards[len(cards)-3],
	})
	// 进入下一阶段
	r.startRobLandowner()
	r.event.Emit(&twoonone_pb.EventRoomResponse{
		Body: &twoonone_pb.EventRoomResponse_RoomStarted_{
			RoomStarted: &twoonone_pb.EventRoomResponse_RoomStarted{
				NextOperatorInfo: player.FormatInternalPlayer2Protobuf(r.operatorNow),
			},
		},
	})
	return nil
}

func (r *Room) startRobLandowner() {
	r.stage = twoonone_pb.RoomStage_ROOM_STAGE_ROB_LANDOWNERING
	r.operatorNow = r.players[rand.Intn(len(r.players))]
}

func (r *Room) setLandownerCards(cards [3]card.Card) {
	copy(r.landownerCards[:], cards[:])
}

func (r *Room) RobLandowner(logger logx.Logger, p *player.Player) *types.AppError {
	r.lock.Lock()
	r.lock.Unlock()
	if p.GetRoomHash() != r.GetHash() {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST_PLAYER, "")
	}
	if r.stage != twoonone_pb.RoomStage_ROOM_STAGE_ROB_LANDOWNERING {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_ROB_LANDOWNERING, "")
	}
	if r.operatorNow.GetId() != p.GetId() {
		return types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_OPERATOR, "")
	}
	robs := r.getRobLandowners()
	p.SetRobLandownerAction(twoonone_pb.RobLandownerInfo_ACTION_ROB)
	if len(r.getRobLandowners()) > len(robs) {
		r.multiple *= 2
		r.event.Emit(&twoonone_pb.EventRoomResponse{
			Body: &twoonone_pb.EventRoomResponse_RoblandownerContinuousRob{
				RoblandownerContinuousRob: &twoonone_pb.EventRoomResponse_RobLandownerContinuousRob{
					Multiple: int32(r.multiple),
				},
			},
		})
	}

	takes := r.getTakeRobLandowners()
	if len(takes) == len(r.players) { // 全部已参与
		if len(takes) == 0 { //无人抢地主
			return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_ROB_LANDOWNER, "")
		} else if len(takes) == 1 { //仅一人抢地主
			r.landowner = takes[0]
		} else { //一人以上抢地主
			// 按操作时间降序
			sort.Slice(takes, func(i, j int) bool {
				return takes[i].GetRobLandownerActionTime().UnixNano() > takes[j].GetRobLandownerActionTime().UnixNano()
			})
			r.landowner = takes[0]
		}
		r.event.Emit(&twoonone_pb.EventRoomResponse{
			Body: &twoonone_pb.EventRoomResponse_RoomRobLandowner_{
				RoomRobLandowner: &twoonone_pb.EventRoomResponse_RoomRobLandowner{
					OperatorInfo:     player.FormatInternalPlayer2Protobuf(p),
					NextOperatorInfo: nil,
				},
			},
		})
	} else { //还有未参与
		next := r.nextRobLandownerOperator()
		r.event.Emit(&twoonone_pb.EventRoomResponse{
			Body: &twoonone_pb.EventRoomResponse_RoomRobLandowner_{
				RoomRobLandowner: &twoonone_pb.EventRoomResponse_RoomRobLandowner{
					OperatorInfo:     player.FormatInternalPlayer2Protobuf(p),
					NextOperatorInfo: player.FormatInternalPlayer2Protobuf(next),
				},
			},
		})
		return nil
	}
	r.startSendCard()
	r.event.Emit(&twoonone_pb.EventRoomResponse{
		Body: &twoonone_pb.EventRoomResponse_RoblandownerIntoSendingcard{
			RoblandownerIntoSendingcard: &twoonone_pb.EventRoomResponse_RobLandownerIntoSendingCard{
				SendcarderInfo: player.FormatInternalPlayer2Protobuf(r.operatorNow),
				LandownerCards: card.FormatInternalCards2Protobuf(r.landownerCards[:]),
			},
		},
	})
	return nil
}

func (r *Room) NoRobLandowner(logger logx.Logger, p *player.Player) *types.AppError {
	r.lock.Lock()
	r.lock.Unlock()
	if p.GetRoomHash() != r.GetHash() {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST_PLAYER, "")
	}
	if r.stage != twoonone_pb.RoomStage_ROOM_STAGE_ROB_LANDOWNERING {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_ROB_LANDOWNERING, "")
	}
	if r.operatorNow.GetId() != p.GetId() {
		return types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_OPERATOR, "")
	}
	p.SetRobLandownerAction(twoonone_pb.RobLandownerInfo_ACTION_NO_ROB)

	takes := r.getTakeRobLandowners()
	if len(takes) == len(r.players) { // 全部已参与
		if len(takes) == 0 { //无人抢地主
			return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_ROB_LANDOWNER, "")
		} else if len(takes) == 1 { //仅一人抢地主
			r.landowner = takes[0]
		} else { //一人以上抢地主
			// 按操作时间降序
			sort.Slice(takes, func(i, j int) bool {
				return takes[i].GetRobLandownerActionTime().UnixNano() > takes[j].GetRobLandownerActionTime().UnixNano()
			})
			r.landowner = takes[0]
		}
		r.event.Emit(&twoonone_pb.EventRoomResponse{
			Body: &twoonone_pb.EventRoomResponse_RoomNorobLandowner{
				RoomNorobLandowner: &twoonone_pb.EventRoomResponse_RoomNoRobLandowner{
					OperatorInfo:     player.FormatInternalPlayer2Protobuf(p),
					NextOperatorInfo: nil,
				},
			},
		})
	} else { //还有未参与
		next := r.nextRobLandownerOperator()
		r.event.Emit(&twoonone_pb.EventRoomResponse{
			Body: &twoonone_pb.EventRoomResponse_RoomNorobLandowner{
				RoomNorobLandowner: &twoonone_pb.EventRoomResponse_RoomNoRobLandowner{
					OperatorInfo:     player.FormatInternalPlayer2Protobuf(p),
					NextOperatorInfo: player.FormatInternalPlayer2Protobuf(next),
				},
			},
		})
		return nil
	}
	r.startSendCard()
	r.event.Emit(&twoonone_pb.EventRoomResponse{
		Body: &twoonone_pb.EventRoomResponse_RoblandownerIntoSendingcard{
			RoblandownerIntoSendingcard: &twoonone_pb.EventRoomResponse_RobLandownerIntoSendingCard{
				SendcarderInfo: player.FormatInternalPlayer2Protobuf(r.operatorNow),
				LandownerCards: card.FormatInternalCards2Protobuf(r.landownerCards[:]),
			},
		},
	})
	return nil
}

func (r *Room) SendCard(logger logx.Logger, p *player.Player, sendcards []*twoonone_pb.Card) *types.AppError {
	r.lock.Lock()
	r.lock.Unlock()
	if p.GetRoomHash() != r.hash {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST_PLAYER, "")
	}
	if r.stage != twoonone_pb.RoomStage_ROOM_STAGE_SENDING_CARD {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_SENDING_CARD, "")
	}
	if r.operatorNow.GetId() != p.GetId() {
		return types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_OPERATOR, "")
	}
	cards := card.FormatProtobuf2InternalCards(sendcards)
	if serr := r.sendCard(logger, p, cards); serr != nil {
		return serr
	}
	return nil
}

func (r *Room) NoSendCard(logger logx.Logger, p *player.Player) *types.AppError {
	r.lock.Lock()
	r.lock.Unlock()
	if p.GetRoomHash() != r.hash {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST_PLAYER, "")
	}
	if r.stage != twoonone_pb.RoomStage_ROOM_STAGE_SENDING_CARD {
		return types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_SENDING_CARD, "")
	}
	if r.operatorNow.GetId() != p.GetId() {
		return types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_OPERATOR, "")
	}
	if serr := r.noSendCard(logger, p); serr != nil {
		return serr
	}
	return nil
}

func (r *Room) sendCard(logger logx.Logger, p *player.Player, sendcards []card.Card) *types.AppError {
	lastcard := r.GetLastCard()
	// 正常出牌
	cardtype := r.matchCardType(sendcards)
	cardcontious := r.calcCardContinous(sendcards, cardtype)
	cardsize := r.calcCardSize(sendcards, cardtype)
	// 特殊情况判断
	if cardtype == twoonone_pb.CardType_CARD_TYPE_UNKNOWN { //未知牌型
		return types.NewError(twoonone_pb.Error_ERROR_SEND_CARD_TYPE_UNKNOWN, "")
	} else if lastcard == nil { //本局第一次出牌
		if serr := r.playerSendCard(logger, p, sendcards, cardtype, cardsize, cardcontious); serr != nil {
			return serr
		}
		return nil
	} else if lastcard.SenderInfo.GetId() == p.GetId() { //同一人操作
		if serr := r.playerSendCard(logger, p, sendcards, cardtype, cardsize, cardcontious); serr != nil {
			return serr
		}
		return nil
	} else if cardtype == twoonone_pb.CardType_CARD_TYPE_KING_BOOM { // 王炸
		if serr := r.playerSendCard(logger, p, sendcards, cardtype, cardsize, cardcontious); serr != nil {
			return serr
		}
		return nil
	} else if cardtype == twoonone_pb.CardType_CARD_TYPE_BOOM && lastcard.SendCardType == twoonone_pb.CardType_CARD_TYPE_BOOM { //上一副与当前都为炸弹
		if cardsize <= lastcard.SendCardSize {
			types.NewError(twoonone_pb.Error_ERROR_SEND_CARD_SIZE_LE_LAST_CARD_SIZE, "")
		}
		if serr := r.playerSendCard(logger, p, sendcards, cardtype, cardsize, cardcontious); serr != nil {
			return serr
		}
		return nil
	} else if cardtype == twoonone_pb.CardType_CARD_TYPE_BOOM && lastcard.SendCardType != twoonone_pb.CardType_CARD_TYPE_KING_BOOM { //上一副不为王炸且也不为炸弹
		if serr := r.playerSendCard(logger, p, sendcards, cardtype, cardsize, cardcontious); serr != nil {
			return serr
		}
		return nil
	}
	//完全正常出牌
	if cardtype != lastcard.SendCardType {
		types.NewError(twoonone_pb.Error_ERROR_SEND_CARD_TYPE_NE_LAST_CARD_TYPE, "")
	}
	if cardcontious != lastcard.SendCardContinous {
		types.NewError(twoonone_pb.Error_ERROR_SEND_CARD_CONTINUOUS_NE_LAST_CARD_CONTINUOUS, "")
	}
	if cardsize <= lastcard.SendCardSize {
		types.NewError(twoonone_pb.Error_ERROR_SEND_CARD_SIZE_LE_LAST_CARD_SIZE, "")
	}
	if serr := r.playerSendCard(logger, p, sendcards, cardtype, cardsize, cardcontious); serr != nil {
		return serr
	}
	return nil
}

func (r *Room) noSendCard(logger logx.Logger, p *player.Player) *types.AppError {
	lastcard := r.GetLastCard()
	// 特殊情况判断
	if lastcard == nil { //第一次出牌
		types.NewError(twoonone_pb.Error_ERROR_PLAYER_IS_ONLY_OPERATOR, "")
	} else if lastcard.SenderInfo.GetId() == p.GetId() { //上一次出牌为同一人
		types.NewError(twoonone_pb.Error_ERROR_PLAYER_IS_ONLY_OPERATOR, "")
	}
	next := r.nextSendCardOperator()
	r.operatorNow = next
	r.event.Emit(&twoonone_pb.EventRoomResponse{
		Body: &twoonone_pb.EventRoomResponse_RoomNoSendcard{
			RoomNoSendcard: &twoonone_pb.EventRoomResponse_RoomNoSendCard{
				OperatorInfo:     player.FormatInternalPlayer2Protobuf(p),
				NextOperatorInfo: player.FormatInternalPlayer2Protobuf(next),
			},
		},
	})
	return nil
}

func (r *Room) GetEvent() *event.EventStream {
	return r.event
}

func (r *Room) gameFinish(logger logx.Logger) *types.AppError {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.operatorNow = nil
	var winis twoonone_pb.Role
	if len(r.landowner.GetCards()) == 0 {
		winis = twoonone_pb.Role_ROLE_LANDOWNER
	} else {
		winis = twoonone_pb.Role_ROLE_FARMER
	}
	//春天
	if winis == twoonone_pb.Role_ROLE_LANDOWNER {
		if len(r.farmers[0].GetCards()) == 17 && len(r.farmers[1].GetCards()) == 17 {
			r.multiple *= 2
			r.event.Emit(&twoonone_pb.EventRoomResponse{Body: &twoonone_pb.EventRoomResponse_SendcardSpringNotice{
				SendcardSpringNotice: &twoonone_pb.EventRoomResponse_SendCardSpringNotice{
					Multiple: int32(r.multiple),
				},
			}})
		}
	}
	changeCoin := r.basicCoin * float64(r.multiple)
	if winis == twoonone_pb.Role_ROLE_LANDOWNER { //地主获胜
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
	// 结算至数据库
	for _, v := range r.players {
		if serr := v.UpdateToDatabaseAndClean(logger); serr != nil {
			return serr
		}
	}
	r.event.Emit(&twoonone_pb.EventRoomResponse{Body: &twoonone_pb.EventRoomResponse_GameFinish_{
		GameFinish: &twoonone_pb.EventRoomResponse_GameFinish{
			LandownerInfo: &twoonone_pb.EventRoomResponse_GameFinish_PlayerInfoExtra{
				PlayerInfo: player.FormatInternalPlayer2Protobuf(r.landowner),
				HandCards:  card.FormatInternalCards2Protobuf(r.landowner.GetCards()),
			},
			Farmer1Info: &twoonone_pb.EventRoomResponse_GameFinish_PlayerInfoExtra{
				PlayerInfo: player.FormatInternalPlayer2Protobuf(r.farmers[0]),
				HandCards:  card.FormatInternalCards2Protobuf(r.farmers[0].GetCards()),
			},
			Farmer2Info: &twoonone_pb.EventRoomResponse_GameFinish_PlayerInfoExtra{
				PlayerInfo: player.FormatInternalPlayer2Protobuf(r.farmers[1]),
				HandCards:  card.FormatInternalCards2Protobuf(r.farmers[1].GetCards()),
			},
			Winner: winis,
		},
	}})
	return nil
}

func (r *Room) playerSendCard(logger logx.Logger, p *player.Player, sendcards []card.Card, ct twoonone_pb.CardType, cs int, cc int) *types.AppError {
	if !p.DeleteCards(sendcards) {
		types.NewError(twoonone_pb.Error_ERROR_PLAYER_CARD_NO_EXIST, "")
	}
	// 特殊牌型倍率翻倍
	switch ct {
	case twoonone_pb.CardType_CARD_TYPE_KING_BOOM, twoonone_pb.CardType_CARD_TYPE_BOOM:
		r.multiple *= 2
		switch ct {
		case twoonone_pb.CardType_CARD_TYPE_BOOM:
			r.event.Emit(&twoonone_pb.EventRoomResponse{Body: &twoonone_pb.EventRoomResponse_SendcardBoomNotice{
				SendcardBoomNotice: &twoonone_pb.EventRoomResponse_SendCardBoomNotice{
					Multiple:       int32(r.multiple),
					SendcarderInfo: player.FormatInternalPlayer2Protobuf(p),
				},
			}})
		case twoonone_pb.CardType_CARD_TYPE_KING_BOOM:
			r.event.Emit(&twoonone_pb.EventRoomResponse{Body: &twoonone_pb.EventRoomResponse_SendcardKingboomNotice{
				SendcardKingboomNotice: &twoonone_pb.EventRoomResponse_SendCardKingBoomNotice{
					Multiple:       int32(r.multiple),
					SendcarderInfo: player.FormatInternalPlayer2Protobuf(p),
				},
			}})
		}
	}
	r.sendCards = append(r.sendCards, &SendCard{
		SenderInfo:        p,
		SendCards:         sendcards,
		SendCardType:      ct,
		SendCardSize:      cs,
		SendCardContinous: cc,
	})
	next := r.nextSendCardOperator()
	r.operatorNow = next
	r.event.Emit(&twoonone_pb.EventRoomResponse{
		Body: &twoonone_pb.EventRoomResponse_RoomSendcard{
			RoomSendcard: &twoonone_pb.EventRoomResponse_RoomSendCard{
				OperatorInfo:     player.FormatInternalPlayer2Protobuf(p),
				NextOperatorInfo: player.FormatInternalPlayer2Protobuf(next),
				Sendcards:        card.FormatInternalCards2Protobuf(sendcards),
			},
		},
	})
	switch l := len(p.GetCards()); l {
	case 1, 2:
		r.event.Emit(&twoonone_pb.EventRoomResponse{
			Body: &twoonone_pb.EventRoomResponse_SendcardCardnumberNotice{
				SendcardCardnumberNotice: &twoonone_pb.EventRoomResponse_SendCardCardNumberNotice{
					Number:           int32(l),
					NoticeTargetInfo: player.FormatInternalPlayer2Protobuf(p),
				},
			},
		})
	case 0:
		if serr := r.gameFinish(logger); serr != nil {
			return serr
		}
	}
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
func (r *Room) matchCardType(cards []card.Card) twoonone_pb.CardType {
	//结果经降序后选maths[0]，abcd...都是递增关系，abcd...代表card
	var (
		matchs []twoonone_pb.CardType
	)
	//升序
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Number < cards[j].Number
	})
	//计数
	counts := r.countCard(cards)
	//匹配简单牌型
	switch len(counts) {
	case 1: //可能为单牌，双牌，三牌，炸弹
		switch counts[0].count {
		case 1: //单牌
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_SINGLE)
		case 2: //双牌
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_DOUBLE)
		case 3: //三牌
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_THREE_CARD)
		case 4: //炸弹
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_BOOM)
		}
	case 2: //可能为三带单牌，三带双牌，王炸，纯飞机(aaabbb)
		//匹配三带双牌，前正向匹配(aaabb)，后反向匹配(bbbaa)
		if counts[0].count == 3 && counts[1].count == 2 || counts[0].count == 2 && counts[1].count == 3 {
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_THREE_WITH_DOUBLE)
		}
		//匹配三带单牌
		if counts[0].card.Number == card.JOKER || counts[1].card.Number == card.KING {
		} else if counts[0].count == 3 && counts[1].count == 1 || counts[0].count == 1 && counts[1].count == 3 {
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_THREE_WITH_SINGLE)
		}
		//匹配王炸
		if counts[0].card.Number == card.JOKER && counts[1].card.Number == card.KING {
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_KING_BOOM)
		}
		//匹配纯飞机(aaabbb)
		if counts[0].card.Number == card.TWO || counts[1].card.Number != card.TWO { //确定ab不为2
		} else if counts[0].card.Number+1 == counts[1].card.Number { //确定a+1 == b
			if counts[0].count == 3 && counts[1].count == 3 { //确定a、b数量都为3
				matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE)
			}
		}
	case 3: //可能为四带两单牌，四带两双牌，连对(aabbcc)，纯飞机(aaabbbccc)
		//匹配四带两单牌，前正向匹配(aaaabc)，中复杂匹配(abbbbc)，后反向匹配(abcccc)
		if counts[0].count == 4 && counts[1].count == 1 && counts[2].count == 1 || counts[0].count == 1 && counts[1].count == 4 && counts[2].count == 1 || counts[0].count == 1 && counts[1].count == 1 && counts[2].count == 4 {
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_FOUR_WITH_TWO_SINGLE)
		}
		//匹配四带两双牌，前正向匹配(aaaabbcc)，中复杂匹配(aabbbbcc)，后反向匹配(aabbcccc)
		if counts[0].count == 4 && counts[1].count == 2 && counts[2].count == 2 || counts[0].count == 2 && counts[1].count == 4 && counts[2].count == 2 || counts[0].count == 2 && counts[1].count == 2 && counts[2].count == 4 {
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_FOUR_WITH_TWO_DOUBLE)
		}
		//匹配连对(aabbcc)
		switch {
		case counts[0].card.Number == card.TWO || counts[1].card.Number == card.TWO || counts[2].card.Number == card.TWO: //确定abc都不为2
		case counts[0].card.Number+1 != counts[1].card.Number && counts[0].card.Number+2 != counts[2].card.Number: //确定a+1 != b && a+2 != c
		case counts[1].card.Number+1 != counts[2].card.Number: //确定b+1 != c
		default:
			if counts[0].count == 2 && counts[1].count == 2 && counts[2].count == 2 { //确定a、b、c值都为3
				matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_DOUBLE_SEQUENCE)
			}
		}
		//匹配纯飞机(aaabbbccc)
		switch {
		case counts[0].card.Number == card.TWO || counts[1].card.Number == card.TWO || counts[2].card.Number == card.TWO: //确定abc都不为2
		case counts[0].card.Number+1 != counts[1].card.Number && counts[0].card.Number+2 != counts[2].card.Number: //确定a+1 != b && a+2 != c
		case counts[1].card.Number+1 != counts[2].card.Number: //确定b+1 != c
		default:
			if counts[0].count == 3 && counts[1].count == 3 && counts[2].count == 3 { //确定a、b、c值都为3
				matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE)
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
					if counts[n].card.Number+1 != counts[n+1].card.Number { //确定k(n)+1 == k(n+1)
						return false
					}
				}
				if counts[n].count != 1 { //确定v(n) == 1
					return false
				}
				switch counts[n].card.Number {
				case card.JOKER, card.KING, card.TWO: //确定k[n]不为小王/大王/2
					return false
				}
			}
			return true
		}
		if match(counts) {
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_SINGLE_SEQUENCE)
		}
		//匹配连对
		match = func(counts []*cardCount) bool {
			//确保counts >= 3, 至少: aabbcc
			if len(counts) < 3 {
				return false
			}
			for n := 0; n != len(counts); n++ {
				if counts[n].card.Number == card.TWO { //确定k(n) != 2
					return false
				}
				if n != len(counts)-1 {
					if counts[n].card.Number+1 != counts[n+1].card.Number { //确定k(n)+1 == k(n+1)
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
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_DOUBLE_SEQUENCE)
		}
		//匹配纯飞机
		match = func(counts []*cardCount) bool {
			//确保counts >= 2, 至少: aaabbb
			if len(counts) < 2 {
				return false
			}
			for n := 0; n != len(counts); n++ {
				if n != len(counts)-1 {
					if counts[n].card.Number+1 != counts[n+1].card.Number { //确定k(n)+1 == k(n+1)
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
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE)
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
				switch counts[len(counts)-1].card.Number { //确定y不为小王/大王
				case card.JOKER, card.KING:
					return no_match
				default:
					switch counts[(len(counts)-1)-1].card.Number { //确定x不为小王/大王
					case card.JOKER, card.KING:
						return no_match
					default:
						return fd
					}
				}
			case counts[0].count == 1 && counts[len(counts)-1].count == 1: //复杂
				switch counts[0].card.Number { //确定x不为小王/大王
				case card.JOKER, card.KING:
					return no_match
				default:
					switch counts[len(counts)-1].card.Number { //确定y不为小王/大王
					case card.JOKER, card.KING:
						return no_match
					default:
						return complex
					}
				}
			case counts[0].count == 1 && counts[1].count == 1: //反向
				switch counts[0].card.Number { //确定y不为小王/大王
				case card.JOKER, card.KING:
					return no_match
				default:
					switch counts[1].card.Number { //确定x不为小王/大王
					case card.JOKER, card.KING:
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
				if counts[n].card.Number == card.TWO { //确定k(n) != 2
					return false
				}
				if n != len(counts)-1 {
					if counts[n].card.Number+1 != counts[n+1].card.Number { //确定k(n)+1 == k(n+1)
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
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE_WITH_TWO_SINGLE)
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
				if counts[n].card.Number == card.TWO { //确定k(n) != 2
					return false
				}
				if n != len(counts)-1 {
					if counts[n].card.Number+1 != counts[n+1].card.Number { //确定k(n)+1 == k(n+1)
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
			matchs = append(matchs, twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE_WITH_TWO_DOUBLE)
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
		return twoonone_pb.CardType_CARD_TYPE_UNKNOWN
	}
}

// 计算指定牌类型的连续次数
func (r *Room) calcCardContinous(cards []card.Card, cardtype twoonone_pb.CardType) int {
	switch cardtype {
	default:
		return 0
	case twoonone_pb.CardType_CARD_TYPE_SINGLE_SEQUENCE:
		return len(cards)
	case twoonone_pb.CardType_CARD_TYPE_DOUBLE_SEQUENCE:
		return len(cards) / 2
	case twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE:
		return len(cards) / 3
	case twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE_WITH_TWO_SINGLE:
		return (len(cards) - 2) / 3
	case twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE_WITH_TWO_DOUBLE:
		return (len(cards) - 2*2) / 3
	}
}

// 计算指定牌类型的大小
func (r *Room) calcCardSize(cards []card.Card, cardtype twoonone_pb.CardType) int {
	//升序
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Number < cards[j].Number
	})
	//计数
	counts := r.countCard(cards)
	var cardsize int
	switch cardtype {
	case twoonone_pb.CardType_CARD_TYPE_THREE_WITH_SINGLE:
		//上正向匹配(aaab)，下反向匹配(bbba)
		switch {
		case counts[0].count == 3 && counts[1].count == 1:
			cardsize += int(counts[0].card.Number) * counts[0].count
		case counts[0].count == 1 && counts[1].count == 3:
			cardsize += int(counts[1].card.Number) * counts[1].count
		}
	case twoonone_pb.CardType_CARD_TYPE_THREE_WITH_DOUBLE:
		//上正向匹配(aaabb)，下反向匹配(bbbaa)
		switch {
		case counts[0].count == 3 && counts[1].count == 2:
			cardsize += counts[0].count * int(counts[0].card.Number)
		case counts[0].count == 2 && counts[1].count == 3:
			cardsize += counts[1].count * int(counts[1].card.Number)
		}
	case twoonone_pb.CardType_CARD_TYPE_FOUR_WITH_TWO_SINGLE:
		//上正向匹配(aaaabc)，中复杂匹配(abbbbc)，下反向匹配(abcccc)
		switch {
		case counts[0].count == 4 && counts[1].count == 1 && counts[2].count == 1:
			cardsize += counts[0].count * int(counts[0].card.Number)
		case counts[0].count == 1 && counts[1].count == 4 && counts[2].count == 1:
			cardsize += counts[1].count * int(counts[1].card.Number)
		case counts[0].count == 1 && counts[1].count == 1 && counts[2].count == 4:
			cardsize += counts[2].count * int(counts[2].card.Number)
		}
	case twoonone_pb.CardType_CARD_TYPE_FOUR_WITH_TWO_DOUBLE:
		//前正向匹配(aaaabbcc)，中复杂匹配(aabbbbcc)，后反向匹配(aabbcccc)
		switch {
		case counts[0].count == 4 && counts[1].count == 2 && counts[2].count == 2:
			cardsize += counts[0].count * int(counts[0].card.Number)
		case counts[0].count == 2 && counts[1].count == 4 && counts[2].count == 2:
			cardsize += counts[1].count * int(counts[1].card.Number)
		case counts[0].count == 2 && counts[1].count == 2 && counts[2].count == 4:
			cardsize += counts[2].count * int(counts[2].card.Number)
		}
	case twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE_WITH_TWO_SINGLE:
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
				size += v.count * int(v.card.Number)
			}
			return size
		}
		size := calculate(counts, match_type(counts))
		if size != -1 {
			cardsize += size
		}
	case twoonone_pb.CardType_CARD_TYPE_AIR_SEQUENCE_WITH_TWO_DOUBLE:
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
				size += v.count * int(v.card.Number)
			}
			return size
		}
		size := calculate(counts, match_type(counts))
		if size != -1 {
			cardsize += size
		}
	default:
		for _, v := range cards {
			cardsize += int(v.Number)
		}
	}
	return cardsize
}

type cardCount struct {
	card  card.Card
	count int
}

// 牌计数器
func (r *Room) countCard(cards []card.Card) []*cardCount {
	var (
		counter [15]int //0-14代表对应的card
		counts  []*cardCount
	)
	for _, card := range cards {
		counter[card.Number]++
	}
	for c, count := range counter {
		//跳过不存在的牌
		if count == 0 {
			continue
		}
		counts = append(counts, &cardCount{
			card:  card.Card{Number: card.CardNumber(c)},
			count: count,
		})
	}
	return counts
}

func (r *Room) nextRobLandownerOperator() *player.Player {
	for _, v := range r.players {
		if v.GetRobLandownerAction() == twoonone_pb.RobLandownerInfo_ACTION_EMPTY {
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
	r.stage = twoonone_pb.RoomStage_ROOM_STAGE_SENDING_CARD
	// 重排玩家
	r.resequencePlayerSinceSendingCard(&r.players)
	r.operatorNow = lo
}

func (r *Room) getTakeRobLandowners() []*player.Player {
	pis := []*player.Player{}
	for _, v := range r.players {
		if v.GetRobLandownerAction() != twoonone_pb.RobLandownerInfo_ACTION_EMPTY {
			pis = append(pis, v)
		}
	}
	return pis
}

func (r *Room) getRobLandowners() []*player.Player {
	pis := []*player.Player{}
	for _, v := range r.players {
		if v.GetRobLandownerAction() == twoonone_pb.RobLandownerInfo_ACTION_ROB {
			pis = append(pis, v)
		}
	}
	return pis
}

func generateCards() [54]card.Card {
	cards := []card.Card{}
	for n := 0; n < int(card.KING)-2+1; n++ {
		for nn := 0; nn < 4; nn++ {
			cards = append(cards, card.Card{Number: card.CardNumber(n)})
		}
	}
	cards = append(cards, card.Card{Number: card.JOKER})
	cards = append(cards, card.Card{Number: card.KING})
	ret := [54]card.Card{}
	copy(ret[:], cards)
	return ret
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

func (r *Room) resequencePlayerSinceSendingCard(xp *[]*player.Player) {
	x := *xp
	lo := r.GetLandowner()
	if lo != nil { //游戏已开始，重排顺序
		x[0] = lo
		x[1] = r.farmers[0]
		x[2] = r.farmers[1]
	}
}

func (r *Room) resequenceCards(xp *[]card.Card) {
	x := *xp
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
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
	if len(r.sendCards) == 0 {
		return nil
	}
	return r.sendCards[len(r.sendCards)-1]
}

func (r *Room) GetSendCards() []*SendCard {
	return r.sendCards
}

func (r *Room) GetBasicCoin() float64 {
	return r.basicCoin
}

func (r *Room) GetMultiple() int {
	return r.multiple
}

func (r *Room) GetLandownerCard() [3]card.Card {
	return r.landownerCards
}

func FormatInternalRoom2Protobuf(x *Room) *twoonone_pb.RoomInfo {
	if x == nil {
		return nil
	}
	return &twoonone_pb.RoomInfo{
		Hash:        x.hash,
		PlayerInfos: player.FormatInternalPlayers2Protobuf(x.players),
		BasicCoin:   x.basicCoin,
		Multiple:    int32(x.multiple),
		Stage:       x.stage,
		Sendcards:   formatInternalSendCards2Protobuf(x.sendCards),
		OperatorNow: player.FormatInternalPlayer2Protobuf(x.operatorNow),
	}
}

func formatInternalSendCard2Protobuf(x *SendCard) *twoonone_pb.SendCard {
	if x == nil {
		return nil
	}
	return &twoonone_pb.SendCard{
		SenderInfo:         player.FormatInternalPlayer2Protobuf(x.SenderInfo),
		Sendcards:          card.FormatInternalCards2Protobuf(x.SendCards),
		SendcardType:       x.SendCardType,
		SendcardSize:       int32(x.SendCardSize),
		SendcardContinuous: int32(x.SendCardContinous),
	}
}

func formatInternalSendCards2Protobuf(xs []*SendCard) []*twoonone_pb.SendCard {
	if xs == nil {
		return nil
	}
	ret := []*twoonone_pb.SendCard{}
	for _, v := range xs {
		ret = append(ret, formatInternalSendCard2Protobuf(v))
	}
	return ret
}

func FormatInternalRooms2Protobuf(xs []*Room) []*twoonone_pb.RoomInfo {
	if xs == nil {
		return nil
	}
	ret := []*twoonone_pb.RoomInfo{}
	for _, v := range xs {
		ret = append(ret, FormatInternalRoom2Protobuf(v))
	}
	return ret
}
