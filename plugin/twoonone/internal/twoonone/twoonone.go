package twoonone

import (
	"net/http"
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/middleware/sql"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone/player"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone/room"
	internal_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/types"
	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	rooms []*room.Room
	db    = sql.NewHandler()
)

// type GRPCRequest struct {
// 	logger logx.Logger
// }

// func NewRequest(l logx.Logger) *Request {
// 	return &Request{logger: l}
// }

// func (r *Request) CreateRoom(in *types.CreateRoomRequest) (*twoonone_pb.CreateRoomResponse, error) {
// 	return createRoom(in), nil
// }

// func (r *Request) ExitRoom(in *types.ExitRoomRequest) (*twoonone_pb.ExitRoomResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" {
// 		return &twoonone_pb.ExitRoomResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return exitRoom(r.logger, in), nil
// }

// func (r *Request) EventRoom(in *types.EventRoomRequest, stream twoonone_pb.Twoonone_EventRoomServer) error {
// 	if in.RoomHash == "" {
// 		return status.Error(codes.InvalidArgument, "")
// 	}
// 	e, ok := event.FindEventStream(in.RoomHash)
// 	if !ok {
// 		if err := stream.Send(&twoonone_pb.EventRoomResponse{
// 			Body: &twoonone_pb.EventRoomResponse_Error{Error: utils.GenerateError(
// 				r.logger,
// 				twoonone_pb.ErrorType_ERROR_TYPE_GAME,
// 				&twoonone_pb.GameError{Error: twoonone_pb.GameError_ERROR_ROOM_NO_EXIST},
// 			)},
// 		}); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	for {
// 		resp, err := e.Read()
// 		if err != nil {
// 			return nil
// 		}
// 		if err := stream.Send(resp); err != nil {
// 			return err
// 		}
// 	}
// }

// func (r *Request) GetDailyCoin(in *types.GetDailyCoinRequest) (*twoonone_pb.GetDailyCoinResponse, error) {
// 	if in.UserId == "" {
// 		return &twoonone_pb.GetDailyCoinResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return getDailyCoin(r.logger, in), nil
// }

// func (r *Request) GetRoom(in *types.GetRoomRequest) (*twoonone_pb.GetRoomResponse, error) {
// 	if in.RoomHash == "" {
// 		return &twoonone_pb.GetRoomResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return getRoom(r.logger, in), nil
// }

// func (r *Request) GetRooms(in *types.GetRoomsRequest) (*twoonone_pb.GetRoomsResponse, error) {
// 	return getRooms(), nil
// }

// func (r *Request) JoinRoom(in *types.JoinRoomRequest) (*twoonone_pb.JoinRoomResponse, error) {
// 	if in.RoomHash == "" || in.UserId == "" {
// 		return &twoonone_pb.JoinRoomResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return joinRoom(r.logger, in), nil
// }

// func (r *Request) RobLandowner(in *types.RobLandownerRequest) (*twoonone_pb.RobLandownerResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" {
// 		return &twoonone_pb.RobLandownerResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return robLandowner(r.logger, in), nil
// }

// func (r *Request) NoRobLandowner(in *types.NoRobLandownerRequest) (*twoonone_pb.NoRobLandownerResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" {
// 		return &twoonone_pb.NoRobLandownerResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return noRobLandowner(r.logger, in), nil
// }

// func (r *Request) SendCard(in *types.SendCardRequest) (*twoonone_pb.SendCardResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" || len(in.Sendcards) == 0 {
// 		return &twoonone_pb.SendCardResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return sendCard(r.logger, in), nil
// }

// func (r *Request) NoSendCard(in *types.NoSendCardRequest) (*twoonone_pb.NoSendCardResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" {
// 		return &twoonone_pb.NoSendCardResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return noSendCard(r.logger, in), nil
// }

// func (r *Request) StartRoom(in *types.StartRoomRequest) (*twoonone_pb.StartRoomResponse, error) {
// 	if in.RoomHash == "" {
// 		return &twoonone_pb.StartRoomResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return startRoom(r.logger, in), nil
// }

type APIRequest struct {
	logger logx.Logger
}

func NewAPIRequest(l logx.Logger) *APIRequest {
	return &APIRequest{
		logger: l,
	}
}

func setToMap(m map[string]string, key, value string) {
	if m == nil {
		m = make(map[string]string)
	}
	m[key] = value
}

func (r *APIRequest) GetRoom(req *types.GetRoomRequest) (any, error) {
	if req.RoomHash == "" {
		return nil, types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	resp, err := getRoom(r.logger, req)
	if err != nil {
		return nil, err
	}
	if req.Extra.NewExtra != "" {
		setToMap(resp.Extra, internal_types.EXTRA_KEY_extra, req.Extra.NewExtra)
	}
	return resp, nil
}

func (r *APIRequest) ExitRoom(req *types.ExitRoomRequest) (any, error) {
	if req.RoomHash == "" {
		return nil, types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	resp, err := exitRoom(r.logger, req)
	if err != nil {
		return nil, err
	}
	if req.Extra.NewExtra != "" {
		setToMap(resp.Extra, internal_types.EXTRA_KEY_extra, req.Extra.NewExtra)
	}
	return resp, nil
}

func (r *APIRequest) GetDailyCoin(req *types.GetDailyCoinRequest) (any, error) {
	return getDailyCoin(r.logger, req)
}

func (r *APIRequest) GetRooms() (any, error) {
	return getRooms(), nil
}

func (r *APIRequest) JoinRoom(req *types.JoinRoomRequest) (any, error) {
	if req.RoomHash == "" {
		return nil, types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return joinRoom(r.logger, req)
}

func (r *APIRequest) NoRobLandowner(req *types.NoRobLandownerRequest) (any, error) {
	if req.RoomHash == "" {
		return nil, types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return noRobLandowner(r.logger, req)
}

func (r *APIRequest) NoSendCard(req *types.NoSendCardRequest) (any, error) {
	if req.RoomHash == "" {
		return nil, types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return noSendCard(r.logger, req)
}

func (r *APIRequest) RobLandowner(req *types.RobLandownerRequest) (any, error) {
	if req.RoomHash == "" {
		return nil, types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return robLandowner(r.logger, req)
}

func (r *APIRequest) SendCard(req *types.SendCardRequest) (any, error) {
	if req.RoomHash == "" || len(req.SendCards) == 0 {
		return nil, types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return sendCard(r.logger, req)
}

func (r *APIRequest) StartRoom(req *types.StartRoomRequest) (any, error) {
	if req.RoomHash == "" {
		return nil, types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return startRoom(r.logger, req)
}

func (r *APIRequest) CreateRoom() (any, error) {
	return createRoom()
}

func getPlayerFromRooms(id string) (*player.Player, bool) {
	for _, v := range rooms {
		pi, ok := v.GetPlayer(id)
		if ok {
			return pi, ok
		}
	}
	return nil, false
}

func getPlayerFromRoom(playerId string, tableHash string) (*player.Player, bool) {
	for _, v := range rooms {
		if v.GetHash() == tableHash {
			return v.GetPlayer(playerId)
		}
	}
	return nil, false
}

func findRoom(hash string) (*room.Room, bool) {
	for _, v := range rooms {
		if v.GetHash() == hash {
			return v, true
		}
	}
	return nil, false
}

func parseCard(c []types.Card) []*twoonone_pb.Card {
	cards := []*twoonone_pb.Card{}
	for _, v := range c {
		cards = append(cards, &twoonone_pb.Card{
			Number: twoonone_pb.Card_Number(v.Number),
		})
	}
	return cards
}

func getRooms() *twoonone_pb.GetRoomsResponse {
	return &twoonone_pb.GetRoomsResponse{
		RoomInfos: room.FormatInternalRooms2Protobuf(rooms),
	}
}

func getDailyCoin(logger logx.Logger, in *types.GetDailyCoinRequest) (*twoonone_pb.GetDailyCoinResponse, error) {
	if checkDailyCoin(in.Extra) {
		if err := db.UpdateUser(logger, in.Extra.UserId, sql.IncCoin(500), sql.UpdateGetDailyTime()); err != nil {
			return nil, err
		}
		return &twoonone_pb.GetDailyCoinResponse{}, nil
	} else {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_ALREADY_GET_DALIY_COIN, "")
	}
}

func checkDailyCoin(u types.Extra) bool {
	last_time := time.Unix(u.LastGetDaliyTime, 0)
	if last_time.IsZero() { //第一次领取
		return true
	} else {
		now := time.Now()
		// 判断是否为同一年同一月
		if now.Month() == last_time.Month() && now.Year() == last_time.Year() {
			return now.Day() > last_time.Day()
		} else {
			//非同一年同一月必然不是同一天
			return true
		}
	}
}

func getRoom(logger logx.Logger, in *types.GetRoomRequest) (*twoonone_pb.GetRoomResponse, error) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	return &twoonone_pb.GetRoomResponse{
		RoomInfo: room.FormatInternalRoom2Protobuf(r),
	}, nil
}

func createRoom() (*twoonone_pb.CreateRoomResponse, error) {
	newRoom := room.New(200, 1)
	rooms = append(rooms, newRoom)
	return &twoonone_pb.CreateRoomResponse{
		RoomHash: newRoom.GetHash(),
	}, nil
}

func joinRoom(logger logx.Logger, in *types.JoinRoomRequest) (*twoonone_pb.JoinRoomResponse, error) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	if _, ok := getPlayerFromRooms(in.Extra.UserId); ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_EXISTED_A_ROOM, "")
	}
	if err := r.Join(logger, in.Extra.UserId, in.Extra.Name, in.Extra.Coin); err != nil {
		return nil, err
	}
	return &twoonone_pb.JoinRoomResponse{}, nil
}

func exitRoom(logger logx.Logger, in *types.ExitRoomRequest) (*twoonone_pb.ExitRoomResponse, error) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	if err := r.Exit(logger, in.Extra.UserId); err != nil {
		return nil, err
	}
	return &twoonone_pb.ExitRoomResponse{}, nil
}

func startRoom(logger logx.Logger, in *types.StartRoomRequest) (*twoonone_pb.StartRoomResponse, error) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	if _, ok := r.GetPlayer(in.Extra.UserId); !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST, "")
	}
	if err := r.Start(logger); err != nil {
		return nil, err
	}
	return &twoonone_pb.StartRoomResponse{}, nil
}

func robLandowner(logger logx.Logger, in *types.RobLandownerRequest) (*twoonone_pb.RobLandownerResponse, error) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	p, ok := r.GetPlayer(in.Extra.UserId)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST, "")
	}
	if err := r.RobLandowner(logger, p); err != nil {
		return nil, err
	}
	return &twoonone_pb.RobLandownerResponse{}, nil
}

func noRobLandowner(logger logx.Logger, in *types.NoRobLandownerRequest) (*twoonone_pb.NoRobLandownerResponse, error) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	p, ok := r.GetPlayer(in.Extra.UserId)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST, "")
	}
	if err := r.NoRobLandowner(logger, p); err != nil {
		return nil, err
	}
	return &twoonone_pb.NoRobLandownerResponse{}, nil
}

func sendCard(logger logx.Logger, in *types.SendCardRequest) (*twoonone_pb.SendCardResponse, error) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	p, ok := r.GetPlayer(in.Extra.UserId)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST, "")
	}
	scs := parseCard(in.SendCards)
	if err := r.SendCard(logger, p, scs); err != nil {
		return nil, err
	}
	return &twoonone_pb.SendCardResponse{}, nil
}

func noSendCard(logger logx.Logger, in *types.NoSendCardRequest) (*twoonone_pb.NoSendCardResponse, error) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	p, ok := r.GetPlayer(in.Extra.UserId)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST, "")
	}
	if serr := r.NoSendCard(logger, p); serr != nil {
		return nil, serr
	}
	return &twoonone_pb.NoSendCardResponse{}, nil
}

func deleteRoomFromRooms(r *room.Room) bool {
	for i, v := range rooms {
		if v.GetHash() == r.GetHash() {
			rooms = append(rooms[:i], rooms[i+1:]...)
			return true
		}
	}
	return false
}
