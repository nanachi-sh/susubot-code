package uno

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"

	"github.com/nanachi-sh/susubot-code/plugin/uno/db"
	"github.com/nanachi-sh/susubot-code/plugin/uno/define"
	uno "github.com/nanachi-sh/susubot-code/plugin/uno/protos/qqverifier"
	uno_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
	"github.com/nanachi-sh/susubot-code/plugin/uno/uno/player"
	"github.com/nanachi-sh/susubot-code/plugin/uno/uno/room"
	"github.com/twmb/murmur3"
	"google.golang.org/grpc"
)

var (
	rooms []*room.Room
)

// type GameEvents struct {
// 	GameFinish            *uno_pb.RoomEventResponse_GameFinishEvent
// 	DrawCard_IntoSendCard *uno_pb.RoomEventResponse_DrawCard_IntoSendCardEvent
// 	DrawCard_Skipped      *uno_pb.RoomEventResponse_DrawCard_SkippedEvent
// 	IndicateUNO_Success   *uno_pb.RoomEventResponse_IndicateUNO_SuccessEvent
// 	HandCardUpdate        *uno_pb.RoomEventResponse_HandCardUpdate
// }

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

func playerHash(id string, roomHash string) string {
	h1, h2 := murmur3.SeedStringSum128(rand.Uint64(), rand.Uint64(), id+roomHash)
	return fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}

func CreateRoom(cs []*http.Cookie) (*uno_pb.CreateRoomResponse, error) {
	if len(cs) == 0 {
		return &uno_pb.CreateRoomResponse{Err: uno_pb.Errors_NoFoundAccountHash.Enum()}, nil
	}
	isNormal, err := CheckNormalUserFromSource(cs)
	if err != nil {
		if err == sql.ErrNoRows {
			return &uno_pb.CreateRoomResponse{Err: uno_pb.Errors_NoValidAccountHash.Enum()}, nil
		}
		return nil, err
	}
	if !isNormal && !CheckPrivilegeUser(cs) {
		return &uno_pb.CreateRoomResponse{Err: uno_pb.Errors_NoValidAccountHash.Enum()}, nil
	}
	newRoom := room.New()
	rooms = append(rooms, newRoom)
	ctx, cancel := context.WithCancel(context.Background())
	roomEvents = append(roomEvents, &roomEvent{
		roomHash: newRoom.GetHash(),
		ctx:      ctx,
		cancel:   cancel,
		block:    sync.RWMutex{},
		wait:     sync.RWMutex{},
	})
	return &uno_pb.CreateRoomResponse{
		RoomHash: newRoom.GetHash(),
	}, nil
}

func JoinRoom(cs []*http.Cookie, req *uno_pb.JoinRoomRequest) (*uno_pb.JoinRoomResponse, error) {
	isTemp := CheckTempUser(req.PlayerInfo.Id)
	isPrivilege := CheckPrivilegeUser(cs)
	isNormal, err := CheckNormalUserFromSource(cs)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	if isTemp || isPrivilege || isNormal {
	} else {
		return &uno_pb.JoinRoomResponse{Err: uno_pb.Errors_AbnormalAccount.Enum()}, nil
	}
	r, ok := getRoom(req.RoomHash)
	if !ok {
		return &uno_pb.JoinRoomResponse{
			Err: uno_pb.Errors_RoomNoExist.Enum(),
		}, nil
	}
	if _, ok := getPlayerFromRooms(req.PlayerInfo.Id); ok {
		return &uno_pb.JoinRoomResponse{
			Err: uno_pb.Errors_RoomExistPlayer.Enum(),
		}, nil
	}
	hash := playerHash(req.PlayerInfo.Id, req.RoomHash)
	if serr := r.Join(&uno_pb.PlayerAccountInfo{
		Id:   req.PlayerInfo.Id,
		Name: req.PlayerInfo.Name,
	}, hash); serr != nil {
		return &uno_pb.JoinRoomResponse{
			Err: serr,
		}, nil
	}
	ps := []*uno_pb.PlayerAccountInfo{}
	for _, v := range r.GetPlayers() {
		ps = append(ps, v.FormatToProtoBuf().PlayerAccountInfo)
	}
	return &uno_pb.JoinRoomResponse{
		Players:    ps,
		VerifyHash: hash,
	}, nil
}

func GetRooms() *uno_pb.GetRoomsResponse {
	rs := []*uno_pb.RoomSimple{}
	for _, v := range rooms {
		rs = append(rs, v.FormatToProtoBufSimple())
	}
	return &uno_pb.GetRoomsResponse{
		Rooms: rs,
	}
}

func GetRoom(cs []*http.Cookie, req *uno_pb.GetRoomRequest) *uno_pb.GetRoomResponse {
	isPrivilege := CheckPrivilegeUser(cs)
	if r, ok := getRoom(req.RoomHash); !ok {
		return &uno_pb.GetRoomResponse{
			Err: uno_pb.Errors_RoomNoExist.Enum(),
		}
	} else {
		if isPrivilege {
			return &uno_pb.GetRoomResponse{Extra: r.FormatToProtoBufExtra()}
		} else {
			return &uno_pb.GetRoomResponse{Simple: r.FormatToProtoBufSimple()}
		}
	}
}

func ExitRoom(cs []*http.Cookie, req *uno_pb.ExitRoomRequest) *uno_pb.ExitRoomResponse {
	playerHash, ok := GetPlayerHash(cs)
	if !ok {
		return &uno_pb.ExitRoomResponse{
			Err: uno_pb.Errors_NoFoundPlayerHash.Enum(),
		}
	}
	r, ok := getRoom(req.RoomHash)
	if !ok {
		return &uno_pb.ExitRoomResponse{
			Err: uno_pb.Errors_RoomNoExist.Enum(),
		}
	}
	p, ok := r.GetPlayer(req.PlayerId)
	if !ok {
		return &uno_pb.ExitRoomResponse{
			Err: uno_pb.Errors_RoomNoExistPlayer.Enum(),
		}
	}
	if p.GetHash() != playerHash {
		return &uno_pb.ExitRoomResponse{
			Err: uno_pb.Errors_PlayerHashNE.Enum(),
		}
	}
	if serr := r.Exit(req.PlayerId); serr != nil {
		return &uno_pb.ExitRoomResponse{Err: serr}
	}
	ps := []*uno_pb.PlayerAccountInfo{}
	for _, v := range r.GetPlayers() {
		ps = append(ps, v.FormatToProtoBuf().PlayerAccountInfo)
	}
	return &uno_pb.ExitRoomResponse{
		Players: ps,
	}
}

func StartRoom(cs []*http.Cookie, req *uno_pb.StartRoomRequest) *uno_pb.BasicResponse {
	playerHash, ok := GetPlayerHash(cs)
	if !ok {
		return &uno_pb.BasicResponse{Err: uno_pb.Errors_NoFoundPlayerHash.Enum()}
	}
	r, ok := getRoom(req.RoomHash)
	if !ok {
		return &uno_pb.BasicResponse{
			Err: uno_pb.Errors_Unexpected.Enum(),
		}
	}
	if _, ok := r.GetPlayerFromHash(playerHash); !ok {
		return &uno_pb.BasicResponse{Err: uno_pb.Errors_NoValidPlayerHash.Enum()}
	}
	if serr := r.Start(); serr != nil {
		return &uno_pb.BasicResponse{Err: serr}
	}
	return &uno_pb.BasicResponse{}
}

func DrawCard(cs []*http.Cookie, req *uno_pb.DrawCardRequest) *uno_pb.DrawCardResponse {
	playerHash, ok := GetPlayerHash(cs)
	if !ok {
		return &uno_pb.DrawCardResponse{Err: uno_pb.Errors_NoFoundPlayerHash.Enum()}
	}
	p, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &uno_pb.DrawCardResponse{
			Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
		}
	}
	if p.GetHash() != playerHash {
		return &uno_pb.DrawCardResponse{Err: uno_pb.Errors_PlayerHashNE.Enum()}
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
		ge, ok := findRoomEvent(r.GetHash())
		if !ok {
			return &uno_pb.DrawCardResponse{Err: uno_pb.Errors_Unexpected.Enum()}
		}
		switch {
		case e.IntoSendCard:
			ge.update(&uno_pb.RoomEventResponse{
				DrawCard_IntoSendCard: e.IntoSendCardE,
			})
		case e.Skipped:
			ge.update(&uno_pb.RoomEventResponse{
				DrawCard_Skipped: e.SkippedE,
			})
		}
	}
	switch stage := r.GetStage(); stage {
	case uno_pb.Stage_ElectingBanker:
		return &uno_pb.DrawCardResponse{
			ElectingBanker: &uno_pb.DrawCardResponse_DrawCard_ElectingBanker{
				ElectBankerCard: p.GetElectBankerCard(),
			},
		}
	case uno_pb.Stage_SendingCard:
		cs := []*uno_pb.Card{}
		for _, v := range p.GetCards() {
			cs = append(cs, &v)
		}
		return &uno_pb.DrawCardResponse{
			SendingCard: &uno_pb.DrawCardResponse_DrawCard_SendingCard{
				PlayerCard: cs,
				DrawCard:   p.GetDrawCard(),
			},
		}
	default:
		return &uno_pb.DrawCardResponse{Err: uno_pb.Errors_Unexpected.Enum()}
	}
}

func SendCard(cs []*http.Cookie, req *uno_pb.SendCardRequest) *uno_pb.SendCardResponse {
	if req.SendCard == nil {
		return &uno_pb.SendCardResponse{Err: uno_pb.Errors_Unexpected.Enum()}
	}
	playerHash, ok := GetPlayerHash(cs)
	if !ok {
		return &uno_pb.SendCardResponse{Err: uno_pb.Errors_NoFoundPlayerHash.Enum()}
	}
	r, ok := getRoom(req.RoomHash)
	if !ok {
		return &uno_pb.SendCardResponse{Err: uno_pb.Errors_RoomNoExist.Enum()}
	}
	p, ok := r.GetPlayer(req.PlayerId)
	if !ok {
		return &uno_pb.SendCardResponse{
			Err: uno_pb.Errors_RoomNoExistPlayer.Enum(),
		}
	}
	if p.GetHash() != playerHash {
		return &uno_pb.SendCardResponse{Err: uno_pb.Errors_PlayerHashNE.Enum()}
	}
	sc := uno_pb.Card{}
	if req.SendCard != nil {
		sc = *req.SendCard
	}
	next, e, serr := r.SendCardAction(p, sc, uno_pb.SendCardActions_Send)
	if serr != nil {
		return &uno_pb.SendCardResponse{Err: serr}
	}
	resp := new(uno_pb.SendCardResponse)
	if next != nil {
		resp.NextOperator = next.FormatToProtoBuf().PlayerAccountInfo
	}
	if e != nil {
		ge, ok := findRoomEvent(r.GetHash())
		if !ok {
			return &uno_pb.SendCardResponse{Err: uno_pb.Errors_Unexpected.Enum()}
		}
		if e.GameFinish {
			if !deleteRoom(r) {
				return &uno_pb.SendCardResponse{Err: uno_pb.Errors_Unexpected.Enum()}
			}
			ge.update(&uno_pb.RoomEventResponse{
				GameFinish: e.GameFinishE,
			})
		}
	}
	return resp
}

func NoSendCard(cs []*http.Cookie, req *uno_pb.NoSendCardRequest) *uno_pb.NoSendCardResponse {
	playerHash, ok := GetPlayerHash(cs)
	if !ok {
		return &uno_pb.NoSendCardResponse{Err: uno_pb.Errors_NoFoundPlayerHash.Enum()}
	}
	r, ok := getRoom(req.RoomHash)
	if !ok {
		return &uno_pb.NoSendCardResponse{Err: uno_pb.Errors_RoomNoExist.Enum()}
	}
	p, ok := r.GetPlayer(req.PlayerId)
	if !ok {
		return &uno_pb.NoSendCardResponse{
			Err: uno_pb.Errors_RoomNoExistPlayer.Enum(),
		}
	}
	if p.GetHash() != playerHash {
		return &uno_pb.NoSendCardResponse{Err: uno_pb.Errors_PlayerHashNE.Enum()}
	}
	next, e, serr := r.SendCardAction(p, uno_pb.Card{}, uno_pb.SendCardActions_NoSend)
	if serr != nil {
		return &uno_pb.NoSendCardResponse{Err: serr}
	}
	resp := new(uno_pb.NoSendCardResponse)
	if next != nil {
		resp.NextOperator = next.FormatToProtoBuf().PlayerAccountInfo
	}
	if e != nil {
		ge, ok := findRoomEvent(r.GetHash())
		if !ok {
			return &uno_pb.NoSendCardResponse{Err: uno_pb.Errors_Unexpected.Enum()}
		}
		if e.GameFinish {
			if !deleteRoom(r) {
				return &uno_pb.NoSendCardResponse{Err: uno_pb.Errors_Unexpected.Enum()}
			}
			ge.update(&uno_pb.RoomEventResponse{
				GameFinish: e.GameFinishE,
			})
		}
	}
	return resp
}

func CallUNO(cs []*http.Cookie, req *uno_pb.CallUNORequest) *uno_pb.CallUNOResponse {
	playerHash, ok := GetPlayerHash(cs)
	if !ok {
		return &uno_pb.CallUNOResponse{
			Err: uno_pb.Errors_NoFoundPlayerHash.Enum(),
		}
	}
	r, ok := getRoom(req.RoomHash)
	if !ok {
		return &uno_pb.CallUNOResponse{
			Err: uno_pb.Errors_RoomNoExist.Enum(),
		}
	}
	p, ok := r.GetPlayer(req.PlayerId)
	if !ok {
		return &uno_pb.CallUNOResponse{
			Err: uno_pb.Errors_RoomNoExistPlayer.Enum(),
		}
	}
	if p.GetHash() != playerHash {
		return &uno_pb.CallUNOResponse{
			Err: uno_pb.Errors_PlayerHashNE.Enum(),
		}
	}
	cards, serr := r.CallUNO(p)
	cardPs := []*uno_pb.Card{}
	for _, v := range cards {
		cardPs = append(cardPs, &v)
	}
	return &uno_pb.CallUNOResponse{
		PlayerCard: cardPs,
		Err:        serr,
	}
}

func Challenge(cs []*http.Cookie, req *uno_pb.ChallengeRequest) *uno_pb.ChallengeResponse {
	playerHash, ok := GetPlayerHash(cs)
	if !ok {
		return &uno_pb.ChallengeResponse{
			Err: uno_pb.Errors_NoFoundPlayerHash.Enum(),
		}
	}
	r, ok := getRoom(req.RoomHash)
	if !ok {
		return &uno_pb.ChallengeResponse{
			Err: uno_pb.Errors_RoomNoExist.Enum(),
		}
	}
	p, ok := r.GetPlayer(req.PlayerId)
	if !ok {
		return &uno_pb.ChallengeResponse{
			Err: uno_pb.Errors_RoomNoExistPlayer.Enum(),
		}
	}
	if p.GetHash() != playerHash {
		return &uno_pb.ChallengeResponse{
			Err: uno_pb.Errors_PlayerHashNE.Enum(),
		}
	}
	win, pi, serr := r.Challenge(p)
	if serr != nil {
		return &uno_pb.ChallengeResponse{Err: serr}
	}
	if win {
		ge, ok := findRoomEvent(r.GetHash())
		if !ok {
			return &uno_pb.ChallengeResponse{Err: uno_pb.Errors_Unexpected.Enum()}
		}
		ge.update(&uno_pb.RoomEventResponse{
			HandCardUpdate: &uno_pb.RoomEventResponse_HandCardUpdateEvent{Updated: pi.PlayerAccountInfo},
		})
	}
	return &uno_pb.ChallengeResponse{
		Err:   serr,
		IsWin: win,
	}
}

func IndicateUNO(cs []*http.Cookie, req *uno_pb.IndicateUNORequest) *uno_pb.IndicateUNOResponse {
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
	if ok {
		ge, ok := findRoomEvent(r.GetHash())
		if !ok {
			return &uno_pb.IndicateUNOResponse{Err: uno_pb.Errors_Unexpected.Enum()}
		}
		ge.update(&uno_pb.RoomEventResponse{
			IndicateUNO_Success: &uno_pb.RoomEventResponse_IndicateUNO_SuccessEvent{
				Punished: pi.PlayerAccountInfo,
			},
		})
	}
	return &uno_pb.IndicateUNOResponse{
		IndicateSuccessed: ok,
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

func TEST_SetPlayerCard(cs []*http.Cookie, req *uno_pb.TEST_SetPlayerCardRequest) *uno_pb.BasicResponse {
	if !CheckPrivilegeUser(cs) {
		return &uno_pb.BasicResponse{Err: uno_pb.Errors_NoPrivilegeAccount.Enum()}
	}
	p, ok := getPlayerFromRooms(req.PlayerId)
	if !ok {
		return &uno_pb.BasicResponse{
			Err: uno_pb.Errors_PlayerNoExistAnyRoom.Enum(),
		}
	}
	p.ClearHandCard()
	cards := []uno_pb.Card{}
	for _, v := range req.Cards {
		cards = append(cards, *v)
	}
	p.AddCards(cards)
	return nil
}

func CreateUser(cs []*http.Cookie, req *uno_pb.CreateUserRequest) (*uno_pb.BasicResponse, error) {
	if req.UserInfo == nil || (req.UserInfo.Id == "" || req.UserInfo.Name == "") {
		return nil, errors.New("Id或名字不能为空")
	}
	if CheckPrivilegeUser(cs) {
		// 特权用户无需验证
		if err := db.CreateUser(req.UserInfo.Id, req.UserInfo.Name, req.Source); err != nil {
			return nil, err
		}
	} else {
		if req.VerifyHash == "" {
			return nil, errors.New("VerifyHash不能为空")
		}
		switch req.Source {
		case uno_pb.Source_QQ:
			resp, err := define.QQVerifierC.Verified(define.QQVerifierCtx, &uno.VerifiedRequest{
				VerifyHash: req.VerifyHash,
			})
			if err != nil {
				return nil, err
			}
			if resp.Err != nil {
				switch *resp.Err {
				case uno.Errors_VerifyNoFound:
					return nil, errors.New("验证哈希不正确")
				case uno.Errors_Expired:
					return nil, errors.New("验证请求已过期")
				case uno.Errors_UnVerified:
					return nil, errors.New("还未验证")
				}
			}
			if *resp.Result != uno.Result_Verified {
				return nil, errors.New("还未验证")
			}
			if req.UserInfo.Id != resp.VarifyId {
				return nil, errors.New("验证请求的QQID与申请请求的QQID不符")
			}
		}
	}
	return nil, errors.New("非预期错误")
}

func RoomEvent(req *uno_pb.RoomEventRequest, stream grpc.ServerStreamingServer[uno_pb.RoomEventResponse]) (*uno_pb.RoomEventResponse, error) {
	if req.PlayerHash == "" {
		return &uno_pb.RoomEventResponse{Err: uno_pb.Errors_NoFoundPlayerHash.Enum()}, nil
	}
	r, ok := getRoom(req.RoomHash)
	if !ok {
		return &uno_pb.RoomEventResponse{Err: uno_pb.Errors_RoomNoExist.Enum()}, nil
	}
	if _, ok := r.GetPlayerFromHash(req.PlayerHash); !ok {
		return &uno_pb.RoomEventResponse{Err: uno_pb.Errors_NoValidPlayerHash.Enum()}, nil
	}
	e, ok := findRoomEvent(req.RoomHash)
	if !ok {
		return &uno_pb.RoomEventResponse{Err: uno_pb.Errors_Unexpected.Enum()}, nil
	}
	for {
		re := e.read()
		if err := stream.Send(re); err != nil {
			return nil, err
		}
	}
}

var roomEvents []*roomEvent

type roomEvent struct {
	roomHash string
	ctx      context.Context
	cancel   context.CancelFunc
	block    sync.RWMutex
	wait     sync.RWMutex
}
type myRoomEvent struct{}

func findRoomEvent(hash string) (*roomEvent, bool) {
	for _, v := range roomEvents {
		if v.roomHash == hash {
			return v, true
		}
	}
	return nil, false
}

func deleteRoomEvent(hash string) {
	for i, v := range roomEvents {
		if v.roomHash == hash {
			if len(roomEvents) == 1 {
				roomEvents = []*roomEvent{}
			} else {
				roomEvents = append(roomEvents[:i], roomEvents[i+1:]...)
			}
		}
	}
}

func (re *roomEvent) update(e *uno_pb.RoomEventResponse) {
	if e == nil {
		return
	}
	re.ctx = context.WithValue(re.ctx, myRoomEvent{}, e)
	re.cancel()
	if e.GameFinish != nil {
		go func() {
			re.wait.RLock()
			deleteRoomEvent(re.roomHash)
		}()
	}
}

func (re *roomEvent) read() *uno_pb.RoomEventResponse {
	//若等待队列关闭，则加入并阻塞
	re.wait.RLock()
	re.wait.RUnlock()
	//进入阻塞队列
	re.block.RLock()
	//检查阻塞队列是否为空
	defer func() {
		if re.block.TryLock() { //阻塞队列为空
			//打开等待队列
			re.wait.Unlock()
			//重置
			re.readReset()
			//等待Wait队列空
			re.wait.Lock()
			re.wait.Unlock()
			//打开阻塞队列
			re.block.Unlock()
		}
	}()
	defer re.block.RUnlock()
	<-re.ctx.Done()
	//第一个会话负责关闭等待队列
	re.wait.TryLock()
	v := re.ctx.Value(myRoomEvent{})
	return v.(*uno_pb.RoomEventResponse)
}

func (re *roomEvent) readReset() {
	re.ctx = context.WithoutCancel(re.ctx)
	re.ctx, re.cancel = context.WithCancel(re.ctx)
}
