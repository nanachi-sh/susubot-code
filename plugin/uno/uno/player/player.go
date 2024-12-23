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
}

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

func (r *Player) GetId() string {
	return r.id
}

func (r *Player) GetName() string {
	return r.name
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
