// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package types

type Card struct {
	Number int `form:"number,range=[0:14]"`
}

type ExitRoomRequest struct {
	RoomHash string `path:"room_hash"`
	Extra    Extra  `custom:"extra"`
}

type Extra struct {
	UserId       string  `custom:"user_id"`
	Name         string  `custom:"name"`
	Email        string  `custom:"email"`
	WinCount     int     `custom:"wincount"`
	LoseCount    int     `custom:"losecount"`
	Coin         float64 `custom:"coin"`
	Extra_update bool    `custom:"extra_update"`
}

type GetDailyCoinRequest struct {
	Extra Extra `custom:"extra"`
}

type GetRoomRequest struct {
	RoomHash string `path:"room_hash"`
	Extra    Extra  `custom:"extra"`
}

type JoinRoomRequest struct {
	RoomHash string `path:"room_hash"`
	Extra    Extra  `custom:"extra"`
}

type NoRobLandownerRequest struct {
	RoomHash string `path:"room_hash"`
	Extra    Extra  `custom:"extra"`
}

type NoSendCardRequest struct {
	RoomHash string `path:"room_hash"`
	Extra    Extra  `custom:"extra"`
}

type RobLandownerRequest struct {
	RoomHash string `path:"room_hash"`
	Extra    Extra  `custom:"extra"`
}

type SendCardRequest struct {
	RoomHash  string `path:"room_hash"`
	SendCards []Card `form:"sendcards"`
	Extra     Extra  `custom:"extra"`
}

type StartRoomRequest struct {
	RoomHash string `path:"room_hash"`
	Extra    Extra  `custom:"extra"`
}
