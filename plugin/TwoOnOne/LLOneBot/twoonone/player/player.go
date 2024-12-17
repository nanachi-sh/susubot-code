package player

import (
	"fmt"
	"sync"
	"time"

	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/protos/twoonone"
)

type Player struct {
	setRobLandownerActionTime time.Time
	lock                      sync.Mutex

	id        string
	name      string
	winCount  int
	loseCount int
	coin      float64

	originPlayerInfo *twoonone_pb.PlayerAccountInfo

	roomHash           string
	cards              []twoonone_pb.Card
	robLandownerAction *twoonone_pb.RobLandownerActions
}

func New(pi *twoonone_pb.PlayerInfo) *Player {
	if pi.AccountInfo == nil {
		return nil
	}
	p := &Player{
		id:               pi.AccountInfo.Id,
		name:             pi.AccountInfo.Name,
		winCount:         int(pi.AccountInfo.WinCount),
		loseCount:        int(pi.AccountInfo.LoseCount),
		coin:             pi.AccountInfo.Coin,
		originPlayerInfo: pi.AccountInfo,
	}
	if pi.TableInfo != nil {
		p.roomHash = pi.TableInfo.RoomHash
		p.cards = pi.TableInfo.Cards
		p.robLandownerAction = pi.TableInfo.RobLandownerAction
	}
	return p
}

func (p *Player) IsEmpty() bool {
	return p.id == ""
}

func (p *Player) GetId() string {
	return p.id
}

func (p *Player) GetName() string {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.name
}

func (p *Player) GetCoin() float64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.coin
}

func (p *Player) GetWinCount() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.winCount
}

func (p *Player) GetLoseCount() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.loseCount
}

func (p *Player) GetRoomHash() string {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.roomHash
}

func (p *Player) GetCards() []twoonone_pb.Card {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.cards
}

func (p *Player) AddCards(x []twoonone_pb.Card) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.cards = append(p.cards, x...)
}

func (p *Player) DeleteCards(x []twoonone_pb.Card) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	cardsCopy := make([]twoonone_pb.Card, len(p.cards))
	copy(cardsCopy, p.cards)
	for _, v1 := range x {
		ok := false
		for n, v2 := range cardsCopy {
			if v1 == v2 {
				ok = true
				cardsCopy = append(cardsCopy[:n], cardsCopy[n+1:]...)
				break
			}
		}
		if !ok {
			return false
		}
	}
	p.cards = cardsCopy
	return true
}

func (p *Player) ClearCards() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.cards = p.cards[:0]
}

func (p *Player) SetRobLandownerAction(x *twoonone_pb.RobLandownerActions) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.setRobLandownerActionTime = time.Now()
	p.robLandownerAction = x
}

func (p *Player) GetRobLandownerAction() *twoonone_pb.RobLandownerActions {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.robLandownerAction
}

func (p *Player) GetRobLandownerActionTime() time.Time {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.setRobLandownerActionTime
}

func (p *Player) GetLastGetDailyTimestamp() int64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	fmt.Println(p.originPlayerInfo)
	fmt.Println(p.originPlayerInfo == nil)
	fmt.Println(p.originPlayerInfo.LastGetDailyTimestamp)
	return p.originPlayerInfo.LastGetDailyTimestamp
}

func (p *Player) DecCoin(x float64) {
	p.coin -= x
}

func (p *Player) IncCoin(x float64) {
	p.coin += x
}

func (p *Player) IncWinCount() {
	p.winCount++
}

func (p *Player) IncLoseCount() {
	p.loseCount++
}

func (p *Player) GetOriginInfo() twoonone_pb.PlayerAccountInfo {
	return *p.originPlayerInfo
}
