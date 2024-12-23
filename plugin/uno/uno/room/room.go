package room

import (
	"fmt"
	"math/rand"
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

type SendCard struct {
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
	if r.GetStage() != uno_pb.Stage_WaitingStart {
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
	if !r.delete(playerid) {
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
	cardsArr := r.generateCards()
	cards := cardsArr[:]
	r.resequence(&cards)
	copy(r.cardHeap, cards)
	// 进入下一阶段
	r.startElectBanker()
	return nil
}

func (r *Room) DrawCard(p *player.Player) ([]uno_pb.Card, *uno_pb.Errors) {
	switch r.GetStage() {
	case uno_pb.Stage_ElectingBanker:
		card, serr := r.drawCard_ElectBanker(p)
		if serr != nil {
			return nil, serr
		}
		if len(r.getTakeElectBankerPlayers()) == len(r.GetPlayers()) {

		}
		return []uno_pb.Card{card}, nil
	}
}

func (r *Room) startSendCard() {
	// 重排顺序
	r.stage = uno_pb.Stage_SendingCard
}

func (r *Room) drawCard_ElectBanker(p *player.Player) (uno_pb.Card, *uno_pb.Errors) {
	p.SetElectBankerCard(r.cardHeap[0])
	r.cardHeap = r.cardHeap[1:]
	return *p.GetElectBankerCard(), nil
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
	r.stage = uno_pb.Stage_ElectingBanker
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
					})
				}
			case uno_pb.FeatureCards_Wild, uno_pb.FeatureCards_WildDrawFour:
				cards = append(cards, uno_pb.Card{
					FeatureCard: &uno_pb.FeatureCard{
						Color:       uno_pb.CardColor_Black,
						FeatureCard: uno_pb.FeatureCards(fea),
					},
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
	}
}

func (r *Room) delete(playerid string) bool {
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
