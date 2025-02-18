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
	hash              string
	players           []*player.Player
	operatorNow       *player.Player
	stage             uno_pb.Stage
	cardPool          []*SendCard
	banker            *player.Player
	cardHeap          []uno_pb.Card
	sequenceDirection direction
	sequencePosition  int
	seeds             [2]uint64
}

type direction int

const (
	fd direction = iota
	rd
)

type wildDrawFourStatus int

const (
	wildDrawFourStatus_challengerLose wildDrawFourStatus = iota
	wildDrawFourStatus_challengedLose
)

type SendCard struct {
	SenderId           string
	SendCard           uno_pb.Card
	wildDrawFourStatus *wildDrawFourStatus //若为Wild draw four时有效
	featureEffected    bool
}

func hash() ([2]uint64, string) {
	buf := make([]byte, 100)
	for n := 0; n != len(buf); n++ {
		buf[n] = byte(rand.Intn(256))
	}
	s1, s2 := rand.Uint64(), rand.Uint64()
	h1, h2 := murmur3.SeedSum128(s1, s2, buf)
	return [2]uint64{s1, s2}, fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}

func New() *Room {
	seeds, hash := hash()
	return &Room{
		hash:              hash,
		seeds:             seeds,
		players:           []*player.Player{},
		operatorNow:       nil,
		stage:             uno_pb.Stage_WaitingStart,
		cardPool:          []*SendCard{},
		banker:            nil,
		cardHeap:          []uno_pb.Card{},
		sequenceDirection: fd,
	}
}

func (r *Room) Join(ai *uno_pb.PlayerAccountInfo, playerHash string) *uno_pb.Errors {
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
		PlayerRoomInfo: &uno_pb.PlayerRoomInfo{
			RoomHash: r.hash,
			Hash:     playerHash,
			Cards:    []*uno_pb.Card{},
		},
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

type DrawCardEvent struct {
	IntoSendCard bool
	Skipped      bool

	IntoSendCardE *uno_pb.RoomEventResponse_DrawCard_IntoSendCardEvent
	SkippedE      *uno_pb.RoomEventResponse_DrawCard_SkippedEvent
}

func (r *Room) DrawCard(p *player.Player) (*DrawCardEvent, *uno_pb.Errors) {
	// 仅分流
	switch r.stage {
	case uno_pb.Stage_WaitingStart:
		return nil, uno_pb.Errors_RoomNoStart.Enum()
	case uno_pb.Stage_ElectingBanker:
		return r.drawCard_ElectingBanker(p)
	case uno_pb.Stage_SendingCard:
		return r.drawCard_SendingCard(p)
	default:
		return nil, uno_pb.Errors_Unexpected.Enum()
	}
}

type SendCardEvent struct {
	GameFinish bool

	GameFinishE *uno_pb.RoomEventResponse_GameFinishEvent
}

func (r *Room) gameFinish() *uno_pb.RoomEventResponse_GameFinishEvent {
	ps := []*uno_pb.PlayerInfo{}
	for _, v := range r.players {
		ps = append(ps, v.FormatToProtoBuf())
	}
	return &uno_pb.RoomEventResponse_GameFinishEvent{
		Players: ps,
		Winner:  r.operatorNow.FormatToProtoBuf(),
	}
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

type sendcard_from int

const (
	handCard sendcard_from = iota
	drawCard
)

func (r *Room) sendCard_checkBlackCard(card uno_pb.Card) bool {
	if card.Type != uno_pb.CardType_Feature {
		return false
	}
	switch card.FeatureCard.FeatureCard {
	case uno_pb.FeatureCards_Wild, uno_pb.FeatureCards_WildDrawFour:
		return true
	}
	return false
}

func (r *Room) sendCard(p *player.Player, sendcard uno_pb.Card) (*player.Player, *SendCardEvent, *uno_pb.Errors) {
	if r.stage != uno_pb.Stage_SendingCard {
		return nil, nil, uno_pb.Errors_RoomNoSendingCard.Enum()
	}
	if r.operatorNow.GetId() != p.GetId() {
		return nil, nil, uno_pb.Errors_PlayerNoOperatorNow.Enum()
	}
	sendcardFC := sendcard.FeatureCard
	last := r.GetLastCard()
	from := handCard
	if p.GetDrawCard() != nil {
		from = drawCard
	}
	// 特殊情况判断
	if r.sendCard_checkSkipCard(last) { //上一张为Skip
		return nil, nil, uno_pb.Errors_PlayerCannotSendCard.Enum()
	} else if r.sendCard_checkSkipORReverseCard(last, sendcard) { //上一张不为Skip，Skip, Reverse可无视牌出
		if !r.sendCard_cardCheck(last.SendCard, sendcard) {
			return nil, nil, uno_pb.Errors_SendCardColorORNumberNELastCard.Enum()
		}
		if serr := r.playerSendCard(p, from, sendcard); serr != nil {
			return nil, nil, serr
		}
		if r.sendCard_checkNeedDrawCard(last) {
			jnlast := r.GetLastCard()
			jnlast.featureEffected = true
		}
		r.sendCard_featureCardAction(sendcardFC.FeatureCard)
		if len(p.GetCards()) == 0 {
			return nil, &SendCardEvent{
				GameFinish:  true,
				GameFinishE: r.gameFinish(),
			}, nil
		}
		next := r.nextOperator()
		r.operatorNow = next
		return next, nil, nil
	} else if r.sendCard_checkStackCard(last.SendCard, sendcard) { //上一张为Draw two/Wild Draw four，且这一张也同样
		if !r.sendCard_cardCheck(last.SendCard, sendcard) {
			return nil, nil, uno_pb.Errors_SendCardColorORNumberNELastCard.Enum()
		}
		if serr := r.playerSendCard(p, from, sendcard); serr != nil {
			return nil, nil, serr
		}
		r.sendCard_featureCardAction(sendcardFC.FeatureCard)
		if len(p.GetCards()) == 0 {
			return nil, &SendCardEvent{
				GameFinish:  true,
				GameFinishE: r.gameFinish(),
			}, nil
		}
		next := r.nextOperator()
		r.operatorNow = next
		return next, nil, nil
	} else if r.sendCard_checkNeedDrawCard(last) { //上一张牌未生效，且这次出的不为跳过牌，并且不为同一人，则不允许出牌
		return nil, nil, uno_pb.Errors_PlayerCannotSendCard.Enum()
	}
	// 正常出牌
	if !r.sendCard_cardCheck(last.SendCard, sendcard) {
		return nil, nil, uno_pb.Errors_SendCardColorORNumberNELastCard.Enum()
	}
	if serr := r.playerSendCard(p, from, sendcard); serr != nil {
		return nil, nil, serr
	}
	if sendcardFC != nil {
		r.sendCard_featureCardAction(sendcardFC.FeatureCard)
	}
	if len(p.GetCards()) == 0 {
		return nil, &SendCardEvent{
			GameFinish:  true,
			GameFinishE: r.gameFinish(),
		}, nil
	}
	next := r.nextOperator()
	r.operatorNow = next
	return next, nil, nil
}

func (r *Room) sendCard_checkSkipORReverseCard(last *SendCard, now uno_pb.Card) bool {
	if len(r.players) == 2 {
		return false
	}
	if now.Type != uno_pb.CardType_Feature {
		return false
	}
	nowFC := now.FeatureCard
	switch nowFC.FeatureCard {
	case uno_pb.FeatureCards_Skip, uno_pb.FeatureCards_Reverse:
		return true
	}
	return false
}

// 检查是否为可堆叠卡
func (r *Room) sendCard_checkStackCard(last, now uno_pb.Card) bool {
	if len(r.players) == 2 {
		return false
	}
	if last.Type != uno_pb.CardType_Feature {
		return false
	}
	if now.Type != uno_pb.CardType_Feature {
		return false
	}
	lastFC := last.FeatureCard
	nowFC := now.FeatureCard
	switch lastFC.FeatureCard {
	case uno_pb.FeatureCards_DrawTwo, uno_pb.FeatureCards_WildDrawFour:
		switch nowFC.FeatureCard {
		case uno_pb.FeatureCards_DrawTwo, uno_pb.FeatureCards_WildDrawFour:
			return true
		}
	}
	return false
}

// 检查是否为需摸牌且还未生效的卡
func (r *Room) sendCard_checkNeedDrawCard(last *SendCard) bool {
	if last == nil {
		return false
	}
	if last.featureEffected {
		return false
	}
	if last.SendCard.Type != uno_pb.CardType_Feature {
		return false
	}
	lastFC := last.SendCard.FeatureCard
	switch lastFC.FeatureCard {
	case uno_pb.FeatureCards_DrawTwo:
		return true
	}
	if lastFC.FeatureCard == uno_pb.FeatureCards_WildDrawFour {
		if last.wildDrawFourStatus == nil {
			return true
		}
		return *last.wildDrawFourStatus == wildDrawFourStatus_challengerLose
	}
	return false
}

// 检查上一张牌是否为有效的Skip
func (r *Room) sendCard_checkSkipCard(last *SendCard) bool {
	if last.SendCard.Type == uno_pb.CardType_Feature && !last.featureEffected {
		if last.SendCard.FeatureCard.FeatureCard == uno_pb.FeatureCards_Skip {
			return true
		}
	}
	return false
}

func (r *Room) sendCard_checkNotSPECBlackCard(now uno_pb.Card) bool {
	if now.Type != uno_pb.CardType_Feature {
		return false
	}
	switch now.FeatureCard.FeatureCard {
	case uno_pb.FeatureCards_Wild, uno_pb.FeatureCards_WildDrawFour:
		return now.FeatureCard.Color == uno_pb.CardColor_Black
	}
	return false
}

func (r *Room) sendCard_cardCheck(last, now uno_pb.Card) bool {
	if r.sendCard_checkBlackCard(now) {
		return true
	}
	lastNC := last.NormalCard
	lastFC := last.FeatureCard
	nowNC := now.NormalCard
	nowFC := now.FeatureCard
	switch now.Type {
	case uno_pb.CardType_Normal:
		switch last.Type {
		case uno_pb.CardType_Normal:
			if nowNC.Color == lastNC.Color || nowNC.Number == lastNC.Number {
				return true
			}
		case uno_pb.CardType_Feature:
			if nowNC.Color == lastFC.Color {
				return true
			}
		}
	case uno_pb.CardType_Feature:
		switch last.Type {
		case uno_pb.CardType_Normal:
			if nowFC.Color == lastNC.Color {
				return true
			}
		case uno_pb.CardType_Feature:
			if nowFC.Color == lastFC.Color || nowFC.FeatureCard == lastFC.FeatureCard {
				return true
			}
		}
	}
	return false
}

func (r *Room) noSendCard(p *player.Player) (*player.Player, *SendCardEvent, *uno_pb.Errors) {
	if r.stage != uno_pb.Stage_SendingCard {
		return nil, nil, uno_pb.Errors_RoomNoSendingCard.Enum()
	}
	if r.operatorNow.GetId() != p.GetId() {
		return nil, nil, uno_pb.Errors_PlayerNoOperatorNow.Enum()
	}
	last := r.GetLastCard()
	if r.sendCard_checkNeedDrawCard(last) {
		return nil, nil, uno_pb.Errors_PlayerCannotNoSendCard.Enum()
	} else if r.sendCard_checkSkipCard(last) {
		last.featureEffected = true
		next := r.nextOperator()
		r.operatorNow = next
		return next, nil, nil
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
	case uno_pb.FeatureCards_Reverse:
		r.reverseSequence()
	case uno_pb.FeatureCards_DrawTwo:
	case uno_pb.FeatureCards_WildDrawFour:
	}
}

func (r *Room) nextOperator() *player.Player {
	switch r.sequenceDirection {
	case fd:
		if r.sequencePosition+1 > len(r.players)-1 {
			r.sequencePosition = 0
		} else {
			r.sequencePosition++
		}
	case rd:
		if r.sequencePosition-1 < 0 {
			r.sequencePosition = len(r.players) - 1
		} else {
			r.sequencePosition--
		}
	}
	return r.players[r.sequencePosition]
}

func (r *Room) reverseSequence() {
	switch r.sequenceDirection {
	case fd:
		r.sequenceDirection = rd
	case rd:
		r.sequenceDirection = fd
	}
}

func (r *Room) playerSendCard(p *player.Player, from sendcard_from, sendcard uno_pb.Card) *uno_pb.Errors {
	if r.sendCard_checkNotSPECBlackCard(sendcard) {
		return uno_pb.Errors_BlackCardNoSpecifiedColor.Enum()
	}
	oColor := uno_pb.CardColor_Black
	if r.sendCard_checkBlackCard(sendcard) {
		switch sendcard.Type {
		case uno_pb.CardType_Normal:
			oColor = sendcard.NormalCard.Color
		case uno_pb.CardType_Feature:
			oColor = sendcard.FeatureCard.Color
		}
		sendcard.FeatureCard.Color = uno_pb.CardColor_Black
	}
	switch from {
	case handCard:
		if !p.DeleteCardFromHandCard(sendcard) {
			return uno_pb.Errors_PlayerCardNoExist.Enum()
		}
	case drawCard:
		if p.CheckCardFromHandCard(sendcard) {
			return uno_pb.Errors_PlayerCannotSendCardFromHandCard.Enum()
		}
		if !p.DeleteCardFromDrawCard(sendcard) {
			return uno_pb.Errors_PlayerCardNoExist.Enum()
		}
	default:
		return uno_pb.Errors_Unexpected.Enum()
	}
	p.SetCallUNO(false)
	if oColor != uno_pb.CardColor_Black {
		sendcard.FeatureCard.Color = oColor
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
	r.cardResequence(r.cardHeap)
	// 为玩家发牌
	for _, v := range r.players {
		cards := r.cutCards(7)
		v.AddCards(cards)
	}
	// 设置引牌
	for {
		card := r.cutCards(1)[0]
		if card.Type == uno_pb.CardType_Normal {
			r.addCardToCardPool(SendCard{
				SendCard: card,
			})
			break
		} else {
			r.addCardsToCardHeap([]uno_pb.Card{card})
		}
	}
	r.stage = uno_pb.Stage_SendingCard
	// 重排顺序
	r.playerResequence()
}

func (r *Room) drawCard_ElectingBanker(p *player.Player) (*DrawCardEvent, *uno_pb.Errors) {
	if p.GetElectBankerCard() != nil {
		return nil, uno_pb.Errors_PlayerAlreadyDrawCard.Enum()
	}
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
		ps := []*uno_pb.PlayerAccountInfo{}
		for _, v := range r.players {
			ps = append(ps, v.FormatToProtoBuf().PlayerAccountInfo)
		}
		return &DrawCardEvent{
			IntoSendCard: true,
			IntoSendCardE: &uno_pb.RoomEventResponse_DrawCard_IntoSendCardEvent{
				Players:  ps,
				Banker:   r.banker.FormatToProtoBuf().PlayerAccountInfo,
				LeadCard: &r.GetLastCard().SendCard,
			},
		}, nil
	}
	return nil, nil
}

func (r *Room) effectedStackFeatureCard() {
FOROUT:
	for n := len(r.cardPool) - 1; n >= 0; n-- {
		sc := r.cardPool[n]
		if sc.SendCard.Type != uno_pb.CardType_Feature {
			break
		}
		fC := sc.SendCard.FeatureCard
		switch fC.FeatureCard {
		case uno_pb.FeatureCards_DrawTwo, uno_pb.FeatureCards_WildDrawFour:
			sc.featureEffected = true
		default:
			break FOROUT
		}
	}
}

func (r *Room) drawCard_SendingCard(p *player.Player) (*DrawCardEvent, *uno_pb.Errors) {
	// 可以抽牌的情况：
	// 1.遭到Wild draw four, Draw two(可以打出Skip或Wild draw four跳到下一个玩家，相关见SendCard)
	// 2.轮到该玩家出牌，但不想出或无牌可出
	if r.operatorNow.GetId() != p.GetId() {
		return nil, uno_pb.Errors_PlayerNoOperatorNow.Enum()
	}
	//
	stacks, ok := r.getStackFeatureCard()
	last := r.GetLastCard()
	if r.sendCard_checkSkipCard(last) {
		return nil, uno_pb.Errors_PlayerCannotDrawCard.Enum()
	} else if ok {
		dts := stacks[0]
		wdfs := stacks[1]
		if dts.count > 0 {
			cards := r.cutCards(2 * dts.count)
			p.AddCards(cards)
		}
		if wdfs.count > 0 {
			extra := 0
			if last.wildDrawFourStatus != nil {
				switch *last.wildDrawFourStatus {
				default:
					return nil, uno_pb.Errors_Unexpected.Enum()
				case wildDrawFourStatus_challengerLose: //挑战者失败
					extra = 2
				}
			}
			cards := r.cutCards(4*wdfs.count + extra)
			p.AddCards(cards)
		}
		p.SetCallUNO(false)
		r.effectedStackFeatureCard()
		next := r.nextOperator()
		r.operatorNow = next
		return &DrawCardEvent{
			Skipped: true,
			SkippedE: &uno_pb.RoomEventResponse_DrawCard_SkippedEvent{
				NextOperator: r.operatorNow.FormatToProtoBuf().PlayerAccountInfo,
			},
		}, nil
	} else if p.GetDrawCard() != nil { //已抽过一张牌
		return nil, uno_pb.Errors_PlayerAlreadyDrawCard.Enum()
	}
	// 玩家回合抽牌
	card := r.cutCards(1)[0]
	p.SetDrawCard(card)
	return nil, nil
}

type stackCardInfo struct {
	fc    uno_pb.FeatureCards
	count int
}

func (r *Room) getStackFeatureCard() ([2]stackCardInfo, bool) {
	ret := [2]stackCardInfo{
		stackCardInfo{
			fc:    uno_pb.FeatureCards_DrawTwo,
			count: 0,
		},
		stackCardInfo{
			fc:    uno_pb.FeatureCards_WildDrawFour,
			count: 0,
		},
	}
	last := r.GetLastCard()
	if last == nil {
		return ret, false
	}
FOROUT:
	for n := len(r.cardPool) - 1; n >= 0; n-- {
		sc := r.cardPool[n]
		if sc.SendCard.Type != uno_pb.CardType_Feature {
			break
		}
		fC := sc.SendCard.FeatureCard
		switch fC.FeatureCard {
		case uno_pb.FeatureCards_DrawTwo:
			if sc.featureEffected {
				break FOROUT
			}
			ret[0].count++
		case uno_pb.FeatureCards_WildDrawFour:
			if sc.featureEffected && last.wildDrawFourStatus != nil {
				break
			}
			ret[1].count++
		default:
			break FOROUT
		}
	}
	if ret[0].count == 0 && ret[1].count == 0 {
		return ret, false
	} else {
		return ret, true
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
	r.cardResequence(cards)
	r.cardHeap = make([]uno_pb.Card, 108)
	copy(r.cardHeap, cards)
	r.stage = uno_pb.Stage_ElectingBanker
}

func (r *Room) CallUNO(p *player.Player) ([]uno_pb.Card, *uno_pb.Errors) {
	if r.stage != uno_pb.Stage_SendingCard {
		return nil, uno_pb.Errors_RoomNoSendingCard.Enum()
	}
	l := len(p.GetCards())
	if p.GetDrawCard() != nil {
		l++
	}
	if l != 2 {
		cards := r.cutCards(2)
		p.AddCards(cards)
		return p.GetCards(), uno_pb.Errors_PlayerCannotCallUNO.Enum()
	}
	if p.GetCallUNO() {
		return nil, uno_pb.Errors_PlayerAlreadyCallUNO.Enum()
	}
	p.SetCallUNO(true)
	return nil, nil
}

func (r *Room) Challenge(p *player.Player) (bool, *uno_pb.PlayerInfo, *uno_pb.Errors) {
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
	if last.wildDrawFourStatus != nil {
		return false, nil, uno_pb.Errors_Challenged.Enum()
	}
	clr := last.SendCard.FeatureCard.Color
	lastP, ok := r.GetPlayer(last.SenderId)
	if !ok {
		return false, nil, uno_pb.Errors_Unexpected.Enum()
	}
	win := false
FOROUT:
	for _, v := range lastP.GetCards() {
		switch v.Type {
		case uno_pb.CardType_Normal:
			if v.NormalCard.Color == clr {
				win = true
				break FOROUT
			}
		case uno_pb.CardType_Feature:
			if v.FeatureCard.Color == clr {
				win = true
				break FOROUT
			}
		}
	}
	if win {
		last.wildDrawFourStatus = new(wildDrawFourStatus)
		*last.wildDrawFourStatus = wildDrawFourStatus_challengedLose
		last.featureEffected = true
		cards := r.cutCards(4)
		lastP, ok := r.GetPlayer(last.SenderId)
		if !ok {
			return false, nil, uno_pb.Errors_Unexpected.Enum()
		}
		lastP.AddCards(cards)
		return true, lastP.FormatToProtoBuf(), nil
	} else {
		last.wildDrawFourStatus = new(wildDrawFourStatus)
		*last.wildDrawFourStatus = wildDrawFourStatus_challengerLose
		return false, nil, nil
	}
}

func (r *Room) IndicateUNO(tP *player.Player) (*uno_pb.PlayerInfo, bool, *uno_pb.Errors) {
	if r.stage != uno_pb.Stage_SendingCard {
		return nil, false, uno_pb.Errors_RoomNoSendingCard.Enum()
	}
	if r.operatorNow.GetId() == tP.GetId() {
		return nil, false, uno_pb.Errors_PlayerIsOperatorNow.Enum()
	} else if len(tP.GetCards()) < 2 {
		if tP.GetCallUNO() {
			return nil, false, uno_pb.Errors_PlayerAlreadyCallUNO.Enum()
		}
		cards := r.cutCards(2)
		tP.AddCards(cards)
		return tP.FormatToProtoBuf(), true, nil
	} else {
		return nil, false, uno_pb.Errors_PlayerCannotCallUNO.Enum()
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

func (r *Room) cardResequence(x []uno_pb.Card) {
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
}

func (r *Room) playerResequence() {
	var (
		xCopy        []*player.Player
		sequenceCopy []string
	)
	xCopy = append(xCopy, r.banker)
	sequenceCopy = append(sequenceCopy, r.banker.GetId())
	for _, v := range r.players {
		if v.GetId() == r.banker.GetId() {
			continue
		}
		xCopy = append(xCopy, v)
		sequenceCopy = append(sequenceCopy, v.GetId())
	}
	r.players = xCopy
	r.sequencePosition = 0
}

func (r *Room) deletePlayer(playerid string) bool {
	for i, v := range r.players {
		if v.GetId() == playerid {
			if r.sequencePosition > i {
				r.sequencePosition = 0
			}
			if len(r.players) == 1 {
				r.players = []*player.Player{}
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
	for _, v := range r.players {
		if v.GetId() == id {
			return v, true
		}
	}
	return nil, false
}

func (r *Room) GetPlayerFromHash(playerHash string) (*player.Player, bool) {
	for _, v := range r.players {
		if v.GetHash() == playerHash {
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

func (r *Room) FormatToProtoBufExtra() *uno_pb.RoomExtra {
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
			FeatureEffected:    v.featureEffected,
		})
	}
	for _, v := range r.players {
		ps = append(ps, v.FormatToProtoBuf())
	}
	ur := &uno_pb.RoomExtra{
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

func (r *Room) FormatToProtoBufSimple() *uno_pb.RoomSimple {
	ch := []*uno_pb.Card{}
	cp := []*uno_pb.SendCard{}
	ps := []*uno_pb.PlayerAccountInfo{}
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
			FeatureEffected:    v.featureEffected,
		})
	}
	for _, v := range r.players {
		ps = append(ps, v.FormatToProtoBuf().PlayerAccountInfo)
	}
	ur := &uno_pb.RoomSimple{
		Hash:        r.hash,
		Stage:       r.stage,
		Banker:      nil,
		OperatorNow: nil,
		Players:     ps,
	}
	if r.banker != nil {
		ur.Banker = r.banker.FormatToProtoBuf().PlayerAccountInfo
	}
	if r.operatorNow != nil {
		ur.OperatorNow = r.operatorNow.FormatToProtoBuf().PlayerAccountInfo
	}
	return ur
}

func (r *Room) GetSeeds() [2]uint64 {
	return r.seeds
}
