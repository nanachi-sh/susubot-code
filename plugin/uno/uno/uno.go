package uno

import (
	uno_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
	"github.com/nanachi-sh/susubot-code/plugin/uno/uno/player"
	"github.com/nanachi-sh/susubot-code/plugin/uno/uno/room"
)

var (
	rooms []*room.Room
)

func CreateRoom() *uno_pb.CreateRoomResponse {
	newRoom := room.New()
	rooms = append(rooms, newRoom)
	return &uno_pb.CreateRoomResponse{
		RoomHash: newRoom.GetHash(),
	}
}

func getRoom(hash string) (*room.Room, bool) {
	for _, v := range rooms {
		if v.GetHash() == hash {
			return v, true
		}
	}
	return nil, false
}

func getPlayerFromRooms(id string) (*player.Player, bool) {
	for _, v := range rooms {
		if p, ok := v.GetPlayer(id); ok {
			return p, ok
		}
	}
	return nil, false
}

func JoinRoom(req *uno_pb.JoinRoomRequest) *uno_pb.JoinRoomResponse {
	r, ok := getRoom(req.RoomHash)
	if !ok {
		return &uno_pb.JoinRoomResponse{
			Err: uno_pb.Errors_RoomNoExist.Enum(),
		}
	}
	if _, ok := getPlayerFromRooms(req.PlayerId); ok {
		return &uno_pb.JoinRoomResponse{
			Err: uno_pb.Errors_RoomExistPlayer.Enum(),
		}
	}
	if serr := r.Join(&uno_pb.PlayerAccountInfo{
		Id:   req.PlayerId,
		Name: req.PlayerName,
	}); serr != nil {
		return &uno_pb.JoinRoomResponse{
			Err: serr,
		}
	}
	ps := []*uno_pb.PlayerInfo{}
	for _, v := range r.GetPlayers() {
		ps = append(ps, v.FormatToProtoBuf())
	}
	return &uno_pb.JoinRoomResponse{
		Players: ps,
	}
}

func GetRooms() *uno_pb.GetRoomsResponse {
	rs := []*uno_pb.Room{}
	for _, v := range rooms {
		rs = append(rs, v.FormatToProtoBuf())
	}
	return &uno_pb.GetRoomsResponse{
		Infos: rs,
	}
}

func GetRoom(req *uno_pb.GetRoomRequest) *uno_pb.GetRoomResponse {
	if r, ok := getRoom(req.RoomHash); !ok {
		return &uno_pb.GetRoomResponse{
			Err: uno_pb.Errors_RoomNoExist.Enum(),
		}
	} else {
		return &uno_pb.GetRoomResponse{Info: r.FormatToProtoBuf()}
	}
}

func ExitRoom(req *uno_pb.ExitRoomRequest) *uno_pb.ExitRoomResponse {
	p, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &uno_pb.ExitRoomResponse{
			Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
		}
	}
	r, ok := getRoom(p.GetRoomHash())
	if !ok {
		return &uno_pb.ExitRoomResponse{
			Err: uno_pb.Errors_Unexpected.Enum(),
		}
	}
	if serr := r.Exit(req.PlayerId); serr != nil {
		return &uno_pb.ExitRoomResponse{Err: serr}
	}
	ps := []*uno_pb.PlayerInfo{}
	for _, v := range r.GetPlayers() {
		ps = append(ps, v.FormatToProtoBuf())
	}
	return &uno_pb.ExitRoomResponse{
		Players: ps,
	}
}

// func StartRoom(req *uno_pb.StartRoomRequest) *uno_pb.BasicResponse {
// 	switch {
// 	case req.RoomHash != nil:
// 		r, ok := getRoom(*req.RoomHash)
// 		if !ok {
// 			return &uno_pb.BasicResponse{
// 				Err: uno_pb.Errors_Unexpected.Enum(),
// 			}
// 		}
// 		if serr := r.Start(); serr != nil {
// 			return &uno_pb.BasicResponse{Err: serr}
// 		}
// 		return &uno_pb.BasicResponse{}
// 	case req.PlayerId != nil:
// 		p, ok := getPlayerFromRooms(*req.PlayerId)
// 		if !ok {
// 			return &uno_pb.BasicResponse{
// 				Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
// 			}
// 		}
// 		r, ok := getRoom(p.GetRoomHash())
// 		if !ok {
// 			return &uno_pb.BasicResponse{
// 				Err: uno_pb.Errors_Unexpected.Enum(),
// 			}
// 		}
// 		if serr := r.Start(); serr != nil {
// 			return &uno_pb.BasicResponse{Err: serr}
// 		}
// 		return &uno_pb.BasicResponse{}
// 	default:
// 		return &uno_pb.BasicResponse{Err: uno_pb.Errors_Unexpected.Enum()}
// 	}
// }

func DrawCard(req *uno_pb.DrawCardRequest) *uno_pb.DrawCardResponse {
	p, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &uno_pb.DrawCardResponse{
			Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
		}
	}
	r, ok := getRoom(p.GetRoomHash())
	if !ok {
		return &uno_pb.DrawCardResponse{
			Err: uno_pb.Errors_Unexpected.Enum(),
		}
	}
	e, serr := r.DrawCard(p)
	if serr != nil {
		return &uno_pb.DrawCardResponse{Err: serr}
	}
	if e != nil {
		switch {
		case e.IntoSendCard:
			return &uno_pb.DrawCardResponse{
				ElectBankerCard: p.GetElectBankerCard(),
				Stage:           r.GetStage(),
				IntoSendCard:    e.IntoSendCard,
				IntoSendCardE:   e.IntoSendCardE,
			}
		case e.Skipped:
			cs := []*uno_pb.Card{}
			for _, v := range p.GetCards() {
				cs = append(cs, &v)
			}
			return &uno_pb.DrawCardResponse{
				PlayerCard: cs,
				DrawCard:   p.GetDrawCard(),
				Stage:      r.GetStage(),
				Skipped:    e.Skipped,
				SkippedE:   e.SkippedE,
			}
		}
	}
	switch stage := r.GetStage(); stage {
	case uno_pb.Stage_ElectingBanker:
		return &uno_pb.DrawCardResponse{
			ElectBankerCard: p.GetElectBankerCard(),
			Stage:           stage,
		}
	case uno_pb.Stage_SendingCard:
		cs := []*uno_pb.Card{}
		for _, v := range p.GetCards() {
			cs = append(cs, &v)
		}
		return &uno_pb.DrawCardResponse{
			PlayerCard: cs,
			DrawCard:   p.GetDrawCard(),
			Stage:      stage,
		}
	default:
		return &uno_pb.DrawCardResponse{Err: uno_pb.Errors_Unexpected.Enum()}
	}
}

// func SendCardAction(req *uno_pb.SendCardActionRequest) *uno_pb.SendCardActionResponse {
// 	p, ok := getPlayerFromRooms(req.PlayerId)
// 	if !ok {
// 		return &uno_pb.SendCardActionResponse{
// 			Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
// 		}
// 	}
// 	r, ok := getRoom(p.GetRoomHash())
// 	if !ok {
// 		return &uno_pb.SendCardActionResponse{
// 			Err: uno_pb.Errors_Unexpected.Enum(),
// 		}
// 	}
// 	if req.Action == uno_pb.SendCardActions_Send && req.SendCard == nil {
// 		return &uno_pb.SendCardActionResponse{Err: uno_pb.Errors_Unexpected.Enum()}
// 	}
// 	sc := uno_pb.Card{}
// 	if req.SendCard != nil {
// 		sc = *req.SendCard
// 	}
// 	next, e, serr := r.SendCardAction(p, sc, req.Action)
// 	if serr != nil {
// 		return &uno_pb.SendCardActionResponse{Err: serr}
// 	}
// 	resp := new(uno_pb.SendCardActionResponse)
// 	if next != nil {
// 		resp.NextOperator = next.FormatToProtoBuf()
// 		if req.Action == uno_pb.SendCardActions_Send {
// 			cs := []*uno_pb.Card{}
// 			for _, v := range p.GetCards() {
// 				cs = append(cs, &v)
// 			}
// 			resp.SenderCard = cs
// 		}
// 	}
// 	if e != nil {
// 		if e.GameFinish {
// 			if !deleteRoom(r) {
// 				return &uno_pb.SendCardActionResponse{Err: uno_pb.Errors_Unexpected.Enum()}
// 			}
// 			resp.GameFinish = e.GameFinish
// 			resp.GameFinishE = e.GameFinishE
// 		}
// 	}
// 	return resp
// }

func CallUNO(req *uno_pb.CallUNORequest) *uno_pb.CallUNOResponse {
	p, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &uno_pb.CallUNOResponse{
			Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
		}
	}
	r, ok := getRoom(p.GetRoomHash())
	if !ok {
		return &uno_pb.CallUNOResponse{
			Err: uno_pb.Errors_Unexpected.Enum(),
		}
	}
	cards, serr := r.CallUNO(p)
	cs := []*uno_pb.Card{}
	for _, v := range cards {
		cs = append(cs, &v)
	}
	return &uno_pb.CallUNOResponse{
		PlayerCard: cs,
		Err:        serr,
	}
}

func Challenge(req *uno_pb.ChallengeRequest) *uno_pb.ChallengeResponse {
	p, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &uno_pb.ChallengeResponse{
			Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
		}
	}
	r, ok := getRoom(p.GetRoomHash())
	if !ok {
		return &uno_pb.ChallengeResponse{
			Err: uno_pb.Errors_Unexpected.Enum(),
		}
	}
	win, pi, serr := r.Challenge(p)
	if serr != nil {
		return &uno_pb.ChallengeResponse{Err: serr}
	}
	return &uno_pb.ChallengeResponse{
		Win:        win,
		LastPlayer: pi,
	}
}

func IndicateUNO(req *uno_pb.IndicateUNORequest) *uno_pb.IndicateUNOResponse {
	tp, ok := getPlayerFromRooms(req.TargetId)
	if !ok {
		return &uno_pb.IndicateUNOResponse{
			Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
		}
	}
	r, ok := getRoom(tp.GetRoomHash())
	if !ok {
		return &uno_pb.IndicateUNOResponse{
			Err: uno_pb.Errors_Unexpected.Enum(),
		}
	}
	pi, ok, serr := r.IndicateUNO(tp)
	if serr != nil {
		return &uno_pb.IndicateUNOResponse{Err: serr}
	}
	return &uno_pb.IndicateUNOResponse{
		Punished:   pi,
		IndicateOK: ok,
	}
}

func deleteRoom(r *room.Room) bool {
	for i, v := range rooms {
		if v.GetHash() == r.GetHash() {
			rooms = append(rooms[:i], rooms[i+1:]...)
			return true
		}
	}
	return false
}

func TEST_SetPlayerCard(req *uno_pb.TEST_SetPlayerCardRequest) *uno_pb.BasicResponse {
	p, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &uno_pb.BasicResponse{
			Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
		}
	}
	p.ClearHandCard()
	cs := []uno_pb.Card{}
	for _, v := range req.Cards {
		cs = append(cs, *v)
	}
	p.AddCards(cs)
	return nil
}
