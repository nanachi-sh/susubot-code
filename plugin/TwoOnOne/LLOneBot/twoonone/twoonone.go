package twoonone

import (
	"database/sql"
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/db"
	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/protos/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/twoonone/player"
	"github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/twoonone/room"
)

var (
	rooms []*room.Room
)

func getRoom(hash string) *room.Room {
	for _, v := range rooms {
		if v.GetHash() == hash {
			return v
		}
	}
	return nil
}

func getAccount(id string) (*twoonone_pb.PlayerAccountInfo, error) {
	pi, err := db.GetPlayer(id)
	if err != nil {
		return nil, err
	}
	return pi, nil
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

func GetRooms() []*twoonone_pb.RoomInfo {
	rs := []*twoonone_pb.RoomInfo{}
	for _, v := range rooms {
		rs = append(rs, insideRoomToRoom(v))
	}
	return rs
}

func CreateAccount(req *twoonone_pb.CreateAccountRequest) (*twoonone_pb.Errors, error) {
	ai, _ := db.GetPlayer(req.PlayerId)
	if ai != nil {
		return twoonone_pb.Errors_PlayerAccountExist.Enum(), nil
	}
	if err := db.CreateAccount(req.PlayerId, req.PlayerName, 0); err != nil {
		return nil, err
	}
	serr, err := getDailyCoin(req.PlayerId)
	if err != nil {
		return nil, err
	}
	return serr, nil
}

func GetAccount(req *twoonone_pb.GetAccountRequest) (*twoonone_pb.PlayerAccountInfo, *twoonone_pb.Errors, error) {
	ai, err := getAccount(req.PlayerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, twoonone_pb.Errors_PlayerNoExist.Enum(), nil
		}
		return nil, nil, err
	}
	return ai, nil, nil
}

func getDailyCoin(id string) (*twoonone_pb.Errors, error) {
	ai, err := db.GetPlayer(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return twoonone_pb.Errors_PlayerNoExist.Enum(), nil
		}
		return nil, err
	}
	if ai.LastGetDailyTimestamp == 0 { //第一次领取
		if err := db.IncCoin(ai.Id, 1000); err != nil {
			return nil, err
		}
		if err := db.UpdateLastGetDailyTimestamp(ai.Id, time.Time{}); err != nil {
			return nil, err
		}
	} else if getDailyCoinIF(time.Unix(ai.LastGetDailyTimestamp, 0)) { //符合条件
		if err := db.IncCoin(ai.Id, 500); err != nil {
			return nil, err
		}
		if err := db.UpdateLastGetDailyTimestamp(ai.Id, time.Time{}); err != nil {
			return nil, err
		}
		return nil, nil
	} else { //不符合
		return twoonone_pb.Errors_PlayerAlreadyGetDailyCoin.Enum(), nil
	}
	return nil, nil
}

func GetDailyCoin(req *twoonone_pb.GetDailyCoinRequest) (*twoonone_pb.BasicResponse, error) {
	serr, err := getDailyCoin(req.PlayerId)
	if err != nil {
		return nil, err
	}
	return &twoonone_pb.BasicResponse{
		Err: serr,
	}, nil
}

func insideRoomToRoom(r *room.Room) *twoonone_pb.RoomInfo {
	if r == nil {
		return nil
	}
	ps := []*twoonone_pb.PlayerInfo{}
	for _, v := range r.GetPlayers() {
		ps = append(ps, insidePlayerToPlayerInfo(v))
	}
	cps := []*twoonone_pb.SendCard{}
	for _, v := range r.GetCardPool() {
		cps = append(cps, &twoonone_pb.SendCard{
			SenderInfo:        insidePlayerToPlayerInfo(v.SenderInfo),
			SendCards:         v.SendCards,
			SendCardType:      v.SendCardType,
			SendCardSize:      int32(v.SendCardSize),
			SendCardContinous: int32(v.SendCardContinous),
		})
	}
	operatorNow := insidePlayerToPlayerInfo(r.GetOperatorNow())
	lo := r.GetLandowner()
	landowner := insidePlayerToPlayerInfo(lo)
	var loCards []twoonone_pb.Card
	if lo != nil {
		loCards = lo.GetCards()
	}
	var farmers []*twoonone_pb.PlayerInfo
	for _, v := range r.GetFarmers() {
		if v == nil {
			continue
		}
		farmers = append(farmers, insidePlayerToPlayerInfo(v))
	}
	return &twoonone_pb.RoomInfo{
		Hash:           r.GetHash(),
		Players:        ps,
		BasicCoin:      r.GetBasicCoin(),
		Multiple:       int32(r.GetMultiple()),
		Stage:          r.GetStage(),
		CardPool:       cps,
		OperatorNow:    operatorNow,
		LandownerCards: loCards,
		Landowner:      landowner,
		Farmers:        farmers,
	}
}

func insidePlayerToPlayerInfo(p *player.Player) *twoonone_pb.PlayerInfo {
	if p == nil {
		return nil
	}
	pi := &twoonone_pb.PlayerInfo{
		AccountInfo: &twoonone_pb.PlayerAccountInfo{
			Id:                    p.GetId(),
			Name:                  p.GetName(),
			WinCount:              int32(p.GetWinCount()),
			LoseCount:             int32(p.GetLoseCount()),
			Coin:                  p.GetCoin(),
			LastGetDailyTimestamp: p.GetLastGetDailyTimestamp(),
		},
	}
	if p.GetRoomHash() != "" {
		pi.TableInfo = &twoonone_pb.PlayerTableInfo{
			RoomHash:           p.GetRoomHash(),
			Cards:              p.GetCards(),
			RobLandownerAction: p.GetRobLandownerAction(),
		}
	}
	return pi
}

func GetRoom(req *twoonone_pb.GetRoomInfoRequest) *twoonone_pb.GetRoomInfoResponse {
	switch {
	case req.RoomHash != nil:
		if r := getRoom(*req.RoomHash); r == nil {
			return &twoonone_pb.GetRoomInfoResponse{
				Err:  twoonone_pb.Errors_RoomNoExist.Enum(),
				Info: nil,
			}
		} else {
			return &twoonone_pb.GetRoomInfoResponse{
				Info: insideRoomToRoom(r),
			}
		}
	case req.PlayerId != nil:
		p, ok := getPlayerFromRooms(*req.PlayerId)
		if !ok {
			return &twoonone_pb.GetRoomInfoResponse{
				Err:  twoonone_pb.Errors_PlayerNoExistAnyRoom.Enum(),
				Info: nil,
			}
		}
		if r := getRoom(p.GetRoomHash()); r == nil {
			return &twoonone_pb.GetRoomInfoResponse{
				Err:  twoonone_pb.Errors_RoomNoExist.Enum(),
				Info: nil,
			}
		} else {
			return &twoonone_pb.GetRoomInfoResponse{
				Info: insideRoomToRoom(r),
			}
		}
	default:
		return &twoonone_pb.GetRoomInfoResponse{
			Err:  twoonone_pb.Errors_Unexpected.Enum(),
			Info: nil,
		}
	}
}

func getDailyCoinIF(last time.Time) bool {
	now := time.Now()
	// 判断是否为同一年同一月
	if now.Month() == last.Month() && now.Year() == last.Year() {
		return now.Day() > last.Day()
	} else {
		//非同一年同一月必然不是同一天
		return true
	}
}

func CreateRoom(req *twoonone_pb.CreateRoomRequest) *twoonone_pb.CreateRoomResponse {
	newRoom := room.New(req.BasicCoin, int(req.InitialMultiple))
	rooms = append(rooms, newRoom)
	return &twoonone_pb.CreateRoomResponse{
		RoomHash: newRoom.GetHash(),
	}
}

func JoinRoom(req *twoonone_pb.JoinRoomRequest) (*twoonone_pb.JoinRoomResponse, error) {
	r := getRoom(req.RoomHash)
	if r == nil {
		return &twoonone_pb.JoinRoomResponse{
			Err: twoonone_pb.Errors_RoomNoExist.Enum(),
		}, nil
	}
	if _, ok := getPlayerFromRooms(req.PlayerId); ok {
		return &twoonone_pb.JoinRoomResponse{
			Err: twoonone_pb.Errors_RoomExistPlayer.Enum(),
		}, nil
	}
	pi, err := getAccount(req.PlayerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return &twoonone_pb.JoinRoomResponse{
				Err: twoonone_pb.Errors_PlayerNoExist.Enum(),
			}, nil
		}
		return nil, err
	}
	if err := r.Join(pi); err != nil {
		return &twoonone_pb.JoinRoomResponse{
			Err: err,
		}, nil
	}
	players := r.GetPlayers()
	resp_players := []*twoonone_pb.PlayerInfo{}
	for _, v := range players {
		resp_players = append(resp_players, insidePlayerToPlayerInfo(v))
	}
	return &twoonone_pb.JoinRoomResponse{
		RoomPlayers: resp_players,
	}, nil
}

func ExitRoom(req *twoonone_pb.ExitRoomRequest) (*twoonone_pb.ExitRoomResponse, error) {
	pi, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &twoonone_pb.ExitRoomResponse{
			Err: twoonone_pb.Errors_PlayerNoExist.Enum(),
		}, nil
	}
	r := getRoom(pi.GetRoomHash())
	if err := r.Exit(req.PlayerId); err != nil {
		return &twoonone_pb.ExitRoomResponse{Err: err}, nil
	}
	players := r.GetPlayers()
	resp_players := []*twoonone_pb.PlayerInfo{}
	for _, v := range players {
		resp_players = append(resp_players, insidePlayerToPlayerInfo(v))
	}
	return &twoonone_pb.ExitRoomResponse{
		RoomPlayers: resp_players,
	}, nil
}

func StartRoom(req *twoonone_pb.StartRoomRequest) (*twoonone_pb.StartRoomResponse, error) {
	pi, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &twoonone_pb.StartRoomResponse{
			Err: twoonone_pb.Errors_PlayerNoExist.Enum(),
		}, nil
	}
	r := getRoom(pi.GetRoomHash())
	lastPi, err := r.Start()
	if err != nil {
		return &twoonone_pb.StartRoomResponse{Err: err}, nil
	}
	return &twoonone_pb.StartRoomResponse{
		LastOperator: insidePlayerToPlayerInfo(lastPi),
	}, nil
}

func RobLandownerAction(req *twoonone_pb.RobLandownerActionRequest) (*twoonone_pb.RobLandownerActionResponse, error) {
	p, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &twoonone_pb.RobLandownerActionResponse{
			Err: twoonone_pb.Errors_PlayerNoExist.Enum(),
		}, nil
	}
	r := getRoom(p.GetRoomHash())
	last, err := r.RobLandownerAction(p, req.Action)
	if err != nil {
		return &twoonone_pb.RobLandownerActionResponse{Err: err}, nil
	}
	sendcard := false
	if r.GetStage() == twoonone_pb.RoomStage_SendingCards {
		sendcard = true
	}
	return &twoonone_pb.RobLandownerActionResponse{
		LastOperator:    insidePlayerToPlayerInfo(last),
		IntoSendingCard: sendcard,
	}, nil
}

func SendCardAction(req *twoonone_pb.SendCardRequest) (*twoonone_pb.SendCardResponse, error) {
	p, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &twoonone_pb.SendCardResponse{
			Err: twoonone_pb.Errors_PlayerNoExist.Enum(),
		}, nil
	}
	r := getRoom(p.GetRoomHash())
	next, e, serr, err := r.SendCardAction(p, req.SendCards, req.Action)
	if err != nil {
		return nil, err
	}
	var cn *int32
	if e.CardNumber != nil {
		cn = new(int32)
		*cn = int32(*e.CardNumber)
	}
	return &twoonone_pb.SendCardResponse{
		Err:                    serr,
		SenderCard:             p.GetCards(),
		NextPlayer:             insidePlayerToPlayerInfo(next),
		SenderCardNumberNotice: e.SenderCardNumberNotice,
		GameFinish:             e.GameFinish,
		SenderCardTypeNotice:   e.SenderCardTypeNotice,
		SenderCardNumber:       cn,
		GameFinishE:            e.GameFinishE,
		SenderCardType:         e.CardType,
	}, nil
}
