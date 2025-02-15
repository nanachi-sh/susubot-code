// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package types

type Card struct {
	Number int `json:"number,range=[0:14]"`
}

type ExitRoomRequest struct {
	RoomHash string `path:"room_hash"`
	UserId   string `json:"user_id"`
}

type GetDailyCoinRequest struct {
	UserId string `json:"user_id"`
}

type GetRoomRequest struct {
	RoomHash string `path:"room_hash"`
}

type JoinRoomRequest struct {
	RoomHash string `path:"room_hash"`
	UserId   string `json:"user_id"`
}

type NoRobLandownerRequest struct {
	RoomHash string `path:"room_hash"`
	UserId   string `json:"user_id"`
}

type NoSendCardRequest struct {
	RoomHash string `path:"room_hash"`
	UserId   string `json:"user_id"`
}

type RobLandownerRequest struct {
	RoomHash string `path:"room_hash"`
	UserId   string `json:"user_id"`
}

type SendCardRequest struct {
	RoomHash  string `path:"room_hash"`
	UserId    string `json:"user_id"`
	SendCards []Card `json:"sendcards"`
}

type StartRoomRequest struct {
	RoomHash string `path:"room_hash"`
	UserId   string `json:"user_id"`
}
