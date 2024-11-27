package request

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nanachi-sh/susubot-code/handler/LLOneBot/protos/handler"
	"github.com/nanachi-sh/susubot-code/handler/LLOneBot/protos/handler/request"
	"github.com/nanachi-sh/susubot-code/handler/LLOneBot/request/define"
)

const (
	sendGroupMessage   = "send_group_msg"
	sendFriendMessage  = "send_private_msg"
	recall             = "delete_msg"
	getMessage         = "get_msg"
	getGroupMemberInfo = "get_group_member_info"
	getGroupInfo       = "get_group_info"
	getFriendList      = "get_friend_list"
)

type requestH struct {
	req *request.BotRequestMarshalRequest
}

func New(req *request.BotRequestMarshalRequest) (*requestH, error) {
	if req == nil {
		return nil, errors.New("请求体不能为空")
	}
	return &requestH{req: req}, nil
}

func (rh *requestH) Marshal() ([]byte, error) {
	switch rh.req.Type {
	case request.RequestType_RequestType_GetFriendList:
		return rh.getFriendList()
	case request.RequestType_RequestType_GetGroupInfo:
		return rh.getGroupInfo()
	case request.RequestType_RequestType_GetGroupMemberInfo:
		return rh.getGroupMemberInfo()
	case request.RequestType_RequestType_GetMessage:
		return rh.getMessage()
	case request.RequestType_RequestType_MessageRecall:
		return rh.messageRecall()
	case request.RequestType_RequestType_SendFriendMessage:
		return rh.sendFriendMessage()
	case request.RequestType_RequestType_SendGroupMessage:
		return rh.sendGroupMessage()
	default:
		return nil, fmt.Errorf("请求类型无匹配; RequestType: %v", rh.req.Type.String())
	}
}

func (rh *requestH) getFriendList() ([]byte, error) {
	if rh.req.GetFriendList == nil {
		return nil, errors.New("GetFriendList结构体未定义")
	}
	req := new(define.Request)
	req.Action = getFriendList
	if rh.req.Echo != nil {
		req.Echo = *rh.req.Echo
	}
	return json.Marshal(req)
}

func (rh *requestH) getGroupInfo() ([]byte, error) {
	if rh.req.GetGroupInfo == nil {
		return nil, errors.New("GetGroupInfo结构体未定义")
	}
	ggi := rh.req.GetGroupInfo
	req := new(define.Request)
	req.Action = getGroupInfo
	req.Params = new(define.Request_Params)
	req.Params.GroupId = ggi.GroupId
	if rh.req.Echo != nil {
		req.Echo = *rh.req.Echo
	}
	return json.Marshal(req)
}

func (rh *requestH) getGroupMemberInfo() ([]byte, error) {
	if rh.req.GetGroupMemberInfo == nil {
		return nil, errors.New("GetGroupMemberInfo结构体未定义")
	}
	ggmi := rh.req.GetGroupMemberInfo
	req := new(define.Request)
	req.Action = getGroupMemberInfo
	req.Params = new(define.Request_Params)
	req.Params.GroupId = ggmi.GroupId
	req.Params.UserId = ggmi.UserId
	if rh.req.Echo != nil {
		req.Echo = *rh.req.Echo
	}
	return json.Marshal(req)
}

func (rh *requestH) messageRecall() ([]byte, error) {
	if rh.req.MessageRecall == nil {
		return nil, errors.New("MessageRecall结构体未定义")
	}
	mr := rh.req.MessageRecall
	req := new(define.Request)
	req.Action = recall
	req.Params = new(define.Request_Params)
	req.Params.MessageId = mr.MessageId
	if rh.req.Echo != nil {
		req.Echo = *rh.req.Echo
	}
	return json.Marshal(req)
}

func (rh *requestH) sendGroupMessage() ([]byte, error) {
	if rh.req.SendGroupMessage == nil {
		return nil, errors.New("SendGroupMessage结构体未定义")
	}
	sgm := rh.req.SendGroupMessage
	req := new(define.Request)
	req.Action = sendGroupMessage
	req.Params = new(define.Request_Params)
	req.Params.GroupId = sgm.GroupId
	mcs, err := marshalMessageChain(sgm.MessageChain)
	if err != nil {
		return nil, err
	}
	var mcs_j []map[string]any
	for _, v := range mcs {
		d, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var m map[string]any
		if err := json.Unmarshal(d, &m); err != nil {
			return nil, err
		}
		mcs_j = append(mcs_j, m)
	}
	req.Params.Message = mcs_j
	if rh.req.Echo != nil {
		req.Echo = *rh.req.Echo
	}
	return json.Marshal(req)
}

func (rh *requestH) sendFriendMessage() ([]byte, error) {
	if rh.req.SendFriendMessage == nil {
		return nil, errors.New("SendFriendMessage结构体未定义")
	}
	sfm := rh.req.SendFriendMessage
	req := new(define.Request)
	req.Action = sendFriendMessage
	req.Params = new(define.Request_Params)
	req.Params.UserId = sfm.FriendId
	mcs, err := marshalMessageChain(sfm.MessageChain)
	if err != nil {
		return nil, err
	}
	var mcs_j []map[string]any
	for _, v := range mcs {
		d, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var m map[string]any
		if err := json.Unmarshal(d, &m); err != nil {
			return nil, err
		}
		mcs_j = append(mcs_j, m)
	}
	req.Params.Message = mcs_j
	if rh.req.Echo != nil {
		req.Echo = *rh.req.Echo
	}
	return json.Marshal(req)
}

func (rh *requestH) getMessage() ([]byte, error) {
	if rh.req.GetMessage == nil {
		return nil, errors.New("GetMessage结构体为nil")
	}
	gm := rh.req.GetMessage
	req := new(define.Request)
	req.Action = getMessage
	req.Params = new(define.Request_Params)
	req.Params.MessageId = gm.MessageId
	if rh.req.Echo != nil {
		req.Echo = *rh.req.Echo
	}
	return json.Marshal(req)
}

func marshalMessageChain(mc []*handler.MessageChainObject) ([]*define.MessageChain, error) {
	ret := []*define.MessageChain{}
	for _, v := range mc {
		if v.Type == nil {
			return nil, errors.New("消息链存在Type为nil的消息")
		}
		switch *v.Type {
		case handler.MessageChainType_MessageChainType_Text:
			if v.Text == nil {
				return nil, errors.New("消息链Text结构体为nil")
			}
			ret = append(ret, &define.MessageChain{
				Data: map[string]any{
					"text": v.Text.Text,
				},
				Type: "text",
			})
		case handler.MessageChainType_MessageChainType_At:
			if v.At == nil {
				return nil, errors.New("消息链At结构体为nil")
			}
			ret = append(ret, &define.MessageChain{
				Data: map[string]any{
					"qq": v.At.TargetId,
				},
				Type: "at",
			})
		case handler.MessageChainType_MessageChainType_Reply:
			if v.Reply == nil {
				return nil, errors.New("消息链Reply结构体为nil")
			}
			ret = append(ret, &define.MessageChain{
				Data: map[string]any{
					"id": v.Reply.MessageId,
				},
				Type: "reply",
			})
		case handler.MessageChainType_MessageChainType_Image:
			if v.Image == nil {
				return nil, errors.New("消息链Image结构体为nil")
			}
			ret = append(ret, &define.MessageChain{
				Data: map[string]any{
					"url": v.Image.URL,
				},
				Type: "image",
			})
		case handler.MessageChainType_MessageChainType_Voice: //Voice
			if v.Voice == nil {
				return nil, errors.New("消息链Voice结构体为nil")
			}
			ret = append(ret, &define.MessageChain{
				Data: map[string]any{
					"url": v.Voice.URL,
				},
				Type: "record",
			})
		case handler.MessageChainType_MessageChainType_Video:
			if v.Video == nil {
				return nil, errors.New("消息链Video结构体为nil")
			}
			ret = append(ret, &define.MessageChain{
				Data: map[string]any{
					"url": v.Video.URL,
				},
				Type: "video",
			})
		default:
		}
	}
	return ret, nil
}
