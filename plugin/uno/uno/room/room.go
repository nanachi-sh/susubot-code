package room

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"

	uno_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
	"github.com/nanachi-sh/susubot-code/plugin/uno/uno/player"
	"github.com/twmb/murmur3"
)

type Room struct {
	hash        string
	players     []*player.Player
	operatorNow *player.Player
	stage       uno_pb.Stage
	cardPool    []*SendCard
	banker      *player.Player
	cardHeap    []uno_pb.Card
}

type wildDrawFourStatus int

const (
	wildDrawFourStatus_challengerLose wildDrawFourStatus = iota
	wildDrawFourStatus_challengedLose
)

type SendCard struct {
	SenderId           string
	SendCard           uno_pb.Card
	wildDrawFourStatus *wildDrawFourStatus //若为Wild draw four时有效
}

func hash() string {
	buf := make([]byte, 100)
	for n := 0; n != len(buf); n++ {
		buf[n] = byte(rand.Intn(256))
	}
	h1, h2 := murmur3.SeedSum128(rand.Uint64(), rand.Uint64(), buf)
	return fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}

func New() *Room {
	return &Room{
		hash:        hash(),
		players:     []*player.Player{},
		operatorNow: nil,
		stage:       uno_pb.Stage_WaitingStart,
		cardPool:    []*SendCard{},
		banker:      nil,
	}
}

func (r *Room) Join(ai *uno_pb.PlayerAccountInfo) *uno_pb.Errors {
	if r.stage != uno_pb.Stage_WaitingStart {
		return uno_pb.Errors_RoomStarted.Enum()
	}
	if len(r.GetPlayers()) > 10 {
		return uno_pb.Errors_RoomFull.Enum()
	}
	if _, ok := r.GetPlayer(ai.Id); ok {
		return uno_pb.Errors_RoomExistPlayer.Enum()
	}
	p := player.New(&uno_pb.PlayerInfo{
		PlayerAccountInfo: ai,
		PlayerRoomInfo:    &uno_pb.PlayerRoomInfo{},
	})
	if p == nil {
		return uno_pb.Errors_Unexpected.Enum()
	}
	r.players = append(r.players, p)
	return nil
}

func (r *Room) Exit(playerid string) *uno_pb.Errors {
	if r.GetStage() != uno_pb.Stage_WaitingStart {
		return uno_pb.Errors_RoomStarted.Enum()
	}
	if _, ok := r.GetPlayer(playerid); !ok {
		return uno_pb.Errors_RoomNoExistPlayer.Enum()
	}
	if !r.deletePlayer(playerid) {
		return uno_pb.Errors_RoomNoExistPlayer.Enum()
	}
	return nil
}

func (r *Room) Start() *uno_pb.Errors {
	if r.GetStage() != uno_pb.Stage_WaitingStart {
		return uno_pb.Errors_RoomStarted.Enum()
	}
	if len(r.GetPlayers()) < 2 {
		return uno_pb.Errors_RoomNoReachPlayers.Enum()
	}
	// 进入下一阶段
	r.startElectBanker()
	return nil
}

func (r *Room) DrawCard(p *player.Player) *uno_pb.Errors {
	// 仅分流
	switch r.stage {
	case uno_pb.Stage_WaitingStart:
		return uno_pb.Errors_RoomNoStart.Enum()
	case uno_pb.Stage_ElectingBanker:
		return r.drawCard_ElectingBanker(p)
	case uno_pb.Stage_SendingCard:
		return r.drawCard_SendingCard(p)
	default:
		return uno_pb.Errors_Unexpected.Enum()
	}
}

type SendCardEvent struct {
	GameFinish bool

	GameFinishE *uno_pb.SendCardActionResponse_GameFinishEvent
}

func (r *Room) SendCardAction(p *player.Player, sendcard uno_pb.Card, action uno_pb.SendCardActions) (*player.Player, *SendCardEvent, *uno_pb.Errors) {
	switch action {
	case uno_pb.SendCardActions_Send:
		return r.sendCard(p, sendcard)
	case uno_pb.SendCardActions_NoSend:
		return r.noSendCard(p)
	default:
		return nil, nil, uno_pb.Errors_Unexpected.Enum()
	}
}

func (r *Room) sendCard(p *player.Player, sendcard uno_pb.Card) (*player.Player, *SendCardEvent, *uno_pb.Errors) {
	sendcardFC := sendcard.FeatureCard
	last := r.GetLastCard()
	// 特殊情况判断
	if last == nil { //第一次出牌
		if serr := r.playerSendCard(p, sendcard); serr != nil {
			return nil, nil, serr
		}
		if sendcardFC != nil {
			r.sendCard_featureCardAction(sendcardFC.FeatureCard)
		}
		next := r.nextOperator()
		r.operatorNow = next
		return next, nil, nil
	} else if p.GetDrawCard() != nil { //出摸来的牌
		if !r.sendCard_cardCheck(last.SendCard, sendcard) {
			return nil, nil, uno_pb.Errors_SendCardColorORNumberNELastCard.Enum()
		}
		if !p.DeleteCardFromDrawCard(sendcard) {
			return nil, nil, uno_pb.Errors_PlayerCardNoExist.Enum()
		}
		if sendcardFC != nil {
			r.convertBlackCardColor(last.SendCard, sendcard)
			r.sendCard_featureCardAction(sendcardFC.FeatureCard)
		}
		r.addCardToCardPool(SendCard{
			SenderId: p.GetId(),
			SendCard: sendcard,
		})
		if len(p.GetCards()) == 0 {
			ps := []*uno_pb.PlayerInfo{}
			for _, v := range r.players {
				ps = append(ps, v.FormatToProtoBuf())
			}
			return nil, &SendCardEvent{
				GameFinish: true,
				GameFinishE: &uno_pb.SendCardActionResponse_GameFinishEvent{
					Players: ps,
					Winner:  p.FormatToProtoBuf(),
				},
			}, nil
		}
		next := r.nextOperator()
		r.operatorNow = next
		return next, nil, nil
	} else if sendcardFC != nil && (sendcardFC.FeatureCard == uno_pb.FeatureCards_Skip || sendcardFC.FeatureCard == uno_pb.FeatureCards_Reverse) { //Skip, Reverse无视上一张牌类型
		switch {
		case last.SendCard.FeatureCard != nil:
			lastFC := last.SendCard.FeatureCard
			if lastFC.Color != sendcardFC.Color {
				return nil, nil, uno_pb.Errors_SendCardColorORNumberNELastCard.Enum()
			}
		case last.SendCard.NormalCard != nil:
			lastNC := last.SendCard.NormalCard
			if lastNC.Color != sendcardFC.Color {
				return nil, nil, uno_pb.Errors_SendCardColorORNumberNELastCard.Enum()
			}
		}
		if serr := r.playerSendCard(p, sendcard); serr != nil {
			return nil, nil, serr
		}
		r.sendCard_featureCardAction(sendcardFC.FeatureCard)
		next := r.nextOperator()
		r.operatorNow = next
		return next, nil, nil
	} else if sendcardFC != nil && last.SendCard.FeatureCard != nil && ((sendcardFC.FeatureCard == uno_pb.FeatureCards_DrawTwo || sendcardFC.FeatureCard == uno_pb.FeatureCards_WildDrawFour) && (last.SendCard.FeatureCard.FeatureCard == uno_pb.FeatureCards_DrawTwo || last.SendCard.FeatureCard.FeatureCard == uno_pb.FeatureCards_WildDrawFour)) { //上一张为Draw two/Wild Draw four，且这一张也同样
		if !r.sendCard_cardCheck(last.SendCard, sendcard) {
			return nil, nil, uno_pb.Errors_SendCardColorORNumberNELastCard.Enum()
		}
		if serr := r.playerSendCard(p, sendcard); serr != nil {
			return nil, nil, serr
		}
		r.sendCard_featureCardAction(sendcardFC.FeatureCard)
		next := r.nextOperator()
		r.operatorNow = next
		return next, nil, nil
	} else if last.SendCard.FeatureCard != nil && (last.SendCard.FeatureCard.FeatureCard == uno_pb.FeatureCards_Skip || last.SendCard.FeatureCard.FeatureCard == uno_pb.FeatureCards_DrawTwo || (last.SendCard.FeatureCard.FeatureCard == uno_pb.FeatureCards_WildDrawFour && (last.wildDrawFourStatus != nil && *last.wildDrawFourStatus != wildDrawFourStatus_challengerLose))) { //上一张牌为跳过牌，且这次出的不为跳过牌，则不允许出牌
		return nil, nil, uno_pb.Errors_PlayerCannotSendCard.Enum()
	}
	// 正常出牌
	if !r.sendCard_cardCheck(last.SendCard, sendcard) {
		return nil, nil, uno_pb.Errors_SendCardColorORNumberNELastCard.Enum()
	}
	if serr := r.playerSendCard(p, sendcard); serr != nil {
		return nil, nil, serr
	}
	if sendcardFC != nil {
		r.sendCard_featureCardAction(sendcardFC.FeatureCard)
	}
	next := r.nextOperator()
	r.operatorNow = next
	return next, nil, nil
}

func (r *Room) sendCard_cardCheck(last, now uno_pb.Card) bool {
	lastNC := last.NormalCard
	lastFC := last.FeatureCard
	nowNC := now.NormalCard
	nowFC := now.FeatureCard
	switch now.Type {
	case uno_pb.CardType_Normal:
		switch last.Type {
		case uno_pb.CardType_Normal:
			if nowNC.Color != lastNC.Color && nowNC.Number != lastNC.Number {
				return false
			}
		case uno_pb.CardType_Feature:
			if nowNC.Color != lastFC.Color {
				return false
			}
		}
	case uno_pb.CardType_Feature:
		switch last.Type {
		case uno_pb.CardType_Normal:
			if nowFC.Color != lastNC.Color {
				return false
			}
		case uno_pb.CardType_Feature:
			if nowFC.Color != lastFC.Color && nowFC.FeatureCard != lastFC.FeatureCard {
				return false
			}
		}
	}
	return true
}

func (r *Room) noSendCard(p *player.Player) (*player.Player, *SendCardEvent, *uno_pb.Errors) {
	_, _, ok := r.getStackFeatureCard()
	if ok {
		return nil, nil, uno_pb.Errors_PlayerCardNoExist.Enum()
	} else if p.GetDrawCard() != nil {
		p.AddCards([]uno_pb.Card{*p.GetDrawCard()})
		p.ClearDrawCard()
	} else if p.GetDrawCard() == nil {
		return nil, nil, uno_pb.Errors_PlayerNoDrawCard.Enum()
	}
	next := r.nextOperator()
	r.operatorNow = next
	return next, nil, nil
}

func (r *Room) sendCard_featureCardAction(featureCard uno_pb.FeatureCards) {
	switch featureCard {
	case uno_pb.FeatureCards_Skip:
		// 跳过下一个玩家
		skipP := r.nextOperator()
		r.operatorNow = skipP
	case uno_pb.FeatureCards_Reverse:
		r.reverseSequence()
	case uno_pb.FeatureCards_DrawTwo:
	case uno_pb.FeatureCards_WildDrawFour:
	}
}

func (r *Room) nextOperator() *player.Player {
	ps := r.GetPlayers()
	for i, v := range ps {
		if r.GetOperatorNow().GetId() == v.GetId() {
			if i == len(ps)-1 {
				return ps[0]
			} else {
				return ps[i+1]
			}
		}
	}
	return nil
}

func (r *Room) reverseSequence() {
	playersCopy := make([]*player.Player, len(r.players))
	copy(playersCopy, r.players)
	for n, N := len(playersCopy), 0; ; {
		if n == len(r.players)-1 || N < 0 {
			break
		}
		r.players[N] = playersCopy[n]
	}
}

func (r *Room) convertBlackCardColor(last uno_pb.Card, now uno_pb.Card) {
	if now.Type == uno_pb.CardType_Feature {
		switch now.FeatureCard.FeatureCard {
		case uno_pb.FeatureCards_Wild, uno_pb.FeatureCards_WildDrawFour:
			last := r.GetLastCard()
			if last != nil {
				var clr uno_pb.CardColor
				switch last.SendCard.Type {
				case uno_pb.CardType_Normal:
					clr = last.SendCard.NormalCard.Color
				case uno_pb.CardType_Feature:
					clr = last.SendCard.FeatureCard.Color
				}
				now.FeatureCard.Color = clr
			}
		}
	}
}

func (r *Room) playerSendCard(p *player.Player, sendcard uno_pb.Card) *uno_pb.Errors {
	if last := r.GetLastCard(); last != nil {
		r.convertBlackCardColor(last.SendCard, sendcard)
	}
	if !p.DeleteCardFromHandCard(sendcard) {
		return uno_pb.Errors_PlayerCardNoExist.Enum()
	}
	r.addCardToCardPool(SendCard{
		SenderId: p.GetId(),
		SendCard: sendcard,
	})
	return nil
}

func (r *Room) startSendCard() {
	banker := r.electBankerCardMaxIs()
	r.banker = banker
	r.operatorNow = banker
	// 将牌丢回卡堆，并重新洗牌
	for _, v := range r.players {
		r.addCardsToCardHeap([]uno_pb.Card{*v.GetElectBankerCard()})
	}
	r.resequence(&r.cardHeap)
	// 为玩家发牌
	for _, v := range r.players {
		cards := r.cutCards(7)
		v.AddCards(cards)
	}
	r.stage = uno_pb.Stage_SendingCard
	// 重排顺序
	r.resequence(r.players)
}

func (r *Room) drawCard_ElectingBanker(p *player.Player) *uno_pb.Errors {
	for {
		card := r.cutCards(1)[0]
		if card.NormalCard == nil {
			r.addCardsToCardHeap([]uno_pb.Card{card})
			continue
		}
		p.SetElectBankerCard(card)
		break
	}
	if len(r.getTakeElectBankerPlayers()) == len(r.players) {
		r.startSendCard()
	}
	return nil
}

func (r *Room) drawCard_SendingCard(p *player.Player) *uno_pb.Errors {
	// 可以抽牌的情况：
	// 1.遭到Wild draw four, Draw two(可以打出Skip或Wild draw four跳到下一个玩家，相关见SendCard)
	// 2.轮到该玩家出牌，但不想出或无牌可出
	if r.operatorNow.GetId() != p.GetId() {
		return uno_pb.Errors_PlayerNoOperatorNow.Enum()
	}
	//
	stackFC, count, ok := r.getStackFeatureCard()
	if ok && stackFC == uno_pb.FeatureCards_DrawTwo { //遭到Draw two
		// 摸两张牌，并跳过回合
		cards := r.cutCards(2 * count)
		p.AddCards(cards)
		p.SetCallUNO(false)
		next := r.nextOperator()
		r.operatorNow = next
		return nil
	} else if ok && stackFC == uno_pb.FeatureCards_Wild { //遭到Wild draw four
		last := r.GetLastCard()
		if last.wildDrawFourStatus == nil { //未挑战
			cards := r.cutCards(4 * count)
			p.AddCards(cards)
			next := r.nextOperator()
			p.SetCallUNO(false)
			r.operatorNow = next
			return nil
		} else { //已挑战
			switch *last.wildDrawFourStatus {
			default:
				return uno_pb.Errors_Unexpected.Enum()
			case wildDrawFourStatus_challengedLose: //被挑战者失败
				//被挑战者失败不应在此处理
				return uno_pb.Errors_Unexpected.Enum()
			case wildDrawFourStatus_challengerLose: //挑战者失败
				cards := r.cutCards(4*count + 2)
				p.AddCards(cards)
				p.SetCallUNO(false)
				next := r.nextOperator()
				r.operatorNow = next
				return nil
			}
		}
	} else if p.GetDrawCard() != nil { //已抽过一张牌
		return uno_pb.Errors_PlayerAlreadyDrawCard.Enum()
	}
	// 玩家回合抽牌
	card := r.cutCards(1)[0]
	p.SetCallUNO(false)
	p.SetDrawCard(card)
	return nil
}

func (r *Room) getStackFeatureCard() (uno_pb.FeatureCards, int, bool) {
	last := r.GetLastCard()
	if last == nil {
		return 0, 0, false
	}
	var (
		ct    *uno_pb.FeatureCards
		count int
	)
	for n := len(r.cardPool); n > 0; n-- {
		sc := r.cardPool[n]
		if sc.SendCard.Type != uno_pb.CardType_Feature {
			break
		}
		fC := sc.SendCard.FeatureCard
		switch fC.FeatureCard {
		case uno_pb.FeatureCards_DrawTwo, uno_pb.FeatureCards_WildDrawFour:
		default:
			break
		}
		if ct == nil {
			ct = &fC.FeatureCard
			count++
		} else if *ct == fC.FeatureCard {
			count++
		} else {
			break
		}
	}
	if ct == nil {
		return 0, 0, false
	} else {
		return *ct, count, true
	}
}

func (r *Room) electBankerCardMaxIs() *player.Player {
	ps := r.getTakeElectBankerPlayers()
	if len(ps) == 0 {
		return nil
	}
	psCopy := make([]*player.Player, len(ps))
	copy(psCopy, ps)
	sort.Slice(psCopy, func(i, j int) bool {
		return psCopy[i].GetElectBankerCard().NormalCard.Number > psCopy[j].GetElectBankerCard().NormalCard.Number
	})
	return psCopy[0]
}

func (r *Room) getTakeElectBankerPlayers() []*player.Player {
	ps := []*player.Player{}
	for _, v := range r.GetPlayers() {
		if v.GetElectBankerCard() != nil {
			ps = append(ps, v)
		}
	}
	return ps
}

func (r *Room) startElectBanker() {
	cardsArr := r.generateCards()
	cards := cardsArr[:]
	r.resequence(&cards)
	copy(r.cardHeap, cards)
	r.stage = uno_pb.Stage_ElectingBanker
}

func (r *Room) CallUNO(p *player.Player) ([]uno_pb.Card, *uno_pb.Errors) {
	if r.stage != uno_pb.Stage_SendingCard {
		return nil, uno_pb.Errors_RoomNoSendingCard.Enum()
	}
	if r.operatorNow.GetId() != p.GetId() {
		return nil, uno_pb.Errors_PlayerNoOperatorNow.Enum()
	}
	if len(p.GetCards()) > 2 {
		cards := r.cutCards(2)
		p.AddCards(cards)
		return p.GetCards(), uno_pb.Errors_PlayerCannotCallUNO.Enum()
	}
	p.SetCallUNO(true)
	return nil, nil
}

func (r *Room) Challenge(p *player.Player) (bool, []uno_pb.Card, *uno_pb.Errors) {
	if r.stage != uno_pb.Stage_SendingCard {
		return false, nil, uno_pb.Errors_RoomNoSendingCard.Enum()
	}
	if r.operatorNow.GetId() != p.GetId() {
		return false, nil, uno_pb.Errors_PlayerNoOperatorNow.Enum()
	}
	last := r.GetLastCard()
	if last == nil {
		return false, nil, uno_pb.Errors_RoomNoneSendCard.Enum()
	}
	if last.SendCard.Type != uno_pb.CardType_Feature {
		return false, nil, uno_pb.Errors_CannotChallenge.Enum()
	}
	if last.SendCard.FeatureCard.FeatureCard != uno_pb.FeatureCards_WildDrawFour {
		return false, nil, uno_pb.Errors_CannotChallenge.Enum()
	}
	clr := last.SendCard.FeatureCard.Color
	win := false
	for _, v := range p.GetCards() {
		switch v.Type {
		case uno_pb.CardType_Normal:
			if v.NormalCard.Color == clr {
				win = true
				break
			}
		case uno_pb.CardType_Feature:
			if v.FeatureCard.Color == clr {
				win = true
				break
			}
		}
	}
	if win {
		last.wildDrawFourStatus = new(wildDrawFourStatus)
		*last.wildDrawFourStatus = wildDrawFourStatus_challengedLose
		cards := r.cutCards(4)
		lastP, ok := r.GetPlayer(last.SenderId)
		if !ok {
			return false, nil, uno_pb.Errors_Unexpected.Enum()
		}
		lastP.AddCards(cards)
		return true, lastP.GetCards(), nil
	} else {
		last.wildDrawFourStatus = new(wildDrawFourStatus)
		*last.wildDrawFourStatus = wildDrawFourStatus_challengerLose
		return false, nil, nil
	}
}

func (r *Room) IndicateUNO(tP *player.Player) ([]uno_pb.Card, *uno_pb.Errors) {
	if r.stage != uno_pb.Stage_SendingCard {
		return nil, uno_pb.Errors_RoomNoSendingCard.Enum()
	}
	if len(tP.GetCards()) < 2 {
		cards := r.cutCards(2)
		tP.AddCards(cards)
		return tP.GetCards(), uno_pb.Errors_PlayerCannotCallUNO.Enum()
	} else {
		return nil, uno_pb.Errors_PlayerAlreadyCallUNO.Enum()
	}
}

func (r *Room) generateCards() [108]uno_pb.Card {
	// 普通牌：
	//
	cards := []uno_pb.Card{}
	//添加普通牌，0四色牌各一张，1-9四色牌各两张
	for _, n := range uno_pb.CardNumber_value {
		need := 0
		if n == int32(uno_pb.CardNumber_Zero) {
			need = 1
		} else {
			need = 2
		}
		for i := 0; i < need; i++ {
			for _, color := range uno_pb.CardColor_value {
				if color == int32(uno_pb.CardColor_Black) {
					continue
				}
				cards = append(cards, uno_pb.Card{
					NormalCard: &uno_pb.NormalCard{
						Color:  uno_pb.CardColor(color),
						Number: uno_pb.CardNumber(n),
					},
					Type: uno_pb.CardType_Normal,
				})
			}
		}
	}
	// 添加功能牌，Skip, Reverse, Draw two四色各两张，Wild, Wild draw four各四张黑色
	for _, fea := range uno_pb.FeatureCards_value {
		need := 0
		switch uno_pb.FeatureCards(fea) {
		case uno_pb.FeatureCards_Skip, uno_pb.FeatureCards_Reverse, uno_pb.FeatureCards_DrawTwo:
			need = 2
		case uno_pb.FeatureCards_Wild, uno_pb.FeatureCards_WildDrawFour:
			need = 4
		}
		for i := 0; i < need; i++ {
			switch uno_pb.FeatureCards(fea) {
			case uno_pb.FeatureCards_Skip, uno_pb.FeatureCards_Reverse, uno_pb.FeatureCards_DrawTwo:
				for _, color := range uno_pb.CardColor_value {
					if color == int32(uno_pb.CardColor_Black) {
						continue
					}
					cards = append(cards, uno_pb.Card{
						FeatureCard: &uno_pb.FeatureCard{
							Color:       uno_pb.CardColor(color),
							FeatureCard: uno_pb.FeatureCards(fea),
						},
						Type: uno_pb.CardType_Feature,
					})
				}
			case uno_pb.FeatureCards_Wild, uno_pb.FeatureCards_WildDrawFour:
				cards = append(cards, uno_pb.Card{
					FeatureCard: &uno_pb.FeatureCard{
						Color:       uno_pb.CardColor_Black,
						FeatureCard: uno_pb.FeatureCards(fea),
					},
					Type: uno_pb.CardType_Feature,
				})
			}
		}
	}
	return [108]uno_pb.Card(cards)
}

func (r *Room) resequence(x any) {
	switch xPoint := x.(type) {
	case *[]uno_pb.Card:
		x := *xPoint
		rand.Shuffle(len(x), func(i, j int) {
			x[i], x[j] = x[j], x[i]
		})
	case *[]*player.Player:
		if r.stage == uno_pb.Stage_SendingCard {
			var xCopy []*player.Player
			xCopy = append(xCopy, r.banker)
			for _, v := range xCopy {
				if v.GetId() == r.banker.GetId() {
					continue
				}
				xCopy = append(xCopy, v)
			}
			*xPoint = xCopy
		}
	}
}

func (r *Room) deletePlayer(playerid string) bool {
	for i, v := range r.GetPlayers() {
		if v.GetId() == playerid {
			if len(r.players) == 1 {
				r.players = r.players[0:]
			} else {
				r.players = append(r.players[:i], r.players[i+1:]...)
			}
			return true
		}
	}
	return false
}

func (r *Room) GetOperatorNow() *player.Player {
	return r.operatorNow
}

func (r *Room) GetStage() uno_pb.Stage {
	return r.stage
}

func (r *Room) GetPlayers() []*player.Player {
	return r.players
}

func (r *Room) GetPlayer(id string) (*player.Player, bool) {
	for _, v := range r.GetPlayers() {
		if v.GetId() == id {
			return v, true
		}
	}
	return nil, false
}

func (r *Room) GetBanker() *player.Player {
	return r.banker
}

func (r *Room) GetLastCard() *SendCard {
	if len(r.cardPool) == 0 {
		return nil
	}
	return r.cardPool[len(r.cardPool)-1]
}

func (r *Room) GetHash() string {
	return r.hash
}

func (r *Room) cutCards(count int) []uno_pb.Card {
	if len(r.cardHeap) < count {
		//牌堆牌数不足，重开一副新牌，并重洗牌
		cardsArr := r.generateCards()
		r.addCardsToCardHeap(cardsArr[:])
	}
	cards := r.cardHeap[:count]
	r.cardHeap = r.cardHeap[count:]
	return cards
}

func (r *Room) addCardsToCardHeap(cards []uno_pb.Card) {
	r.cardHeap = append(r.cardHeap, cards...)
}

func (r *Room) addCardToCardPool(sendcard SendCard) {
	r.cardPool = append(r.cardPool, &sendcard)
}

func (r *Room) FormatToProtoBuf() *uno_pb.Room {
	ch := []*uno_pb.Card{}
	cp := []*uno_pb.SendCard{}
	ps := []*uno_pb.PlayerInfo{}
	for _, v := range r.cardHeap {
		ch = append(ch, &v)
	}
	for _, v := range r.cardPool {
		var wdfstate *uno_pb.WildDrawFourStatus
		if v.wildDrawFourStatus != nil {
			switch *v.wildDrawFourStatus {
			case wildDrawFourStatus_challengedLose:
				wdfstate = uno_pb.WildDrawFourStatus_ChallengedLose.Enum()
			case wildDrawFourStatus_challengerLose:
				wdfstate = uno_pb.WildDrawFourStatus_ChallengerLose.Enum()
			}
		}
		cp = append(cp, &uno_pb.SendCard{
			SenderId:           v.SenderId,
			SendCard:           &v.SendCard,
			WildDrawFourStatus: wdfstate,
		})
	}
	for _, v := range r.players {
		ps = append(ps, v.FormatToProtoBuf())
	}
	ur := &uno_pb.Room{
		Hash:        r.hash,
		Stage:       r.stage,
		Banker:      nil,
		CardHeap:    ch,
		CardPool:    cp,
		OperatorNow: nil,
		Players:     ps,
	}
	if r.banker != nil {
		ur.Banker = r.banker.FormatToProtoBuf()
	}
	if r.operatorNow != nil {
		ur.OperatorNow = r.operatorNow.FormatToProtoBuf()
	}
	return ur
}
