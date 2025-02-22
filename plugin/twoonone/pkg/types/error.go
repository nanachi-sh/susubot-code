package types

import (
	"fmt"
	"net/http"

	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
)

var defaultErrorsMap map[int32]string

func init() {
	defaultErrorsMap = map[int32]string{
		int32(twoonone_pb.Error_ERROR_UNKNOWN):          "未知错误",
		int32(twoonone_pb.Error_ERROR_UNDEFINED):        "未定义错误，请查看日志",
		int32(twoonone_pb.Error_ERROR_INVALID_ARGUMENT): "参数有误",

		int32(twoonone_pb.Error_ERROR_NO_AUTH):         "用户未认证",
		int32(twoonone_pb.Error_ERROR_USER_NO_EXIST):   "用户不存在",
		int32(twoonone_pb.Error_ERROR_USER_INCOMPLETE): "用户信息不完整",
		int32(twoonone_pb.Error_ERROR_USER_EXISTED):    "用户已存在",

		int32(twoonone_pb.Error_ERROR_ROOM_NO_EXIST):                                "房间不存在",
		int32(twoonone_pb.Error_ERROR_ROOM_EXIST_PLAYER):                            "玩家已在房间内",
		int32(twoonone_pb.Error_ERROR_ROOM_NO_EXIST_PLAYER):                         "玩家不在房间内",
		int32(twoonone_pb.Error_ERROR_ROOM_NO_ROB_LANDOWNERING):                     "房间不在抢地主阶段",
		int32(twoonone_pb.Error_ERROR_ROOM_NO_SENDING_CARD):                         "房间不在出牌阶段",
		int32(twoonone_pb.Error_ERROR_PLAYER_CARD_NO_EXIST):                         "玩家手牌不足",
		int32(twoonone_pb.Error_ERROR_PLAYER_NO_OPERATOR):                           "玩家不是当前操作者",
		int32(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST):                              "玩家不存在，这个错误不应该发生",
		int32(twoonone_pb.Error_ERROR_PLAYER_ALREADY_GET_DALIY_COIN):                "玩家已领取过每日豆子",
		int32(twoonone_pb.Error_ERROR_ROOM_NO_FULL):                                 "房间还未满",
		int32(twoonone_pb.Error_ERROR_ROOM_STARTED):                                 "房间已开始游戏",
		int32(twoonone_pb.Error_ERROR_ROOM_NO_ROB_LANDOWNER):                        "无人抢地主",
		int32(twoonone_pb.Error_ERROR_SEND_CARD_TYPE_UNKNOWN):                       "未知牌型",
		int32(twoonone_pb.Error_ERROR_SEND_CARD_TYPE_NE_LAST_CARD_TYPE):             "你的牌型与上一位玩家出的不同",
		int32(twoonone_pb.Error_ERROR_SEND_CARD_CONTINUOUS_NE_LAST_CARD_CONTINUOUS): "你的牌连续数与上一位玩家出的不同",
		int32(twoonone_pb.Error_ERROR_SEND_CARD_SIZE_LE_LAST_CARD_SIZE):             "你的牌大小比不过或等于上一位玩家",
		int32(twoonone_pb.Error_ERROR_PLAYER_EXISTED_A_ROOM):                        "玩家已在一个房间内",
		int32(twoonone_pb.Error_ERROR_PLAYER_COIN_LT_ROOM_COIN):                     "玩家豆子数少于房间底分",
		int32(twoonone_pb.Error_ERROR_PLAYER_IS_ONLY_OPERATOR):                      "玩家为唯一可操作者",
		int32(twoonone_pb.Error_ERROR_PLAYER_NO_EXIST_ANY_ROOM):                     "玩家不在任意房间内",
	}
}

func NewError(code twoonone_pb.Error, message string, statusCode ...int) *AppError {
	sc := 0
	if len(statusCode) > 0 {
		sc = statusCode[0]
	}
	return &AppError{
		Code:       code,
		statusCode: sc,
		message:    message,
	}
}

type AppError struct {
	Code       twoonone_pb.Error
	statusCode int
	message    string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Error, Code: %d, Message: %s", e.Code, e.Message())
}

func (e *AppError) Message() string {
	if e.message == "" {
		return e.defaultMessage()
	} else {
		return e.message
	}
}

func (e *AppError) StatusCode() int {
	if e.statusCode == 0 {
		return http.StatusBadRequest
	}
	return e.statusCode
}

func (e *AppError) defaultMessage() string {
	if e, ok := defaultErrorsMap[int32(e.Code)]; ok {
		return e
	} else {
		return "未定义"
	}
}
