// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package types

type Card struct {
	Number int `form:"number,range=[0:14]"`
}

type ExitRoomRequest struct {
	RoomHash  string  `path:"room_hash"`
	UserId    string  `form:"user_id"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	WinCount  int     `custom:"wincount"`
	LoseCount int     `custom:"losecount"`
	Coin      float64 `custom:"coin"`
}

type GetDailyCoinRequest struct {
	UserId    string  `form:"user_id"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	WinCount  int     `custom:"wincount"`
	LoseCount int     `custom:"losecount"`
	Coin      float64 `custom:"coin"`
}

type GetRoomRequest struct {
	RoomHash  string  `path:"room_hash"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	WinCount  int     `custom:"wincount"`
	LoseCount int     `custom:"losecount"`
	Coin      float64 `custom:"coin"`
}

type JoinRoomRequest struct {
	RoomHash  string  `path:"room_hash"`
	UserId    string  `form:"user_id"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	WinCount  int     `custom:"wincount"`
	LoseCount int     `custom:"losecount"`
	Coin      float64 `custom:"coin"`
}

type NoRobLandownerRequest struct {
	RoomHash  string  `path:"room_hash"`
	UserId    string  `form:"user_id"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	WinCount  int     `custom:"wincount"`
	LoseCount int     `custom:"losecount"`
	Coin      float64 `custom:"coin"`
}

type NoSendCardRequest struct {
	RoomHash  string  `path:"room_hash"`
	UserId    string  `form:"user_id"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	WinCount  int     `custom:"wincount"`
	LoseCount int     `custom:"losecount"`
	Coin      float64 `custom:"coin"`
}

type RobLandownerRequest struct {
	RoomHash  string  `path:"room_hash"`
	UserId    string  `form:"user_id"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	WinCount  int     `custom:"wincount"`
	LoseCount int     `custom:"losecount"`
	Coin      float64 `custom:"coin"`
}

type SendCardRequest struct {
	RoomHash  string  `path:"room_hash"`
	UserId    string  `form:"user_id"`
	SendCards []Card  `form:"sendcards"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	WinCount  int     `custom:"wincount"`
	LoseCount int     `custom:"losecount"`
	Coin      float64 `custom:"coin"`
}

type StartRoomRequest struct {
	RoomHash  string  `path:"room_hash"`
	UserId    string  `form:"user_id"`
	Name      string  `custom:"name"`
	Email     string  `custom:"email"`
	WinCount  int     `custom:"wincount"`
	LoseCount int     `custom:"losecount"`
	Coin      float64 `custom:"coin"`
}
