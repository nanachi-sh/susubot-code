// 机器人核心响应处理
package response

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/handler/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/handler/internal/handler/request"
	response_c "github.com/nanachi-sh/susubot-code/basic/handler/internal/handler/response/configs"
	response_t "github.com/nanachi-sh/susubot-code/basic/handler/internal/handler/response/types"
	connector_pb "github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/connector"
	request_pb "github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/request"
	response_pb "github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/response"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Request struct {
	logger logx.Logger
}

func NewRequest(l logx.Logger) *Request {
	return &Request{
		logger: l,
	}
}

func (r *Request) Unmarshal(in *response_pb.UnmarshalRequest) (*response_pb.UnmarshalResponse, error) {
	if len(in.Buf) == 0 {
		return &response_pb.UnmarshalResponse{}, status.Error(codes.InvalidArgument, "")
	}
	responseH, serr := newResponseHandle(r.logger, in.Buf, in.Type, in.CmdEventType, in.ExtraInfo)
	if serr != nil {
		return &response_pb.UnmarshalResponse{Body: &response_pb.UnmarshalResponse_Err{Err: *serr}}, nil
	}
	if in.IgnoreCmdEvent && responseH.ResponseType() == response_pb.ResponseType_ResponseType_CmdEvent {
		return &response_pb.UnmarshalResponse{Body: &response_pb.UnmarshalResponse_Err{Err: response_pb.Errors_TypeNoMatch}}, nil
	}
	response, serr := responseH.MarshalToResponse(r.logger)
	if serr != nil {
		return &response_pb.UnmarshalResponse{Body: &response_pb.UnmarshalResponse_Err{Err: *serr}}, nil
	}
	return response, nil
}

type responseH struct {
	buf   []byte
	rtype response_pb.ResponseType

	cet   *response_pb.CmdEventType
	extra bool
}

type botEventH struct {
	d   response_pb.Response_BotEvent
	buf []byte
}

type (
	cmdEventH struct {
		d   response_pb.Response_CmdEvent
		buf []byte
	}

	cmdEvent_GetMessageH struct {
		j   *response_t.JSON_cmdEvent_GetMessage
		d   response_pb.Response_CmdEvent_GetMessage
		buf []byte
	}
)

type messageH struct {
	d   response_pb.Response_Message
	buf []byte
}

type (
	qqeventH struct {
		d     response_pb.Response_QQEvent
		buf   []byte
		extra bool
	}

	qqevent_groupAddH struct {
		d     response_pb.Response_QQEvent_GroupAdd
		buf   []byte
		extra bool
	}

	qqevent_groupRemoveH struct {
		d     response_pb.Response_QQEvent_GroupRemove
		buf   []byte
		extra bool
	}

	qqevent_messageRecallH struct {
		d     response_pb.Response_QQEvent_MessageRecall
		buf   []byte
		extra bool
	}
)

func newResponseHandle(logger logx.Logger, buf []byte, targetType *response_pb.ResponseType, cet *response_pb.CmdEventType, extra bool) (*responseH, *response_pb.Errors) {
	rh := &responseH{
		buf:   buf,
		cet:   cet,
		extra: extra,
	}
	if targetType == nil {
		ret, serr := rh.matchType(logger)
		if serr != nil {
			return nil, serr
		}
		rh.rtype = ret
	} else {
		rh.rtype = *targetType
	}
	return rh, nil
}

func unmarshalMessageChain(logger logx.Logger, mc []*response_t.MessageChain) ([]*response_pb.MessageChainObject, *response_pb.Errors) {
	ret := []*response_pb.MessageChainObject{}
	for _, json := range mc {
		switch json.Type {
		case "text":
			d, ok := json.Data["text"]
			if !ok {
				logger.Error("text不存在")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response_pb.MessageChainObject{
					Type: response_pb.MessageChainType_MessageChainType_Text,
					Text: &response_pb.MessageChain_Text{
						Text: d,
					},
				})
			default:
				logger.Error("text不为string")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
		case "at":
			d, ok := json.Data["qq"]
			if !ok {
				logger.Error("qq(at)不存在")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response_pb.MessageChainObject{
					Type: response_pb.MessageChainType_MessageChainType_At,
					At: &response_pb.MessageChain_At{
						TargetId: d,
					},
				})
			default:
				logger.Error("qq(at)不为string")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
		case "reply":
			d, ok := json.Data["id"]
			if !ok {
				logger.Error("id(reply)不存在")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response_pb.MessageChainObject{
					Type: response_pb.MessageChainType_MessageChainType_Reply,
					Reply: &response_pb.MessageChain_Reply{
						MessageId: d,
					},
				})
			default:
				logger.Error("id(reply)不为string")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
		case "image":
			d, ok := json.Data["url"]
			if !ok {
				logger.Error("url(image)不存在")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response_pb.MessageChainObject{
					Type: response_pb.MessageChainType_MessageChainType_Image,
					Image: &response_pb.MessageChain_Image{
						URL: d,
					},
				})
			default:
				logger.Error("url(image)不为string")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
		case "record": //Voice
			d, ok := json.Data["url"]
			if !ok {
				logger.Error("url(record)不存在")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response_pb.MessageChainObject{
					Type: response_pb.MessageChainType_MessageChainType_Voice,
					Voice: &response_pb.MessageChain_Voice{
						URL: d,
					},
				})
			default:
				logger.Error("url(record)不为string")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
		case "video":
			d, ok := json.Data["url"]
			if !ok {
				logger.Error("url(video)不存在")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response_pb.MessageChainObject{
					Type: response_pb.MessageChainType_MessageChainType_Video,
					Video: &response_pb.MessageChain_Video{
						URL: d,
					},
				})
			default:
				logger.Error("url(video)不为string")
				return nil, response_pb.Errors_MessageChainError.Enum()
			}
		default:
		}
	}
	return ret, nil
}

func (rh *responseH) ResponseType() response_pb.ResponseType {
	return rh.rtype
}

func (rh *responseH) matchType(logger logx.Logger) (response_pb.ResponseType, *response_pb.Errors) {
	j := new(response_t.JSON_responseType)
	if err := json.Unmarshal(rh.buf, j); err != nil {
		logger.Error(err)
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.Echo != nil {
		return response_pb.ResponseType_ResponseType_CmdEvent, nil
	}
	if j.PostType == nil {
		logger.Error("响应事件类型无匹配，PostType为nil")
		return -1, response_pb.Errors_TypeNoMatch.Enum()
	} else {
		switch pt := *j.PostType; pt {
		case "message":
			return response_pb.ResponseType_ResponseType_Message, nil
		case "notice":
			return response_pb.ResponseType_ResponseType_QQEvent, nil
		case "meta_event":
			return response_pb.ResponseType_ResponseType_BotEvent, nil
		default:
			logger.Errorf("响应事件类型无匹配; PostType: %v", pt)
			return -1, response_pb.Errors_TypeNoMatch.Enum()
		}
	}
}

func (rh *responseH) MarshalToResponse(logger logx.Logger) (*response_pb.UnmarshalResponse, *response_pb.Errors) {
	ret := new(response_pb.UnmarshalResponse_ResponseDefine)
	ret.Type = rh.rtype
	switch rh.rtype {
	case response_pb.ResponseType_ResponseType_BotEvent:
		be, err := rh.BotEvent(logger)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		ret.Body = &response_pb.UnmarshalResponse_ResponseDefine_BotEvent{
			BotEvent: be,
		}
	case response_pb.ResponseType_ResponseType_CmdEvent:
		ce, err := rh.CmdEvent(logger)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		ret.Body = &response_pb.UnmarshalResponse_ResponseDefine_CmdEvent{
			CmdEvent: ce,
		}
	case response_pb.ResponseType_ResponseType_Message:
		m, err := rh.Message(logger)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		ret.Body = &response_pb.UnmarshalResponse_ResponseDefine_Message{
			Message: m,
		}
	case response_pb.ResponseType_ResponseType_QQEvent:
		qqe, err := rh.QQEvent(logger)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		ret.Body = &response_pb.UnmarshalResponse_ResponseDefine_QQEvent{
			QQEvent: qqe,
		}
	}
	return &response_pb.UnmarshalResponse{
		Body: &response_pb.UnmarshalResponse_Response{Response: ret},
	}, nil
}

func (rh *responseH) BotEvent(logger logx.Logger) (*response_pb.Response_BotEvent, *response_pb.Errors) {
	beh := new(botEventH)
	beh.buf = rh.buf
	t, serr := beh.matchType(logger)
	if serr != nil {
		return nil, serr
	}
	beh.d.Type = &t
	switch t {
	case response_pb.BotEventType_BotEventType_Connected:
		ret, serr := beh.Connected(logger)
		if serr != nil {
			return nil, serr
		}
		beh.d.Connected = ret
	case response_pb.BotEventType_BotEventType_HeartPacket:
		ret, serr := beh.HeartPacket(logger)
		if serr != nil {
			return nil, serr
		}
		beh.d.HeartPacket = ret
	default:
		logger.Errorf("机器人响应类型无匹配; MetaEventType: %v", beh.d.Type)
		return nil, response_pb.Errors_TypeNoMatch.Enum()
	}
	return &beh.d, nil
}

func (beh *botEventH) matchType(logger logx.Logger) (response_pb.BotEventType, *response_pb.Errors) {
	j := new(response_t.JSON_botEventType)
	if err := json.Unmarshal(beh.buf, j); err != nil {
		logger.Error(err)
		return -1, response_pb.Errors_Undefined.Enum()
	}
	met, st := j.MetaEventType, j.SubType
	switch {
	case met == "lifecycle" && st == "connect":
		return response_pb.BotEventType_BotEventType_Connected, nil
	case met == "heartbeat":
		return response_pb.BotEventType_BotEventType_HeartPacket, nil
	default:
		logger.Errorf("机器人响应类型无匹配; MetaEventType: %v, SubType: %v", met, st)
		return -1, response_pb.Errors_TypeNoMatch.Enum()
	}
}

func (beh *botEventH) Connected(logger logx.Logger) (*response_pb.Response_BotEvent_Connected, *response_pb.Errors) {
	j := new(response_t.JSON_botEvent_Connected)
	if err := json.Unmarshal(beh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	return &response_pb.Response_BotEvent_Connected{
		Timestamp: j.Timestamp,
		BotId:     strconv.FormatInt(j.SelfId, 10),
	}, nil
}

func (beh *botEventH) HeartPacket(logger logx.Logger) (*response_pb.Response_BotEvent_HeartPacket, *response_pb.Errors) {
	j := new(response_t.JSON_botEvent_HeartPacket)
	if err := json.Unmarshal(beh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	return &response_pb.Response_BotEvent_HeartPacket{
		Timestamp: j.Timestamp,
		BotId:     strconv.FormatInt(j.SelfId, 10),
		Interval:  j.Interval,
		Status: &response_pb.Response_BotEvent_HeartPacketStatus{
			Online: j.Status.Online,
			Good:   j.Status.Good,
		},
	}, nil
}

func (rh *responseH) CmdEvent(logger logx.Logger) (*response_pb.Response_CmdEvent, *response_pb.Errors) {
	ceh := new(cmdEventH)
	ceh.buf = rh.buf
	if rh.cet == nil {
		logger.Error("未指定命令响应类型")
		return nil, response_pb.Errors_CmdEventTypeNoSet.Enum()
	}
	ceh.d.Type = rh.cet
	e, serr := ceh.Echo(logger)
	if serr != nil {
		return nil, serr
	}
	ceh.d.Echo = e
	switch *rh.cet {
	case response_pb.CmdEventType_CmdEventType_GetFriendList:
		gfl, serr := ceh.GetFriendList(logger)
		if serr != nil {
			return nil, serr
		}
		ceh.d.GetFriendList = gfl
	case response_pb.CmdEventType_CmdEventType_GetGroupInfo:
		ggi, serr := ceh.GetGroupInfo(logger)
		if serr != nil {
			return nil, serr
		}
		ceh.d.GetGroupInfo = ggi
	case response_pb.CmdEventType_CmdEventType_GetGroupMemberInfo:
		ggmi, serr := ceh.GetGroupMemberInfo(logger)
		if serr != nil {
			return nil, serr
		}
		ceh.d.GetGroupMemberInfo = ggmi
	case response_pb.CmdEventType_CmdEventType_GetMessage:
		gm, serr := ceh.GetMessage(logger)
		if serr != nil {
			return nil, serr
		}
		ceh.d.GetMessage = gm
	case response_pb.CmdEventType_CmdEventType_GetFriendInfo:
		gfi, serr := ceh.GetFriendInfo(logger)
		if serr != nil {
			return nil, serr
		}
		ceh.d.GetFriendInfo = gfi
	}
	return &ceh.d, nil
}

func (ceh *cmdEventH) Echo(logger logx.Logger) (string, *response_pb.Errors) {
	j := new(response_t.JSON_cmdEvent_Echo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Error(err)
		return "", response_pb.Errors_Undefined.Enum()
	}
	return j.Echo, nil
}

func (ceh *cmdEventH) GetFriendInfo(logger logx.Logger) (*response_pb.Response_CmdEvent_GetFriendInfo, *response_pb.Errors) {
	j := new(response_t.JSON_cmdEvent_GetFriendInfo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	ok := false
	if j.Status == response_t.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response_pb.Response_CmdEvent_GetFriendInfo)
	if ok {
		ret = &response_pb.Response_CmdEvent_GetFriendInfo{
			OK:       true,
			UserName: j.NickName,
			UserId:   strconv.FormatInt(j.UserId, 10),
			Remark:   &j.Remark,
		}
	} else {
		ret.OK = false
		rc := strconv.FormatInt(int64(j.Retcode), 10)
		ret.Retcode = &rc
	}
	return ret, nil
}

func (ceh *cmdEventH) GetFriendList(logger logx.Logger) (*response_pb.Response_CmdEvent_GetFriendList, *response_pb.Errors) {
	j := new(response_t.JSON_cmdEvent_GetFriendList)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	ok := false
	if j.Status == response_t.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response_pb.Response_CmdEvent_GetFriendList)
	if ok {
		ret.OK = true
		ret.Retcode = nil
		friends := []*response_pb.Response_CmdEvent_GetFriendList_FriendInfo{}
		for _, v := range j.Data {
			remark := new(string)
			if v.REmark == "" {
				remark = nil
			} else {
				remark = &v.REmark
			}
			friends = append(friends, &response_pb.Response_CmdEvent_GetFriendList_FriendInfo{
				UserName: v.NickName,
				UserId:   strconv.FormatInt(v.User_Id, 10),
				Remark:   remark,
			})
		}
		ret.Friends = friends
	} else {
		ret.OK = false
		rc := strconv.FormatInt(int64(j.Retcode), 10)
		ret.Retcode = &rc
		ret.Friends = nil
	}
	return ret, nil
}

func (ceh *cmdEventH) GetGroupInfo(logger logx.Logger) (*response_pb.Response_CmdEvent_GetGroupInfo, *response_pb.Errors) {
	j := new(response_t.JSON_cmdEvent_GetGroupInfo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	ok := false
	if j.Status == response_t.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response_pb.Response_CmdEvent_GetGroupInfo)
	if ok {
		ret = &response_pb.Response_CmdEvent_GetGroupInfo{
			OK:        true,
			Retcode:   nil,
			GroupId:   strconv.FormatInt(j.Group.GroupId, 10),
			GroupName: j.Group.GroupName,
			MemberMax: int32(j.Group.MemberMax),
			MemberNow: int32(j.Group.MemberNow),
		}
	} else {
		ret.OK = false
		rc := strconv.FormatInt(int64(j.Retcode), 10)
		ret.Retcode = &rc
	}
	return ret, nil
}

func (ceh *cmdEventH) GetGroupMemberInfo(logger logx.Logger) (*response_pb.Response_CmdEvent_GetGroupMemberInfo, *response_pb.Errors) {
	j := new(response_t.JSON_cmdEvent_GetGroupMemberInfo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	ok := false
	if j.Status == response_t.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response_pb.Response_CmdEvent_GetGroupMemberInfo)
	if ok {
		card := new(string)
		if j.Member.Card == "" {
			card = nil
		} else {
			card = &j.Member.Card
		}
		jointime := new(int64)
		if j.Member.JoinTime == 0 {
			jointime = nil
		} else {
			jointime = &j.Member.JoinTime
		}
		lastactivetime := new(int64)
		if j.Member.LastActiveTime == 0 {
			lastactivetime = nil
		} else {
			lastactivetime = &j.Member.LastActiveTime
		}
		lastsenttime := new(int64)
		if j.Member.LastSentTime == 0 {
			lastsenttime = nil
		} else {
			lastsenttime = &j.Member.LastSentTime
		}
		var role response_pb.GroupRole
		switch j.Member.Role {
		case response_c.Role_Member:
			role = response_pb.GroupRole_GroupRole_Member
		case response_c.Role_Admin:
			role = response_pb.GroupRole_GroupRole_Admin
		case response_c.Role_Owner:
			role = response_pb.GroupRole_GroupRole_Owner
		}
		ret = &response_pb.Response_CmdEvent_GetGroupMemberInfo{
			OK:             true,
			Retcode:        nil,
			GroupId:        strconv.FormatInt(j.Member.GroupId, 10),
			UserId:         strconv.FormatInt(j.Member.UserId, 10),
			UserName:       j.Member.UserName,
			Card:           card,
			JoinTime:       jointime,
			LastActiveTime: lastactivetime,
			LastSentTime:   lastsenttime,
			Role:           role,
		}
	} else {
		ret.OK = false
		rc := strconv.FormatInt(int64(j.Retcode), 10)
		ret.Retcode = &rc
	}
	return ret, nil
}

func (ceh *cmdEventH) GetMessage(logger logx.Logger) (*response_pb.Response_CmdEvent_GetMessage, *response_pb.Errors) {
	j := new(response_t.JSON_cmdEvent_GetMessage)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	ok := false
	if j.Status == response_t.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response_pb.Response_CmdEvent_GetMessage)
	if ok {
		cegmh := &cmdEvent_GetMessageH{
			j:   nil,
			d:   response_pb.Response_CmdEvent_GetMessage{},
			buf: ceh.buf,
		}
		m, serr := cegmh.Message(logger)
		if serr != nil {
			return nil, serr
		}
		ret = &response_pb.Response_CmdEvent_GetMessage{
			OK:      ok,
			Retcode: nil,
			Message: m,
		}
	} else {
		ret.OK = false
		rc := strconv.FormatInt(int64(j.Retcode), 10)
		ret.Retcode = &rc
	}
	return ret, nil
}

func (cegmh *cmdEvent_GetMessageH) matchType(logger logx.Logger) (response_pb.MessageType, *response_pb.Errors) {
	j := cegmh.j
	switch mt, st := j.Data.MessageType, j.Data.SubType; {
	case mt == "private" && st == "friend":
		return response_pb.MessageType_MessageType_Private, nil
	case mt == "group" && st == "normal":
		return response_pb.MessageType_MessageType_Group, nil
	default:
		logger.Errorf("命令响应获取消息事件类型无匹配; MessageType: %v, SubType: %v", mt, st)
		return -1, response_pb.Errors_TypeNoMatch.Enum()
	}
}

func (cegmh *cmdEvent_GetMessageH) Message(logger logx.Logger) (*response_pb.Response_CmdEvent_Message, *response_pb.Errors) {
	if cegmh.j == nil {
		j := new(response_t.JSON_cmdEvent_GetMessage)
		if err := json.Unmarshal(cegmh.buf, j); err != nil {
			logger.Error(err)
			return nil, response_pb.Errors_Undefined.Enum()
		}
		cegmh.j = j
	}
	m := new(response_pb.Response_CmdEvent_Message)
	t, serr := cegmh.matchType(logger)
	if serr != nil {
		return nil, serr
	}
	m.Type = &t
	switch t {
	case response_pb.MessageType_MessageType_Group:
		g, serr := cegmh.group(logger)
		if serr != nil {
			return nil, serr
		}
		m.Group = g
	case response_pb.MessageType_MessageType_Private:
		p, serr := cegmh.private()
		if serr != nil {
			return nil, serr
		}
		m.Private = p
	}
	return m, nil
}

func (cegmh *cmdEvent_GetMessageH) group(logger logx.Logger) (*response_pb.Response_CmdEvent_Message_Group, *response_pb.Errors) {
	j := cegmh.j
	sname := &j.Data.Sender.Nickname
	var srole, brole *response_pb.GroupRole
	if j.Data.Sender.Role != "" {
		switch j.Data.Sender.Role {
		case response_c.Role_Member:
			srole = response_pb.GroupRole_GroupRole_Member.Enum()
		case response_c.Role_Admin:
			srole = response_pb.GroupRole_GroupRole_Admin.Enum()
		case response_c.Role_Owner:
			srole = response_pb.GroupRole_GroupRole_Owner.Enum()
		}
	}
	jmc := []*response_t.MessageChain{}
	for _, v := range j.Data.MessageChain {
		jmc = append(jmc, &response_t.MessageChain{
			Data: v.Data,
			Type: v.Type,
		})
	}
	mc, serr := unmarshalMessageChain(logger, jmc)
	if serr != nil {
		return nil, serr
	}
	return &response_pb.Response_CmdEvent_Message_Group{
		SenderId:     strconv.FormatInt(j.Data.Sender.UserId, 10),
		SenderName:   sname,
		MessageId:    strconv.FormatInt(j.Data.MessageId, 10),
		Timestamp:    j.Data.Timestamp,
		BotId:        strconv.FormatInt(j.Data.SelfId, 10),
		GroupId:      strconv.FormatInt(j.Data.GroupId, 10),
		SenderRole:   srole,
		BotRole:      brole,
		MessageChain: mc,
	}, nil
}

func (cegmh *cmdEvent_GetMessageH) private() (*response_pb.Response_CmdEvent_Message_Private, *response_pb.Errors) {
	return nil, response_pb.Errors_Undefined.Enum()
}

func (rh *responseH) Message(logger logx.Logger) (*response_pb.Response_Message, *response_pb.Errors) {
	mh := new(messageH)
	mh.buf = rh.buf
	t, serr := mh.matchType(logger)
	if serr != nil {
		return nil, serr
	}
	mh.d.Type = &t
	switch t {
	case response_pb.MessageType_MessageType_Group:
		g, serr := mh.group(logger)
		if serr != nil {
			return nil, serr
		}
		mh.d.Group = g
	case response_pb.MessageType_MessageType_Private:
		p, serr := mh.private(logger)
		if serr != nil {
			return nil, serr
		}
		mh.d.Private = p
	}
	return &mh.d, nil
}

func (mh *messageH) matchType(logger logx.Logger) (response_pb.MessageType, *response_pb.Errors) {
	j := new(response_t.JSON_messageType)
	if err := json.Unmarshal(mh.buf, j); err != nil {
		logger.Error(err)
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.MessageType == nil || j.SubType == nil {
		logger.Error("MessageType/SubType为nil")
		return -1, response_pb.Errors_Undefined.Enum()
	}
	switch mt, st := *j.MessageType, *j.SubType; {
	case mt == "private" && st == "friend":
		return response_pb.MessageType_MessageType_Private, nil
	case mt == "group" && st == "normal":
		return response_pb.MessageType_MessageType_Group, nil
	default:
		logger.Errorf("消息事件类型无匹配; MessageType: %v, SubType: %v", mt, st)
		return -1, response_pb.Errors_TypeNoMatch.Enum()
	}
}

func (mh *messageH) group(logger logx.Logger) (*response_pb.Response_Message_Group, *response_pb.Errors) {
	j := new(response_t.JSON_message_Group)
	if err := json.Unmarshal(mh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	sname := &j.Sender.Nickname
	var srole, brole *response_pb.GroupRole
	if j.Sender.Role != "" {
		switch j.Sender.Role {
		case response_c.Role_Member:
			srole = response_pb.GroupRole_GroupRole_Member.Enum()
		case response_c.Role_Admin:
			srole = response_pb.GroupRole_GroupRole_Admin.Enum()
		case response_c.Role_Owner:
			srole = response_pb.GroupRole_GroupRole_Owner.Enum()
		}
	}
	jmc := []*response_t.MessageChain{}
	for _, v := range j.MessageChain {
		jmc = append(jmc, &response_t.MessageChain{
			Data: v.Data,
			Type: v.Type,
		})
	}
	mc, serr := unmarshalMessageChain(logger, jmc)
	if serr != nil {
		return nil, serr
	}
	return &response_pb.Response_Message_Group{
		SenderId:     strconv.FormatInt(j.Sender.UserId, 10),
		SenderName:   sname,
		MessageId:    strconv.FormatInt(j.MessageId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		SenderRole:   srole,
		BotRole:      brole,
		MessageChain: mc,
	}, nil
}

func (mh *messageH) private(logger logx.Logger) (*response_pb.Response_Message_Private, *response_pb.Errors) {
	j := new(response_t.JSON_message_Private)
	if err := json.Unmarshal(mh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	sname := &j.Sender.Nickname
	jmc := []*response_t.MessageChain{}
	for _, v := range j.MessageChain {
		jmc = append(jmc, &response_t.MessageChain{
			Data: v.Data,
			Type: v.Type,
		})
	}
	mc, serr := unmarshalMessageChain(logger, jmc)
	if serr != nil {
		return nil, serr
	}
	return &response_pb.Response_Message_Private{
		SenderId:     strconv.FormatInt(j.Sender.UserId, 10),
		SenderName:   sname,
		MessageId:    strconv.FormatInt(j.MessageId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		MessageChain: mc,
	}, nil
}

func (rh *responseH) QQEvent(logger logx.Logger) (*response_pb.Response_QQEvent, *response_pb.Errors) {
	qeh := new(qqeventH)
	qeh.buf = rh.buf
	t, serr := qeh.matchType(logger)
	if serr != nil {
		return nil, serr
	}
	qeh.extra = rh.extra
	qeh.d.Type = &t
	switch t {
	case response_pb.QQEventType_QQEventType_GroupAdd:
		qegah := new(qqevent_groupAddH)
		qegah.buf = rh.buf
		qegah.extra = rh.extra
		ga, serr := qegah.GroupAdd(logger)
		if serr != nil {
			return nil, serr
		}
		qeh.d.GroupAdd = ga
	case response_pb.QQEventType_QQEventType_GroupRemove:
		qegrh := new(qqevent_groupRemoveH)
		qegrh.buf = rh.buf
		qegrh.extra = rh.extra
		gr, serr := qegrh.GroupRemove(logger)
		if serr != nil {
			return nil, serr
		}
		qeh.d.GroupRemove = gr
	case response_pb.QQEventType_QQEventType_GroupMute:
		gm, serr := qeh.groupMute(logger)
		if serr != nil {
			return nil, serr
		}
		qeh.d.GroupMute = gm
	case response_pb.QQEventType_QQEventType_GroupUnmute:
		gum, serr := qeh.groupUnmute(logger)
		if serr != nil {
			return nil, serr
		}
		qeh.d.GroupUnmute = gum
	case response_pb.QQEventType_QQEventType_MessageRecall:
		qemrh := new(qqevent_messageRecallH)
		qemrh.buf = rh.buf
		qemrh.extra = rh.extra
		mr, serr := qemrh.MessageRecall(logger)
		if serr != nil {
			return nil, serr
		}
		qeh.d.MessageRecall = mr
	}
	return &qeh.d, nil
}

func (qeh *qqeventH) matchType(logger logx.Logger) (response_pb.QQEventType, *response_pb.Errors) {
	j := new(response_t.JSON_qqEventType)
	if err := json.Unmarshal(qeh.buf, j); err != nil {
		logger.Error(err)
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.NoticeType == nil {
		logger.Error("NoticeType为nil")
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.SubType == nil {
		j.SubType = new(string)
	}
	switch mt, st := *j.NoticeType, *j.SubType; {
	default:
		logger.Errorf("QQ事件类型无匹配; NoticeType: %v, SubType: %v", mt, st)
		return -1, response_pb.Errors_TypeNoMatch.Enum()
	case mt == "group_increase":
		return response_pb.QQEventType_QQEventType_GroupAdd, nil
	case mt == "group_decrease":
		return response_pb.QQEventType_QQEventType_GroupRemove, nil
	case mt == "group_ban" && st == "ban":
		return response_pb.QQEventType_QQEventType_GroupMute, nil
	case mt == "group_ban" && st == "lift_ban":
		return response_pb.QQEventType_QQEventType_GroupUnmute, nil
	case mt == "group_recall", mt == "friend_recall":
		return response_pb.QQEventType_QQEventType_MessageRecall, nil
	}
}

func (qeh *qqeventH) groupMute(logger logx.Logger) (*response_pb.Response_QQEvent_GroupMute, *response_pb.Errors) {
	j := new(response_t.JSON_qqEvent_groupMute)
	if err := json.Unmarshal(qeh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	var (
		targetName   *string
		operatorName *string
	)
	if qeh.extra {
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		useridStr := strconv.FormatInt(j.UserId, 10)
		operatoridStr := strconv.FormatInt(j.OperatorId, 10)
		ctx, cancel := context.WithTimeout(configs.DefaultCtx, time.Second*15)
		defer cancel()
		stream, err := configs.Call_Connector.Read(ctx, &connector_pb.Empty{})
		if err != nil {
			logger.Error(err)
			return nil, response_pb.Errors_Undefined.Enum()
		}
		user_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, useridStr, stream)
		if serr != nil {
			return nil, serr
		}
		targetName = &user_ggmi.UserName
		operator_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, operatoridStr, stream)
		if serr != nil {
			return nil, serr
		}
		operatorName = &operator_ggmi.UserName
	}
	return &response_pb.Response_QQEvent_GroupMute{
		TargetId:     strconv.FormatInt(j.UserId, 10),
		TargetName:   targetName,
		Timestamp:    j.Timestamp,
		OperatorId:   strconv.FormatInt(j.OperatorId, 10),
		OperatorName: operatorName,
		Duration:     int32(j.Duration),
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		BotId:        strconv.FormatInt(j.SelfId, 10),
	}, nil
}

func (qeh *qqeventH) groupUnmute(logger logx.Logger) (*response_pb.Response_QQEvent_GroupUnmute, *response_pb.Errors) {
	j := new(response_t.JSON_qqEvent_groupUnmute)
	if err := json.Unmarshal(qeh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	var (
		targetName   *string
		operatorName *string
	)
	if qeh.extra {
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		useridStr := strconv.FormatInt(j.UserId, 10)
		operatoridStr := strconv.FormatInt(j.OperatorId, 10)
		ctx, cancel := context.WithTimeout(configs.DefaultCtx, time.Second*15)
		defer cancel()
		stream, err := configs.Call_Connector.Read(ctx, &connector_pb.Empty{})
		if err != nil {
			logger.Error(err)
			return nil, response_pb.Errors_Undefined.Enum()
		}
		user_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, useridStr, stream)
		if serr != nil {
			return nil, serr
		}
		targetName = &user_ggmi.UserName
		operator_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, operatoridStr, stream)
		if serr != nil {
			return nil, serr
		}
		operatorName = &operator_ggmi.UserName
	}
	return &response_pb.Response_QQEvent_GroupUnmute{
		TargetId:     strconv.FormatInt(j.UserId, 10),
		TargetName:   targetName,
		Timestamp:    j.Timestamp,
		OperatorId:   strconv.FormatInt(j.OperatorId, 10),
		OperatorName: operatorName,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		BotId:        strconv.FormatInt(j.SelfId, 10),
	}, nil
}

// WIP
func (qegah *qqevent_groupAddH) matchType(logger logx.Logger) (response_pb.QQEventType_GroupAddType, *response_pb.Errors) {
	j := new(response_t.JSON_qqEventType)
	if err := json.Unmarshal(qegah.buf, j); err != nil {
		logger.Error(err)
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.NoticeType == nil {
		logger.Error("NoticeType为nil")
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.SubType == nil {
		j.SubType = new(string)
	}
	switch mt, st := *j.NoticeType, *j.SubType; {
	default:
		logger.Errorf("QQ事件群增加类型无匹配; NoticeType: %v, SubType: %v", mt, st)
		return -1, response_pb.Errors_TypeNoMatch.Enum()
	case mt == "group_increase":
		return response_pb.QQEventType_GroupAddType_QQEventType_GroupAddType_Direct, nil
	}
}

func (qegah *qqevent_groupAddH) GroupAdd(logger logx.Logger) (*response_pb.Response_QQEvent_GroupAdd, *response_pb.Errors) {
	t, serr := qegah.matchType(logger)
	if serr != nil {
		return nil, serr
	}
	qegah.d.Type = &t
	switch t {
	case response_pb.QQEventType_GroupAddType_QQEventType_GroupAddType_Direct:
		d, serr := qegah.direct(logger)
		if serr != nil {
			return nil, serr
		}
		qegah.d.Direct = d
	case response_pb.QQEventType_GroupAddType_QQEventType_GroupAddType_Invite:
		i, serr := qegah.invite()
		if serr != nil {
			return nil, serr
		}
		qegah.d.Invite = i
	}
	return &qegah.d, nil
}

func (qegah *qqevent_groupAddH) direct(logger logx.Logger) (*response_pb.Response_QQEvent_GroupAdd_Direct, *response_pb.Errors) {
	j := new(response_t.JSON_qqEvent_groupAdd_direct)
	if err := json.Unmarshal(qegah.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	var (
		joinerName   *string
		approverName *string
	)
	if qegah.extra {
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		useridStr := strconv.FormatInt(j.UserId, 10)
		operatoridStr := strconv.FormatInt(j.OperatorId, 10)
		ctx, cancel := context.WithTimeout(configs.DefaultCtx, time.Second*15)
		defer cancel()
		stream, err := configs.Call_Connector.Read(ctx, &connector_pb.Empty{})
		if err != nil {
			logger.Error(err)
			return nil, response_pb.Errors_Undefined.Enum()
		}
		user_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, useridStr, stream)
		if serr != nil {
			return nil, serr
		}
		operator_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, operatoridStr, stream)
		if serr != nil {
			return nil, serr
		}
		joinerName = &user_ggmi.UserName
		approverName = &operator_ggmi.UserName
	}
	return &response_pb.Response_QQEvent_GroupAdd_Direct{
		JoinerId:     strconv.FormatInt(j.UserId, 10),
		JoinerName:   joinerName,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		ApproverId:   strconv.FormatInt(j.OperatorId, 10),
		ApproverName: approverName,
	}, nil
}

func (qegah *qqevent_groupAddH) invite() (*response_pb.Response_QQEvent_GroupAdd_Invite, *response_pb.Errors) {
	return nil, response_pb.Errors_Undefined.Enum()
}

func (qegrh *qqevent_groupRemoveH) GroupRemove(logger logx.Logger) (*response_pb.Response_QQEvent_GroupRemove, *response_pb.Errors) {
	t, serr := qegrh.matchType(logger)
	if serr != nil {
		return nil, serr
	}
	qegrh.d.Type = &t
	switch t {
	case response_pb.QQEventType_GroupRemoveType_QQEventType_GroupRemoveType_Kick:
		k, serr := qegrh.kick(logger)
		if serr != nil {
			return nil, serr
		}
		qegrh.d.Kick = k
	case response_pb.QQEventType_GroupRemoveType_QQEventType_GroupRemoveType_Manual:
		m, serr := qegrh.manual(logger)
		if serr != nil {
			return nil, serr
		}
		qegrh.d.Manual = m
	}
	return &qegrh.d, nil
}

func (qegrh *qqevent_groupRemoveH) matchType(logger logx.Logger) (response_pb.QQEventType_GroupRemoveType, *response_pb.Errors) {
	j := new(response_t.JSON_qqEventType)
	if err := json.Unmarshal(qegrh.buf, j); err != nil {
		logger.Error(err)
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.NoticeType == nil {
		logger.Error("NoticeType为nil")
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.SubType == nil {
		j.SubType = new(string)
	}
	switch mt, st := *j.NoticeType, *j.SubType; {
	default:
		logger.Errorf("QQ事件群减少类型无匹配; NoticeType: %v, SubType: %v", mt, st)
		return -1, response_pb.Errors_Undefined.Enum()
	case mt == "group_decrease" && st == "leave":
		return response_pb.QQEventType_GroupRemoveType_QQEventType_GroupRemoveType_Manual, nil
	case mt == "group_decrease" && st == "kick":
		return response_pb.QQEventType_GroupRemoveType_QQEventType_GroupRemoveType_Kick, nil
	}
}

func (qegrh *qqevent_groupRemoveH) manual(logger logx.Logger) (*response_pb.Response_QQEvent_GroupRemove_Manual, *response_pb.Errors) {
	j := new(response_t.JSON_qqEvent_groupRemove_manual)
	if err := json.Unmarshal(qegrh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	var quiterName *string
	if qegrh.extra {
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		useridStr := strconv.FormatInt(j.UserId, 10)
		ctx, cancel := context.WithTimeout(configs.DefaultCtx, time.Second*15)
		defer cancel()
		ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, useridStr, nil)
		if serr != nil {
			return nil, serr
		}
		quiterName = &ggmi.UserName
	}
	return &response_pb.Response_QQEvent_GroupRemove_Manual{
		QuiterId:   strconv.FormatInt(j.UserId, 10),
		QuiterName: quiterName,
		GroupId:    strconv.FormatInt(j.GroupId, 10),
		Timestamp:  j.Timestamp,
		BotId:      strconv.FormatInt(j.SelfId, 10),
	}, nil
}

func (qegrh *qqevent_groupRemoveH) kick(logger logx.Logger) (*response_pb.Response_QQEvent_GroupRemove_Kick, *response_pb.Errors) {
	j := new(response_t.JSON_qqEvent_groupRemove_kick)
	if err := json.Unmarshal(qegrh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	var (
		quiterName   *string
		operatorName *string
	)
	if qegrh.extra {
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		quiteridStr := strconv.FormatInt(j.UserId, 10)
		operatoridStr := strconv.FormatInt(j.OperatorId, 10)
		ctx, cancel := context.WithTimeout(configs.DefaultCtx, time.Second*15)
		defer cancel()
		stream, err := configs.Call_Connector.Read(ctx, &connector_pb.Empty{})
		if err != nil {
			logger.Error(err)
			return nil, response_pb.Errors_Undefined.Enum()
		}
		quiter_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, quiteridStr, stream)
		if serr != nil {
			return nil, serr
		}
		operator_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, operatoridStr, stream)
		if serr != nil {
			return nil, serr
		}
		quiterName = &quiter_ggmi.UserName
		operatorName = &operator_ggmi.UserName
	}
	return &response_pb.Response_QQEvent_GroupRemove_Kick{
		TargetId:     strconv.FormatInt(j.UserId, 10),
		TargetName:   quiterName,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		OperatorId:   strconv.FormatInt(j.OperatorId, 10),
		OperatorName: operatorName,
	}, nil
}

func (qemrh *qqevent_messageRecallH) MessageRecall(logger logx.Logger) (*response_pb.Response_QQEvent_MessageRecall, *response_pb.Errors) {
	t, serr := qemrh.matchType(logger)
	if serr != nil {
		return nil, serr
	}
	qemrh.d.Type = &t
	switch t {
	case response_pb.QQEventType_MessageRecallType_QQEventType_MessageRecallType_Group:
		g, serr := qemrh.group(logger)
		if serr != nil {
			return nil, serr
		}
		qemrh.d.Group = g
	case response_pb.QQEventType_MessageRecallType_QQEventType_MessageRecallType_Private:
		p, serr := qemrh.private(logger)
		if serr != nil {
			return nil, serr
		}
		qemrh.d.Private = p
	}
	return &qemrh.d, nil
}

func (qemrh *qqevent_messageRecallH) matchType(logger logx.Logger) (response_pb.QQEventType_MessageRecallType, *response_pb.Errors) {
	j := new(response_t.JSON_qqEventType)
	if err := json.Unmarshal(qemrh.buf, j); err != nil {
		logger.Error(err)
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.NoticeType == nil {
		logger.Error("NoticeType为nil")
		return -1, response_pb.Errors_Undefined.Enum()
	}
	if j.SubType == nil {
		j.SubType = new(string)
	}
	switch mt, st := *j.NoticeType, *j.SubType; {
	default:
		logger.Errorf("QQ事件类型消息撤回无匹配; NoticeType: %v, SubType: %v", mt, st)
		return -1, response_pb.Errors_TypeNoMatch.Enum()
	case mt == "group_recall":
		return response_pb.QQEventType_MessageRecallType_QQEventType_MessageRecallType_Group, nil
	case mt == "friend_recall":
		return response_pb.QQEventType_MessageRecallType_QQEventType_MessageRecallType_Private, nil
	}
}

func (qemrh *qqevent_messageRecallH) group(logger logx.Logger) (*response_pb.Response_QQEvent_MessageRecall_Group, *response_pb.Errors) {
	j := new(response_t.JSON_qqEvent_messageRecall_group)
	if err := json.Unmarshal(qemrh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	var (
		targetName   *string
		operatorName *string
	)
	if qemrh.extra {
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		ctx, cancel := context.WithTimeout(configs.DefaultCtx, time.Second*15)
		defer cancel()
		// 同一人
		if j.UserId == j.OperatorId {
			useridStr := strconv.FormatInt(j.UserId, 10)
			ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, useridStr, nil)
			if serr != nil {
				return nil, serr
			}
			targetName = &ggmi.UserName
		} else { //不同人
			useridStr := strconv.FormatInt(j.UserId, 10)
			operatoridStr := strconv.FormatInt(j.OperatorId, 10)
			stream, err := configs.Call_Connector.Read(ctx, &connector_pb.Empty{})
			if err != nil {
				logger.Error(err)
				return nil, response_pb.Errors_Undefined.Enum()
			}
			user_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, useridStr, stream)
			if serr != nil {
				return nil, serr
			}
			targetName = &user_ggmi.UserName
			operator_ggmi, serr := getGroupMemberInfo(logger, ctx, groupidStr, operatoridStr, stream)
			if serr != nil {
				return nil, serr
			}
			operatorName = &operator_ggmi.UserName
		}
	}
	return &response_pb.Response_QQEvent_MessageRecall_Group{
		TargetId:     strconv.FormatInt(j.UserId, 10),
		TargetName:   targetName,
		Timestamp:    j.Timestamp,
		OperatorId:   strconv.FormatInt(j.OperatorId, 10),
		OperatorName: operatorName,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		BotId:        strconv.FormatInt(j.SelfId, 10),
		MessageId:    strconv.FormatInt(j.MessageId, 10),
	}, nil
}

func (qemrh *qqevent_messageRecallH) private(logger logx.Logger) (*response_pb.Response_QQEvent_MessageRecall_Private, *response_pb.Errors) {
	j := new(response_t.JSON_qqEvent_messageRecall_private)
	if err := json.Unmarshal(qemrh.buf, j); err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	var recallerName *string
	ctx, cancel := context.WithTimeout(configs.DefaultCtx, time.Second*15)
	defer cancel()
	if qemrh.extra {
		useridStr := strconv.FormatInt(j.UserId, 10)
		gfi, serr := getFriendInfo(logger, ctx, useridStr, nil)
		if serr != nil {
			return nil, serr
		}
		recallerName = &gfi.UserName
	}
	return &response_pb.Response_QQEvent_MessageRecall_Private{
		RecallerId:   strconv.FormatInt(j.UserId, 10),
		RecallerName: recallerName,
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		MessageId:    strconv.FormatInt(j.MessageId, 10),
	}, nil
}

func sendCommand(logger logx.Logger, ctx context.Context, buf []byte, requestEcho string, stream grpc.ServerStreamingClient[connector_pb.ReadResponse]) (*cmdEventH, *response_pb.Errors) {
	if stream == nil {
		x, err := configs.Call_Connector.Read(ctx, &connector_pb.Empty{})
		if err != nil {
			logger.Error(err)
			return nil, response_pb.Errors_Undefined.Enum()
		}
		stream = x
	}
	readCh := make(chan any, 1)
	//写入
	go func() {
		if _, err := configs.Call_Connector.Write(ctx, &connector_pb.WriteRequest{
			Buf: buf,
		}); err != nil {
			logger.Error(err)
			return
		}
	}()
	//读取
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				readCh <- struct{}{}
				logger.Error(err)
				return
			}
			if serr := resp.GetErr(); serr != connector_pb.Errors_EMPTY {
				readCh <- serr
				return
			}
			buf := resp.GetBuf()
			respH := new(responseH)
			respH.buf = buf
			rt, serr := respH.matchType(logger)
			if serr != nil {
				readCh <- serr
				return
			}
			if rt != response_pb.ResponseType_ResponseType_CmdEvent {
				continue
			}
			ceh := new(cmdEventH)
			ceh.buf = buf
			echo, serr := ceh.Echo(logger)
			if serr != nil {
				continue
			}
			if echo != requestEcho {
				continue
			}
			readCh <- ceh
			return
		}
	}()
	select {
	case x := <-readCh:
		switch x := x.(type) {
		case struct{}:
			return nil, response_pb.Errors_Undefined.Enum()
		case *response_pb.Errors:
			return nil, x
		case *cmdEventH:
			return x, nil
		default:
			logger.Error("异常错误")
			return nil, response_pb.Errors_Undefined.Enum()
		}
	case <-ctx.Done():
		logger.Error("获取额外用户信息超时")
		return nil, response_pb.Errors_Undefined.Enum()
	}
}

func getGroupMemberInfo(logger logx.Logger, ctx context.Context, groupid, memberid string, stream grpc.ServerStreamingClient[connector_pb.ReadResponse]) (*response_pb.Response_CmdEvent_GetGroupMemberInfo, *response_pb.Errors) {
	requestEcho := strconv.FormatInt(rand.Int63(), 10)
	//获取写入内容
	request := request.NewRequest(logger)
	resp, err := request.GetGroupMemberInfo(&request_pb.GetGroupMemberInfoRequest{
		Echo:    &requestEcho,
		GroupId: groupid,
		UserId:  memberid,
	})
	if err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	ceh, serr := sendCommand(logger, ctx, resp.Buf, requestEcho, stream)
	if serr != nil {
		return nil, serr
	}
	return ceh.GetGroupMemberInfo(logger)
}

func getFriendInfo(logger logx.Logger, ctx context.Context, friendid string, stream grpc.ServerStreamingClient[connector_pb.ReadResponse]) (*response_pb.Response_CmdEvent_GetFriendInfo, *response_pb.Errors) {
	requestEcho := strconv.FormatInt(rand.Int63(), 10)
	//获取写入内容
	request := request.NewRequest(logger)
	resp, err := request.GetFriendInfo(&request_pb.GetFriendInfoRequest{
		Echo:     &requestEcho,
		FriendId: friendid,
	})
	if err != nil {
		logger.Error(err)
		return nil, response_pb.Errors_Undefined.Enum()
	}
	ceh, serr := sendCommand(logger, ctx, resp.Buf, requestEcho, stream)
	if serr != nil {
		return nil, serr
	}
	return ceh.GetFriendInfo(logger)
}
