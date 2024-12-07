package request

import (
	"encoding/json"
	"errors"

	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/log"
	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/protos/handler"
	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/request/define"
)

const (
	sendGroupMessage   = "send_group_msg"
	sendFriendMessage  = "send_private_msg"
	recall             = "delete_msg"
	getMessage         = "get_msg"
	getGroupMemberInfo = "get_group_member_info"
	getGroupInfo       = "get_group_info"
	getFriendList      = "get_friend_list"
	getFriendInfo      = "get_stranger_info"
)

var logger = log.Get()

func GetFriendInfo(friendid string, echo *string) ([]byte, error) {
	req := new(define.Request)
	req.Action = getFriendInfo
	if echo != nil {
		req.Echo = *echo
	}
	req.Params = new(define.Request_Params)
	req.Params.UserId = friendid
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
}

func GetFriendList(echo *string) ([]byte, error) {
	req := new(define.Request)
	req.Action = getFriendList
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
}

func GetGroupInfo(groupid string, echo *string) ([]byte, error) {
	req := new(define.Request)
	req.Action = getGroupInfo
	req.Params = new(define.Request_Params)
	req.Params.GroupId = groupid
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
}

func GetGroupMemberInfo(groupid, memberid string, echo *string) ([]byte, error) {
	req := new(define.Request)
	req.Action = getGroupMemberInfo
	req.Params = new(define.Request_Params)
	req.Params.GroupId = groupid
	req.Params.UserId = memberid
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
}

func MessageRecall(messageid string, echo *string) ([]byte, error) {
	req := new(define.Request)
	req.Action = recall
	req.Params = new(define.Request_Params)
	req.Params.MessageId = messageid
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
}

func SendGroupMessage(groupid string, inMcs []*handler.MessageChainObject, echo *string) ([]byte, error) {
	req := new(define.Request)
	req.Action = sendGroupMessage
	req.Params = new(define.Request_Params)
	req.Params.GroupId = groupid
	mcs, err := marshalMessageChain(inMcs)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	var mcs_j []map[string]any
	for _, v := range mcs {
		d, err := json.Marshal(v)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		var m map[string]any
		if err := json.Unmarshal(d, &m); err != nil {
			logger.Println(err)
			return nil, err
		}
		mcs_j = append(mcs_j, m)
	}
	req.Params.Message = mcs_j
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
}

func SendFriendMessage(friendid string, inMcs []*handler.MessageChainObject, echo *string) ([]byte, error) {
	req := new(define.Request)
	req.Action = sendFriendMessage
	req.Params = new(define.Request_Params)
	req.Params.UserId = friendid
	mcs, err := marshalMessageChain(inMcs)
	if err != nil {
		return nil, err
	}
	var mcs_j []map[string]any
	for _, v := range mcs {
		d, err := json.Marshal(v)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		var m map[string]any
		if err := json.Unmarshal(d, &m); err != nil {
			logger.Println(err)
			return nil, err
		}
		mcs_j = append(mcs_j, m)
	}
	req.Params.Message = mcs_j
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
}

func GetMessage(messageid string, echo *string) ([]byte, error) {
	req := new(define.Request)
	req.Action = getMessage
	req.Params = new(define.Request_Params)
	req.Params.MessageId = messageid
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
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
					"file": v.Image.URL,
				},
				Type: "image",
			})
		case handler.MessageChainType_MessageChainType_Voice: //Voice
			if v.Voice == nil {
				return nil, errors.New("消息链Voice结构体为nil")
			}
			ret = append(ret, &define.MessageChain{
				Data: map[string]any{
					"file": v.Voice.URL,
				},
				Type: "record",
			})
		case handler.MessageChainType_MessageChainType_Video:
			if v.Video == nil {
				return nil, errors.New("消息链Video结构体为nil")
			}
			ret = append(ret, &define.MessageChain{
				Data: map[string]any{
					"file": v.Video.URL,
				},
				Type: "video",
			})
		default:
		}
	}
	return ret, nil
}
