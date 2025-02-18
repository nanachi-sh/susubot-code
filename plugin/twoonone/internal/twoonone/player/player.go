package player

import (
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/middleware/sql"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone/card"
	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	db = sql.NewHandler()
)

type Player struct {
	created bool

	id   string
	name string

	roomHash         string
	cards            []card.Card
	robLandownerInfo twoonone_pb.RobLandownerInfo

	coinChanged      float64
	winCountChanged  int
	loseCountChanged int
}

func New(pi *twoonone_pb.PlayerInfo) *Player {
	return &Player{
		created:          true,
		id:               pi.User.Id,
		name:             pi.User.Name,
		roomHash:         pi.Table.RoomHash,
		cards:            []card.Card{},
		robLandownerInfo: twoonone_pb.RobLandownerInfo{},
	}
}

func (p *Player) IsEmpty() bool {
	return !p.created
}

func (p *Player) GetId() string {
	return p.id
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetRoomHash() string {
	return p.roomHash
}

func (p *Player) GetCards() []card.Card {
	return p.cards
}

func (p *Player) AddCards(x []card.Card) {
	p.cards = append(p.cards, x...)
}

func (p *Player) DeleteCards(x []card.Card) bool {
	cardsCopy := make([]card.Card, len(p.cards))
	copy(cardsCopy, p.cards)
	delete := func(cards_p *[]card.Card, target card.Card) bool {
		cards := *cards_p
		for n, v := range cards {
			if v.Number == target.Number {
				*cards_p = append(cards[:n], cards[n+1:]...)
				return true
			}
		}
		return false
	}
	for _, v := range x {
		if !delete(&cardsCopy, v) {
			return false
		}
	}
	p.cards = cardsCopy
	return true
}

func (p *Player) ClearCards() {
	p.cards = p.cards[:0]
}

func (p *Player) SetRobLandownerAction(x twoonone_pb.RobLandownerInfo_Action) {
	if x != twoonone_pb.RobLandownerInfo_ACTION_EMPTY {
		p.robLandownerInfo.ActionTime = time.Now().Unix()
	}
	p.robLandownerInfo.Action = x
}

func (p *Player) GetRobLandownerAction() twoonone_pb.RobLandownerInfo_Action {
	return p.robLandownerInfo.Action
}

func (p *Player) GetRobLandownerActionTime() time.Time {
	return time.Unix(p.robLandownerInfo.ActionTime, 0)
}

func (p *Player) DecCoin(x float64) {
	p.coinChanged -= x
}

func (p *Player) IncCoin(x float64) {
	p.coinChanged += x
}

func (p *Player) IncWinCount() {
	p.winCountChanged++
}

func (p *Player) IncLoseCount() {
	p.loseCountChanged++
}

func (p *Player) UpdateToDatabaseAndClean(logger logx.Logger) error {
	if p.coinChanged > 0 {
		if serr := db.UpdateUser(logger, p.GetId(), sql.IncCoin(p.coinChanged)); serr != nil {
			return serr
		}
	} else if p.coinChanged < 0 {
		if serr := db.UpdateUser(logger, p.GetId(), sql.DecCoin(p.coinChanged)); serr != nil {
			return serr
		}
	}

	if p.winCountChanged > 0 {
		if serr := db.UpdateUser(logger, p.GetId(), sql.IncWinCount(uint(p.winCountChanged))); serr != nil {
			return serr
		}
	} else if p.winCountChanged < 0 {
		if serr := db.UpdateUser(logger, p.GetId(), sql.DecWinCount(uint(p.winCountChanged))); serr != nil {
			return serr
		}
	}

	if p.loseCountChanged > 0 {
		if serr := db.UpdateUser(logger, p.GetId(), sql.IncLoseCount(uint(p.loseCountChanged))); serr != nil {
			return serr
		}
	} else if p.loseCountChanged < 0 {
		if serr := db.UpdateUser(logger, p.GetId(), sql.DecLoseCount(uint(p.loseCountChanged))); serr != nil {
			return serr
		}
	}
	p.resetChanged()
	return nil
}

func (p *Player) resetChanged() {
	p.coinChanged = 0
	p.winCountChanged = 0
	p.loseCountChanged = 0
}

func FormatInternalPlayer2Protobuf(x *Player) *twoonone_pb.PlayerInfo {
	if x == nil {
		return nil
	}
	return &twoonone_pb.PlayerInfo{
		User: &twoonone_pb.PlayerInfo_UserInfo{
			Id:   x.id,
			Name: x.name,
		},
		Table: &twoonone_pb.PlayerInfo_TableInfo{
			RoomHash:         x.roomHash,
			RoblandownerInfo: &x.robLandownerInfo,
		},
	}
}

func FormatInternalPlayers2Protobuf(xs []*Player) []*twoonone_pb.PlayerInfo {
	if xs == nil {
		return nil
	}
	ret := []*twoonone_pb.PlayerInfo{}
	for _, v := range xs {
		ret = append(ret, FormatInternalPlayer2Protobuf(v))
	}
	return ret
}
