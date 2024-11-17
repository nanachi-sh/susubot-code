// 机器人核心响应处理
package handler

type ResponseT int

const (
	Message  ResponseT = iota //消息
	Event                     //事件
	CMDEvent                  //命令响应
)

type ResponseH struct {
	Data []byte
	Type ResponseT
}

func New(buf []byte, t ResponseT) *ResponseH {
	return &ResponseH{
		Data: buf,
		Type: t,
	}
}

func (h *ResponseH)  {
	
}