package types

import (
	"fmt"
	"net/http"

	accountmanager_pb "github.com/nanachi-sh/susubot-code/basic/accountmanager/pkg/protos/accountmanager"
)

var defaultErrorsMap map[int32]string

func init() {
	defaultErrorsMap = map[int32]string{
		int32(accountmanager_pb.Error_ERROR_UNKNOWN):                       "未知错误",
		int32(accountmanager_pb.Error_ERROR_UNDEFINED):                     "未定义错误",
		int32(accountmanager_pb.Error_ERROR_INVALID_ARGUMENT):              "参数错误",
		int32(accountmanager_pb.Error_ERROR_NO_VERIFYCODE_AUTH):            "未进行验证码验证",
		int32(accountmanager_pb.Error_ERROR_EMAIL_EXISTED):                 "该邮箱已注册",
		int32(accountmanager_pb.Error_ERROR_VERIFYCODE_ANSWER_FAIL):        "验证码不正确",
		int32(accountmanager_pb.Error_ERROR_EMAIL_VERIFYCODE_ANSWER_FAIL):  "邮箱验证码不正确",
		int32(accountmanager_pb.Error_ERROR_EMAIL_VERIFYCODE_SEND_WAITING): "邮箱验证码发送冷却中",
	}
}

func NewError(code accountmanager_pb.Error, message string, statusCode ...int) *AppError {
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
	Code       accountmanager_pb.Error
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
		return http.StatusOK
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
