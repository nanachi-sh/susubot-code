package request

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/define"
	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/log"
	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/protos/fileweb"
	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/protos/handler/request"
	request_d "github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/request/define"
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

var (
	logger = log.Get()
)

func GetFriendInfo(friendid string, echo *string) ([]byte, error) {
	req := new(request_d.Request)
	req.Action = getFriendInfo
	if echo != nil {
		req.Echo = *echo
	}
	req.Params = new(request_d.Request_Params)
	req.Params.UserId = friendid
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
}

func GetFriendList(echo *string) ([]byte, error) {
	req := new(request_d.Request)
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
	req := new(request_d.Request)
	req.Action = getGroupInfo
	req.Params = new(request_d.Request_Params)
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
	req := new(request_d.Request)
	req.Action = getGroupMemberInfo
	req.Params = new(request_d.Request_Params)
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
	req := new(request_d.Request)
	req.Action = recall
	req.Params = new(request_d.Request_Params)
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

func SendGroupMessage(groupid string, inMcs []*request.MessageChainObject, echo *string) ([]byte, error) {
	req := new(request_d.Request)
	req.Action = sendGroupMessage
	req.Params = new(request_d.Request_Params)
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

func SendFriendMessage(friendid string, inMcs []*request.MessageChainObject, echo *string) ([]byte, error) {
	req := new(request_d.Request)
	req.Action = sendFriendMessage
	req.Params = new(request_d.Request_Params)
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
	req := new(request_d.Request)
	req.Action = getMessage
	req.Params = new(request_d.Request_Params)
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

func marshalMessageChain(mc []*request.MessageChainObject) ([]*request_d.MessageChain, error) {
	ret := []*request_d.MessageChain{}
	for _, v := range mc {
		if v.Type == nil {
			return nil, errors.New("消息链存在Type为nil的消息")
		}
		switch *v.Type {
		case request.MessageChainType_MessageChainType_Text:
			if v.Text == nil {
				return nil, errors.New("消息链Text结构体为nil")
			}
			ret = append(ret, &request_d.MessageChain{
				Data: map[string]any{
					"text": v.Text.Text,
				},
				Type: "text",
			})
		case request.MessageChainType_MessageChainType_At:
			if v.At == nil {
				return nil, errors.New("消息链At结构体为nil")
			}
			ret = append(ret, &request_d.MessageChain{
				Data: map[string]any{
					"qq": v.At.TargetId,
				},
				Type: "at",
			})
		case request.MessageChainType_MessageChainType_Reply:
			if v.Reply == nil {
				return nil, errors.New("消息链Reply结构体为nil")
			}
			ret = append(ret, &request_d.MessageChain{
				Data: map[string]any{
					"id": v.Reply.MessageId,
				},
				Type: "reply",
			})
		case request.MessageChainType_MessageChainType_Image:
			image := v.Image
			if image == nil {
				return nil, errors.New("消息链Image结构体为nil")
			}
			u := ""
			if image.URL != nil {
				u = *image.URL
			} else if image.Buf != nil {
				filewebc := fileweb.NewFileWebClient(define.GRPCClient)
				resp, err := filewebc.Upload(define.FilewebCtx, &fileweb.UploadRequest{
					Buf: image.Buf,
				})
				if err != nil {
					return nil, err
				}
				u = fmt.Sprintf("play6.unturned.fun:1080%v", resp.URLPath)
			}
			ret = append(ret, &request_d.MessageChain{
				Data: map[string]any{
					"file": u,
				},
				Type: "image",
			})
		case request.MessageChainType_MessageChainType_Voice: //Voice
			voice := v.Voice
			if voice == nil {
				return nil, errors.New("消息链Voice结构体为nil")
			}
			u := ""
			if voice.URL != nil {
				u = *voice.URL
			} else if voice.Buf != nil {
				filewebc := fileweb.NewFileWebClient(define.GRPCClient)
				resp, err := filewebc.Upload(define.FilewebCtx, &fileweb.UploadRequest{
					Buf: voice.Buf,
				})
				if err != nil {
					return nil, err
				}
				u = fmt.Sprintf("play6.unturned.fun:1080%v", resp.URLPath)
			}
			ret = append(ret, &request_d.MessageChain{
				Data: map[string]any{
					"file": u,
				},
				Type: "record",
			})
		case request.MessageChainType_MessageChainType_Video:
			video := v.Video
			if video == nil {
				return nil, errors.New("消息链Video结构体为nil")
			}
			u := ""
			if video.URL != nil {
				u = *video.URL
			} else if video.Buf != nil {
				filewebc := fileweb.NewFileWebClient(define.GRPCClient)
				resp, err := filewebc.Upload(define.FilewebCtx, &fileweb.UploadRequest{
					Buf: video.Buf,
				})
				if err != nil {
					return nil, err
				}
				u = fmt.Sprintf("play6.unturned.fun:1080%v", resp.URLPath)
			}
			ret = append(ret, &request_d.MessageChain{
				Data: map[string]any{
					"file": u,
				},
				Type: "video",
			})
		default:
		}
	}
	return ret, nil
}
