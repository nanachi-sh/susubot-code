package player

import (
	uno_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
)

type Player struct {
	id   string
	name string

	roomhash        string
	cards           []uno_pb.Card
	electBankerCard *uno_pb.Card
	drawCard        *uno_pb.Card

	callUNO bool
}

type WildDrawFourStatus int

const (
	WildDrawFourStatus_Novs           WildDrawFourStatus = iota //不挑战
	WildDrawFourStatus_ChallengerLose                           //挑战者失败
	WildDrawFourStatus_ChallengedLose                           //被挑战者失败
)

func New(pi *uno_pb.PlayerInfo) *Player {
	p := &Player{
		id:   pi.PlayerAccountInfo.Id,
		name: pi.PlayerAccountInfo.Name,
	}
	if pi.PlayerRoomInfo != nil {
		p.roomhash = pi.PlayerRoomInfo.RoomHash
	}
	return p
}

func (p *Player) GetId() string {
	return p.id
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetCards() []uno_pb.Card {
	return p.cards
}

func (p *Player) GetRoomHash() string {
	return p.roomhash
}

func (p *Player) SetElectBankerCard(card uno_pb.Card) {
	p.electBankerCard = &card
}

func (p *Player) GetElectBankerCard() *uno_pb.Card {
	return p.electBankerCard
}

func (p *Player) AddCards(cards []uno_pb.Card) {
	p.cards = append(p.cards, cards...)
}

func (p *Player) DeleteCardFromHandCard(x uno_pb.Card) bool {
	for i, v := range p.cards {
		if v.Type != x.Type {
			continue
		}
		switch x.Type {
		case uno_pb.CardType_Normal:
			vNC := v.NormalCard
			xNC := x.NormalCard
			if vNC.Number == xNC.Number && vNC.Color == xNC.Color {
				p.deleteCardFromPosition(i)
				return true
			}
		case uno_pb.CardType_Feature:
			vFC := v.FeatureCard
			xFC := x.FeatureCard
			if vFC.Color == xFC.Color && vFC.FeatureCard == xFC.FeatureCard {
				p.deleteCardFromPosition(i)
				return true
			}
		}
	}
	return false
}

func (p *Player) DeleteCardFromDrawCard(x uno_pb.Card) bool {
	if p.drawCard == nil {
		return false
	}
	dc := *p.drawCard
	if x.Type != dc.Type {
		return false
	}
	switch x.Type {
	case uno_pb.CardType_Normal:
		if x.NormalCard.Color != dc.NormalCard.Color || x.NormalCard.Number != dc.NormalCard.Number {
			return false
		}
	case uno_pb.CardType_Feature:
		if x.FeatureCard.Color != dc.FeatureCard.Color || x.FeatureCard.FeatureCard != dc.FeatureCard.FeatureCard {
			return false
		}
	}
	p.ClearDrawCard()
	return true
}

func (p *Player) deleteCardFromPosition(postion int) {
	if len(p.cards) == 1 {
		p.cards = p.cards[0:]
	} else {
		p.cards = append(p.cards[:postion], p.cards[postion+1:]...)
	}
}

func (p *Player) DeleteCards(x []uno_pb.Card) bool {
	cardsCopy := make([]uno_pb.Card, len(p.cards))
	copy(cardsCopy, p.cards)
	for _, v1 := range x {
		ok := false
	OUTFOR:
		for n, v2 := range cardsCopy {
			switch {
			case v1.FeatureCard != nil:
				v1fc := v1.FeatureCard
				if v1fc.FeatureCard == v2.FeatureCard.FeatureCard && v1fc.Color == v2.FeatureCard.Color {
					ok = true
					cardsCopy = append(cardsCopy[:n], cardsCopy[n+1:]...)
					break OUTFOR
				}
			case v1.NormalCard != nil:
				v1nc := v1.NormalCard
				if v1nc.Number == v2.NormalCard.Number && v1nc.Color == v2.NormalCard.Color {
					ok = true
					cardsCopy = append(cardsCopy[:n], cardsCopy[n+1:]...)
					break OUTFOR
				}
			}
		}
		if !ok {
			return false
		}
	}
	p.cards = cardsCopy
	return true
}

func (p *Player) SetDrawCard(card uno_pb.Card) {
	p.drawCard = &card
}

func (p *Player) GetDrawCard() *uno_pb.Card {
	return p.drawCard
}

func (p *Player) ClearDrawCard() {
	p.drawCard = nil
}

func (p *Player) SetCallUNO(x bool) {
	p.callUNO = x
}

func (p *Player) FormatToProtoBuf() *uno_pb.PlayerInfo {
	pi := &uno_pb.PlayerInfo{
		PlayerAccountInfo: &uno_pb.PlayerAccountInfo{
			Id:   p.id,
			Name: p.name,
		},
	}
	if p.roomhash != "" {
		cards := []*uno_pb.Card{}
		for _, v := range p.cards {
			cards = append(cards, &v)
		}
		pi.PlayerRoomInfo = &uno_pb.PlayerRoomInfo{
			RoomHash: p.roomhash,
			Cards:    cards,
		}
	}
	return pi
}
