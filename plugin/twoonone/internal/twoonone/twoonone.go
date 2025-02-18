package twoonone

import (
	"net/http"
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/middleware/sql"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/model/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone/player"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone/room"
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

// func (r *Request) CreateRoom(in *twoonone_pb.CreateRoomRequest) (*twoonone_pb.CreateRoomResponse, error) {
// 	return createRoom(in), nil
// }

// func (r *Request) ExitRoom(in *twoonone_pb.ExitRoomRequest) (*twoonone_pb.ExitRoomResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" {
// 		return &twoonone_pb.ExitRoomResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return exitRoom(r.logger, in), nil
// }

// func (r *Request) EventRoom(in *twoonone_pb.EventRoomRequest, stream twoonone_pb.Twoonone_EventRoomServer) error {
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

// func (r *Request) GetDailyCoin(in *twoonone_pb.GetDailyCoinRequest) (*twoonone_pb.GetDailyCoinResponse, error) {
// 	if in.UserId == "" {
// 		return &twoonone_pb.GetDailyCoinResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return getDailyCoin(r.logger, in), nil
// }

// func (r *Request) GetRoom(in *twoonone_pb.GetRoomRequest) (*twoonone_pb.GetRoomResponse, error) {
// 	if in.RoomHash == "" {
// 		return &twoonone_pb.GetRoomResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return getRoom(r.logger, in), nil
// }

// func (r *Request) GetRooms(in *twoonone_pb.GetRoomsRequest) (*twoonone_pb.GetRoomsResponse, error) {
// 	return getRooms(), nil
// }

// func (r *Request) JoinRoom(in *twoonone_pb.JoinRoomRequest) (*twoonone_pb.JoinRoomResponse, error) {
// 	if in.RoomHash == "" || in.UserId == "" {
// 		return &twoonone_pb.JoinRoomResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return joinRoom(r.logger, in), nil
// }

// func (r *Request) RobLandowner(in *twoonone_pb.RobLandownerRequest) (*twoonone_pb.RobLandownerResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" {
// 		return &twoonone_pb.RobLandownerResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return robLandowner(r.logger, in), nil
// }

// func (r *Request) NoRobLandowner(in *twoonone_pb.NoRobLandownerRequest) (*twoonone_pb.NoRobLandownerResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" {
// 		return &twoonone_pb.NoRobLandownerResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return noRobLandowner(r.logger, in), nil
// }

// func (r *Request) SendCard(in *twoonone_pb.SendCardRequest) (*twoonone_pb.SendCardResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" || len(in.Sendcards) == 0 {
// 		return &twoonone_pb.SendCardResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return sendCard(r.logger, in), nil
// }

// func (r *Request) NoSendCard(in *twoonone_pb.NoSendCardRequest) (*twoonone_pb.NoSendCardResponse, error) {
// 	if in.PlayerId == "" || in.RoomHash == "" {
// 		return &twoonone_pb.NoSendCardResponse{}, status.Error(codes.InvalidArgument, "")
// 	}
// 	return noSendCard(r.logger, in), nil
// }

// func (r *Request) StartRoom(in *twoonone_pb.StartRoomRequest) (*twoonone_pb.StartRoomResponse, error) {
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

func (r *APIRequest) GetRoom(req *twoonone_pb.GetRoomRequest) (resp any, err error) {
	if req.RoomHash == "" {
		err = types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return getRoom(r.logger, req)
}

func (r *APIRequest) ExitRoom(req *twoonone_pb.ExitRoomRequest) (resp any, err error) {
	if req.RoomHash == "" || req.UserId == "" {
		err = types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return exitRoom(r.logger, req)
}

func (r *APIRequest) GetDailyCoin(req *twoonone_pb.GetDailyCoinRequest) (resp any, err error) {
	if req.UserId == "" {
		err = types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return getDailyCoin(r.logger, req)
}

func (r *APIRequest) GetRooms() (resp any, err error) {
	return getRooms(), nil
}

func (r *APIRequest) JoinRoom(req *twoonone_pb.JoinRoomRequest) (resp any, err error) {
	if req.RoomHash == "" || req.UserId == "" {
		err = types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return joinRoom(r.logger, req)
}

func (r *APIRequest) NoRobLandowner(req *twoonone_pb.NoRobLandownerRequest) (resp any, err error) {
	if req.RoomHash == "" || req.UserId == "" {
		err = types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return noRobLandowner(r.logger, req)
}

func (r *APIRequest) NoSendCard(req *twoonone_pb.NoSendCardRequest) (resp any, err error) {
	if req.RoomHash == "" || req.UserId == "" {
		err = types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return noSendCard(r.logger, req)
}

func (r *APIRequest) RobLandowner(req *twoonone_pb.RobLandownerRequest) (resp any, err error) {
	if req.RoomHash == "" || req.UserId == "" {
		err = types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return robLandowner(r.logger, req)
}

func (r *APIRequest) SendCard(req *twoonone_pb.SendCardRequest) (resp any, err error) {
	if req.RoomHash == "" || req.UserId == "" || len(req.Sendcards) == 0 {
		err = types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return sendCard(r.logger, req)
}

func (r *APIRequest) StartRoom(req *twoonone_pb.StartRoomRequest) (resp any, err error) {
	if req.RoomHash == "" || req.UserId == "" {
		err = types.NewError(twoonone_pb.Error_ERROR_INVALID_ARGUMENT, "", http.StatusBadRequest)
	}
	return startRoom(r.logger, &twoonone_pb.StartRoomRequest{
		RoomHash: req.RoomHash,
		UserId:   req.UserId,
	})
}

func (r *APIRequest) CreateRoom(req *twoonone_pb.CreateRoomRequest) (resp any, err error) {
	defer func() {
		r.logger.Info(err == nil)
		r.logger.Infof("%T", err)
	}()
	return createRoom(&twoonone_pb.CreateRoomRequest{})
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

func getRooms() *twoonone_pb.GetRoomsResponse {
	return &twoonone_pb.GetRoomsResponse{
		RoomInfos: room.FormatInternalRooms2Protobuf(rooms),
	}
}

func getDailyCoin(logger logx.Logger, in *twoonone_pb.GetDailyCoinRequest) (*twoonone_pb.GetDailyCoinResponse, *types.AppError) {
	u, err := db.GetUser(logger, in.UserId)
	if err != nil {
		return nil, err
	}
	if checkDailyCoin(u.UserTwoonone) {
		if err := db.UpdateUser(logger, u.UserTwoonone.Id, sql.IncCoin(500), sql.UpdateGetDailyTime()); err != nil {
			return nil, err
		}
		return &twoonone_pb.GetDailyCoinResponse{}, nil
	} else {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_ALREADY_GET_DALIY_COIN, "")
	}
}

func checkDailyCoin(u twoonone.UserTwoonone) bool {
	if u.LastGetdaliyTime.IsZero() { //第一次领取
		return true
	} else {
		now := time.Now()
		last := u.LastGetdaliyTime
		// 判断是否为同一年同一月
		if now.Month() == last.Month() && now.Year() == last.Year() {
			return now.Day() > last.Day()
		} else {
			//非同一年同一月必然不是同一天
			return true
		}
	}
}

func getRoom(logger logx.Logger, in *twoonone_pb.GetRoomRequest) (*twoonone_pb.GetRoomResponse, *types.AppError) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	return &twoonone_pb.GetRoomResponse{
		RoomInfo: room.FormatInternalRoom2Protobuf(r),
	}, nil
}

func createRoom(in *twoonone_pb.CreateRoomRequest) (*twoonone_pb.CreateRoomResponse, *types.AppError) {
	newRoom := room.New(200, 1)
	rooms = append(rooms, newRoom)
	return &twoonone_pb.CreateRoomResponse{
		RoomHash: newRoom.GetHash(),
	}, nil
}

func joinRoom(logger logx.Logger, in *twoonone_pb.JoinRoomRequest) (*twoonone_pb.JoinRoomResponse, *types.AppError) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	if _, ok := getPlayerFromRooms(in.UserId); ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_EXISTED_A_ROOM, "")
	}
	u, err := db.GetUser(logger, in.UserId)
	if err != nil {
		return nil, err
	}
	if serr := r.Join(logger, u.UserPublic.Id, u.Name, u.Coin); serr != nil {
		return nil, err
	}
	return &twoonone_pb.JoinRoomResponse{}, nil
}

func exitRoom(logger logx.Logger, in *twoonone_pb.ExitRoomRequest) (*twoonone_pb.ExitRoomResponse, *types.AppError) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	if err := r.Exit(logger, in.UserId); err != nil {
		return nil, err
	}
	return &twoonone_pb.ExitRoomResponse{}, nil
}

func startRoom(logger logx.Logger, in *twoonone_pb.StartRoomRequest) (*twoonone_pb.StartRoomResponse, *types.AppError) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	if _, ok := r.GetPlayer(in.UserId); !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST, "")
	}
	if err := r.Start(logger); err != nil {
		return nil, err
	}
	return &twoonone_pb.StartRoomResponse{}, nil
}

func robLandowner(logger logx.Logger, in *twoonone_pb.RobLandownerRequest) (*twoonone_pb.RobLandownerResponse, *types.AppError) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	p, ok := r.GetPlayer(in.UserId)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST, "")
	}
	if err := r.RobLandowner(logger, p); err != nil {
		return nil, err
	}
	return &twoonone_pb.RobLandownerResponse{}, nil
}

func noRobLandowner(logger logx.Logger, in *twoonone_pb.NoRobLandownerRequest) (*twoonone_pb.NoRobLandownerResponse, *types.AppError) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	p, ok := r.GetPlayer(in.UserId)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST, "")
	}
	if err := r.NoRobLandowner(logger, p); err != nil {
		return nil, err
	}
	return &twoonone_pb.NoRobLandownerResponse{}, nil
}

func sendCard(logger logx.Logger, in *twoonone_pb.SendCardRequest) (*twoonone_pb.SendCardResponse, *types.AppError) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	p, ok := r.GetPlayer(in.UserId)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST, "")
	}
	if err := r.SendCard(logger, p, in.Sendcards); err != nil {
		return nil, err
	}
	return &twoonone_pb.SendCardResponse{}, nil
}

func noSendCard(logger logx.Logger, in *twoonone_pb.NoSendCardRequest) (*twoonone_pb.NoSendCardResponse, *types.AppError) {
	r, ok := findRoom(in.RoomHash)
	if !ok {
		return nil, types.NewError(twoonone_pb.Error_ERROR_ROOM_NO_EXIST, "")
	}
	p, ok := r.GetPlayer(in.UserId)
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
