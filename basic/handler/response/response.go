// 机器人核心响应处理
package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/handler/define"
	"github.com/nanachi-sh/susubot-code/basic/handler/log"
	"github.com/nanachi-sh/susubot-code/basic/handler/protos/connector"
	"github.com/nanachi-sh/susubot-code/basic/handler/protos/handler/response"
	"github.com/nanachi-sh/susubot-code/basic/handler/request"
	response_d "github.com/nanachi-sh/susubot-code/basic/handler/response/define"
	"google.golang.org/grpc"
)

var (
	logger = log.Get()
)

type responseH struct {
	buf   []byte
	rtype response.ResponseType

	cet   *response.CmdEventType
	extra bool
}

type botEventH struct {
	d   response.Response_BotEvent
	buf []byte
}

type (
	cmdEventH struct {
		d   response.Response_CmdEvent
		buf []byte
	}

	cmdEvent_GetMessageH struct {
		j   *response_d.JSON_cmdEvent_GetMessage
		d   response.Response_CmdEvent_GetMessage
		buf []byte
	}
)

type messageH struct {
	d   response.Response_Message
	buf []byte
}

type (
	qqeventH struct {
		d     response.Response_QQEvent
		buf   []byte
		extra bool
	}

	qqevent_groupAddH struct {
		d     response.Response_QQEvent_GroupAdd
		buf   []byte
		extra bool
	}

	qqevent_groupRemoveH struct {
		d     response.Response_QQEvent_GroupRemove
		buf   []byte
		extra bool
	}

	qqevent_messageRecallH struct {
		d     response.Response_QQEvent_MessageRecall
		buf   []byte
		extra bool
	}
)

func New(req *response.UnmarshalRequest) (*responseH, error) {
	rh := &responseH{
		buf:   req.Buf,
		cet:   req.CmdEventType,
		extra: req.ExtraInfo,
	}
	t := req.Type
	if t == nil {
		ret, err := rh.matchType()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		rh.rtype = ret
	} else {
		rh.rtype = *t
	}
	return rh, nil
}

func unmarshalMessageChain(mc []*response_d.MessageChain) ([]*response.MessageChainObject, error) {
	ret := []*response.MessageChainObject{}
	for _, json := range mc {
		switch json.Type {
		case "text":
			d, ok := json.Data["text"]
			if !ok {
				return nil, errors.New("text不存在")
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response.MessageChainObject{
					Type: response.MessageChainType_MessageChainType_Text,
					Text: &response.MessageChain_Text{
						Text: d,
					},
				})
			default:
				return nil, errors.New("text不为string")
			}
		case "at":
			d, ok := json.Data["qq"]
			if !ok {
				return nil, errors.New("qq(at)不存在")
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response.MessageChainObject{
					Type: response.MessageChainType_MessageChainType_At,
					At: &response.MessageChain_At{
						TargetId: d,
					},
				})
			default:
				return nil, errors.New("qq(at)不为string")
			}
		case "reply":
			d, ok := json.Data["id"]
			if !ok {
				return nil, errors.New("id(reply)不存在")
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response.MessageChainObject{
					Type: response.MessageChainType_MessageChainType_Reply,
					Reply: &response.MessageChain_Reply{
						MessageId: d,
					},
				})
			default:
				return nil, errors.New("id(reply)不为string")
			}
		case "image":
			d, ok := json.Data["url"]
			if !ok {
				return nil, errors.New("url(image)不存在")
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response.MessageChainObject{
					Type: response.MessageChainType_MessageChainType_Image,
					Image: &response.MessageChain_Image{
						URL: d,
					},
				})
			default:
				return nil, errors.New("url(image)不为string")
			}
		case "record": //Voice
			d, ok := json.Data["url"]
			if !ok {
				return nil, errors.New("url(record)不存在")
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response.MessageChainObject{
					Type: response.MessageChainType_MessageChainType_Voice,
					Voice: &response.MessageChain_Voice{
						URL: d,
					},
				})
			default:
				return nil, errors.New("url(record)不为string")
			}
		case "video":
			d, ok := json.Data["url"]
			if !ok {
				return nil, errors.New("url(video)不存在")
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &response.MessageChainObject{
					Type: response.MessageChainType_MessageChainType_Video,
					Video: &response.MessageChain_Video{
						URL: d,
					},
				})
			default:
				return nil, errors.New("url(video)不为string")
			}
		default:
		}
	}
	return ret, nil
}

func (rh *responseH) ResponseType() response.ResponseType {
	return rh.rtype
}

func (rh *responseH) matchType() (response.ResponseType, error) {
	j := new(response_d.JSON_responseType)
	if err := json.Unmarshal(rh.buf, j); err != nil {
		return -1, err
	}
	if j.Echo != nil {
		return response.ResponseType_ResponseType_CmdEvent, nil
	}
	if j.PostType == nil {
		return -1, errors.New("响应事件类型无匹配，PostType为nil")
	} else {
		switch pt := *j.PostType; pt {
		case "message":
			return response.ResponseType_ResponseType_Message, nil
		case "notice":
			return response.ResponseType_ResponseType_QQEvent, nil
		case "meta_event":
			return response.ResponseType_ResponseType_BotEvent, nil
		default:
			return -1, fmt.Errorf("响应事件类型无匹配; PostType: %v", pt)
		}
	}
}

func (rh *responseH) MarshalToResponse() (*response.UnmarshalResponse, error) {
	ret := new(response.UnmarshalResponse)
	ret.Type = &rh.rtype
	switch rh.rtype {
	case response.ResponseType_ResponseType_BotEvent:
		be, err := rh.BotEvent()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		ret.BotEvent = be
	case response.ResponseType_ResponseType_CmdEvent:
		ce, err := rh.CmdEvent()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		ret.CmdEvent = ce
	case response.ResponseType_ResponseType_Message:
		m, err := rh.Message()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		ret.Message = m
	case response.ResponseType_ResponseType_QQEvent:
		qqe, err := rh.QQEvent()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		ret.QQEvent = qqe
	}
	return ret, nil
}

func (rh *responseH) BotEvent() (*response.Response_BotEvent, error) {
	beh := new(botEventH)
	beh.buf = rh.buf
	t, err := beh.matchType()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	beh.d.Type = &t
	switch t {
	case response.BotEventType_BotEventType_Connected:
		ret, err := beh.Connected()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		beh.d.Connected = ret
	case response.BotEventType_BotEventType_HeartPacket:
		ret, err := beh.HeartPacket()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		beh.d.HeartPacket = ret
	default:
		return nil, fmt.Errorf("机器人响应类型无匹配; MetaEventType: %v", beh.d.Type)
	}
	return &beh.d, nil
}

func (beh *botEventH) matchType() (response.BotEventType, error) {
	j := new(response_d.JSON_botEventType)
	if err := json.Unmarshal(beh.buf, j); err != nil {
		return -1, err
	}
	met, st := j.MetaEventType, j.SubType
	switch {
	case met == "lifecycle" && st == "connect":
		return response.BotEventType_BotEventType_Connected, nil
	case met == "heartbeat":
		return response.BotEventType_BotEventType_HeartPacket, nil
	default:
		return -1, fmt.Errorf("机器人响应类型无匹配; MetaEventType: %v, SubType: %v", met, st)
	}
}

func (beh *botEventH) Connected() (*response.Response_BotEvent_Connected, error) {
	j := new(response_d.JSON_botEvent_Connected)
	if err := json.Unmarshal(beh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	return &response.Response_BotEvent_Connected{
		Timestamp: j.Timestamp,
		BotId:     strconv.FormatInt(j.SelfId, 10),
	}, nil
}

func (beh *botEventH) HeartPacket() (*response.Response_BotEvent_HeartPacket, error) {
	j := new(response_d.JSON_botEvent_HeartPacket)
	if err := json.Unmarshal(beh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	return &response.Response_BotEvent_HeartPacket{
		Timestamp: j.Timestamp,
		BotId:     strconv.FormatInt(j.SelfId, 10),
		Interval:  j.Interval,
		Status: &response.Response_BotEvent_HeartPacketStatus{
			Online: j.Status.Online,
			Good:   j.Status.Good,
		},
	}, nil
}

func (rh *responseH) CmdEvent() (*response.Response_CmdEvent, error) {
	ceh := new(cmdEventH)
	ceh.buf = rh.buf
	if rh.cet == nil {
		return nil, errors.New("未指定命令响应类型")
	}
	ceh.d.Type = rh.cet
	e, err := ceh.Echo()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	ceh.d.Echo = e
	switch *rh.cet {
	case response.CmdEventType_CmdEventType_GetFriendList:
		gfl, err := ceh.GetFriendList()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		ceh.d.GetFriendList = gfl
	case response.CmdEventType_CmdEventType_GetGroupInfo:
		ggi, err := ceh.GetGroupInfo()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		ceh.d.GetGroupInfo = ggi
	case response.CmdEventType_CmdEventType_GetGroupMemberInfo:
		ggmi, err := ceh.GetGroupMemberInfo()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		ceh.d.GetGroupMemberInfo = ggmi
	case response.CmdEventType_CmdEventType_GetMessage:
		gm, err := ceh.GetMessage()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		ceh.d.GetMessage = gm
	case response.CmdEventType_CmdEventType_GetFriendInfo:

	}
	return &ceh.d, nil
}

func (ceh *cmdEventH) Echo() (string, error) {
	j := new(response_d.JSON_cmdEvent_Echo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Println(err)
		return "", err
	}
	return j.Echo, nil
}

func (ceh *cmdEventH) GetFriendInfo() (*response.Response_CmdEvent_GetFriendInfo, error) {
	j := new(response_d.JSON_cmdEvent_GetFriendInfo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	ok := false
	if j.Status == response_d.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response.Response_CmdEvent_GetFriendInfo)
	if ok {
		ret = &response.Response_CmdEvent_GetFriendInfo{
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

func (ceh *cmdEventH) GetFriendList() (*response.Response_CmdEvent_GetFriendList, error) {
	j := new(response_d.JSON_cmdEvent_GetFriendList)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	ok := false
	if j.Status == response_d.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response.Response_CmdEvent_GetFriendList)
	if ok {
		ret.OK = true
		ret.Retcode = nil
		friends := []*response.Response_CmdEvent_GetFriendList_FriendInfo{}
		for _, v := range j.Data {
			remark := new(string)
			if v.REmark == "" {
				remark = nil
			} else {
				remark = &v.REmark
			}
			friends = append(friends, &response.Response_CmdEvent_GetFriendList_FriendInfo{
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

func (ceh *cmdEventH) GetGroupInfo() (*response.Response_CmdEvent_GetGroupInfo, error) {
	j := new(response_d.JSON_cmdEvent_GetGroupInfo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	ok := false
	if j.Status == response_d.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response.Response_CmdEvent_GetGroupInfo)
	if ok {
		ret = &response.Response_CmdEvent_GetGroupInfo{
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

func (ceh *cmdEventH) GetGroupMemberInfo() (*response.Response_CmdEvent_GetGroupMemberInfo, error) {
	j := new(response_d.JSON_cmdEvent_GetGroupMemberInfo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	ok := false
	if j.Status == response_d.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response.Response_CmdEvent_GetGroupMemberInfo)
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
		var role response.GroupRole
		switch j.Member.Role {
		case response_d.Role_Member:
			role = response.GroupRole_GroupRole_Member
		case response_d.Role_Admin:
			role = response.GroupRole_GroupRole_Admin
		case response_d.Role_Owner:
			role = response.GroupRole_GroupRole_Owner
		}
		ret = &response.Response_CmdEvent_GetGroupMemberInfo{
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

func (ceh *cmdEventH) GetMessage() (*response.Response_CmdEvent_GetMessage, error) {
	j := new(response_d.JSON_cmdEvent_GetMessage)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	ok := false
	if j.Status == response_d.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(response.Response_CmdEvent_GetMessage)
	if ok {
		cegmh := &cmdEvent_GetMessageH{
			j:   nil,
			d:   response.Response_CmdEvent_GetMessage{},
			buf: ceh.buf,
		}
		m, err := cegmh.Message()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		ret = &response.Response_CmdEvent_GetMessage{
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

func (cegmh *cmdEvent_GetMessageH) matchType() (response.MessageType, error) {
	j := cegmh.j
	switch mt, st := j.Data.MessageType, j.Data.SubType; {
	case mt == "private" && st == "friend":
		return response.MessageType_MessageType_Private, nil
	case mt == "group" && st == "normal":
		return response.MessageType_MessageType_Group, nil
	default:
		return -1, fmt.Errorf("命令响应获取消息事件类型无匹配; MessageType: %v, SubType: %v", mt, st)
	}
}

func (cegmh *cmdEvent_GetMessageH) Message() (*response.Response_CmdEvent_Message, error) {
	if cegmh.j == nil {
		j := new(response_d.JSON_cmdEvent_GetMessage)
		if err := json.Unmarshal(cegmh.buf, j); err != nil {
			logger.Println(err)
			return nil, err
		}
		cegmh.j = j
	}
	m := new(response.Response_CmdEvent_Message)
	t, err := cegmh.matchType()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	m.Type = &t
	switch t {
	case response.MessageType_MessageType_Group:
		g, err := cegmh.group()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		m.Group = g
	case response.MessageType_MessageType_Private:
		p, err := cegmh.private()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		m.Private = p
	}
	return m, nil
}

func (cegmh *cmdEvent_GetMessageH) group() (*response.Response_CmdEvent_Message_Group, error) {
	j := cegmh.j
	sname := &j.Data.Sender.Nickname
	var srole, brole *response.GroupRole
	if j.Data.Sender.Role != "" {
		switch j.Data.Sender.Role {
		case response_d.Role_Member:
			srole = response.GroupRole_GroupRole_Member.Enum()
		case response_d.Role_Admin:
			srole = response.GroupRole_GroupRole_Admin.Enum()
		case response_d.Role_Owner:
			srole = response.GroupRole_GroupRole_Owner.Enum()
		}
	}
	jmc := []*response_d.MessageChain{}
	for _, v := range j.Data.MessageChain {
		jmc = append(jmc, &response_d.MessageChain{
			Data: v.Data,
			Type: v.Type,
		})
	}
	mc, err := unmarshalMessageChain(jmc)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return &response.Response_CmdEvent_Message_Group{
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

func (cegmh *cmdEvent_GetMessageH) private() (*response.Response_CmdEvent_Message_Private, error) {
	return nil, errors.ErrUnsupported
}

func (rh *responseH) Message() (*response.Response_Message, error) {
	mh := new(messageH)
	mh.buf = rh.buf
	t, err := mh.matchType()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	mh.d.Type = &t
	switch t {
	case response.MessageType_MessageType_Group:
		g, err := mh.group()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		mh.d.Group = g
	case response.MessageType_MessageType_Private:
		p, err := mh.private()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		mh.d.Private = p
	}
	return &mh.d, nil
}

func (mh *messageH) matchType() (response.MessageType, error) {
	j := new(response_d.JSON_messageType)
	if err := json.Unmarshal(mh.buf, j); err != nil {
		logger.Println(err)
		return -1, err
	}
	if j.MessageType == nil || j.SubType == nil {
		return -1, errors.New("MessageType/SubType为nil")
	}
	switch mt, st := *j.MessageType, *j.SubType; {
	case mt == "private" && st == "friend":
		return response.MessageType_MessageType_Private, nil
	case mt == "group" && st == "normal":
		return response.MessageType_MessageType_Group, nil
	default:
		return -1, fmt.Errorf("消息事件类型无匹配; MessageType: %v, SubType: %v", mt, st)
	}
}

func (mh *messageH) group() (*response.Response_Message_Group, error) {
	j := new(response_d.JSON_message_Group)
	if err := json.Unmarshal(mh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	sname := &j.Sender.Nickname
	var srole, brole *response.GroupRole
	if j.Sender.Role != "" {
		switch j.Sender.Role {
		case response_d.Role_Member:
			srole = response.GroupRole_GroupRole_Member.Enum()
		case response_d.Role_Admin:
			srole = response.GroupRole_GroupRole_Admin.Enum()
		case response_d.Role_Owner:
			srole = response.GroupRole_GroupRole_Owner.Enum()
		}
	}
	jmc := []*response_d.MessageChain{}
	for _, v := range j.MessageChain {
		jmc = append(jmc, &response_d.MessageChain{
			Data: v.Data,
			Type: v.Type,
		})
	}
	mc, err := unmarshalMessageChain(jmc)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return &response.Response_Message_Group{
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

func (mh *messageH) private() (*response.Response_Message_Private, error) {
	j := new(response_d.JSON_message_Private)
	if err := json.Unmarshal(mh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	sname := &j.Sender.Nickname
	jmc := []*response_d.MessageChain{}
	for _, v := range j.MessageChain {
		jmc = append(jmc, &response_d.MessageChain{
			Data: v.Data,
			Type: v.Type,
		})
	}
	mc, err := unmarshalMessageChain(jmc)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return &response.Response_Message_Private{
		SenderId:     strconv.FormatInt(j.Sender.UserId, 10),
		SenderName:   sname,
		MessageId:    strconv.FormatInt(j.MessageId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		MessageChain: mc,
	}, nil
}

func (rh *responseH) QQEvent() (*response.Response_QQEvent, error) {
	qeh := new(qqeventH)
	qeh.buf = rh.buf
	t, err := qeh.matchType()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	qeh.extra = rh.extra
	qeh.d.Type = &t
	switch t {
	case response.QQEventType_QQEventType_GroupAdd:
		qegah := new(qqevent_groupAddH)
		qegah.buf = rh.buf
		qegah.extra = rh.extra
		ga, err := qegah.GroupAdd()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qeh.d.GroupAdd = ga
	case response.QQEventType_QQEventType_GroupRemove:
		qegrh := new(qqevent_groupRemoveH)
		qegrh.buf = rh.buf
		qegrh.extra = rh.extra
		gr, err := qegrh.GroupRemove()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qeh.d.GroupRemove = gr
	case response.QQEventType_QQEventType_GroupMute:
		gm, err := qeh.groupMute()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qeh.d.GroupMute = gm
	case response.QQEventType_QQEventType_GroupUnmute:
		gum, err := qeh.groupUnmute()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qeh.d.GroupUnmute = gum
	case response.QQEventType_QQEventType_MessageRecall:
		qemrh := new(qqevent_messageRecallH)
		qemrh.buf = rh.buf
		qemrh.extra = rh.extra
		mr, err := qemrh.MessageRecall()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qeh.d.MessageRecall = mr
	}
	return &qeh.d, nil
}

func (qeh *qqeventH) matchType() (response.QQEventType, error) {
	j := new(response_d.JSON_qqEventType)
	if err := json.Unmarshal(qeh.buf, j); err != nil {
		return -1, err
	}
	if j.NoticeType == nil {
		return -1, errors.New("NoticeType为nil")
	}
	if j.SubType == nil {
		j.SubType = new(string)
	}
	switch mt, st := *j.NoticeType, *j.SubType; {
	default:
		return -1, fmt.Errorf("QQ事件类型无匹配; NoticeType: %v, SubType: %v", mt, st)
	case mt == "group_increase":
		return response.QQEventType_QQEventType_GroupAdd, nil
	case mt == "group_decrease":
		return response.QQEventType_QQEventType_GroupRemove, nil
	case mt == "group_ban" && st == "ban":
		return response.QQEventType_QQEventType_GroupMute, nil
	case mt == "group_ban" && st == "lift_ban":
		return response.QQEventType_QQEventType_GroupUnmute, nil
	case mt == "group_recall", mt == "friend_recall":
		return response.QQEventType_QQEventType_MessageRecall, nil
	}
}

func (qeh *qqeventH) groupMute() (*response.Response_QQEvent_GroupMute, error) {
	j := new(response_d.JSON_qqEvent_groupMute)
	if err := json.Unmarshal(qeh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	var (
		targetName   *string
		operatorName *string
	)
	if qeh.extra {
		connectorClient := connector.NewConnectorClient(define.GRPCClient)
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		useridStr := strconv.FormatInt(j.UserId, 10)
		operatoridStr := strconv.FormatInt(j.OperatorId, 10)
		ctx, cancel := context.WithTimeout(define.ConnectorCtx, time.Second*15)
		defer cancel()
		stream, err := connectorClient.Read(ctx, &connector.Empty{})
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		user_ggmi, err := getGroupMemberInfo(ctx, groupidStr, useridStr, stream)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		targetName = &user_ggmi.UserName
		operator_ggmi, err := getGroupMemberInfo(ctx, groupidStr, operatoridStr, stream)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		operatorName = &operator_ggmi.UserName
	}
	return &response.Response_QQEvent_GroupMute{
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

func (qeh *qqeventH) groupUnmute() (*response.Response_QQEvent_GroupUnmute, error) {
	j := new(response_d.JSON_qqEvent_groupUnmute)
	if err := json.Unmarshal(qeh.buf, j); err != nil {
		return nil, err
	}
	var (
		targetName   *string
		operatorName *string
	)
	if qeh.extra {
		connectorClient := connector.NewConnectorClient(define.GRPCClient)
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		useridStr := strconv.FormatInt(j.UserId, 10)
		operatoridStr := strconv.FormatInt(j.OperatorId, 10)
		ctx, cancel := context.WithTimeout(define.ConnectorCtx, time.Second*15)
		defer cancel()
		stream, err := connectorClient.Read(ctx, &connector.Empty{})
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		user_ggmi, err := getGroupMemberInfo(ctx, groupidStr, useridStr, stream)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		targetName = &user_ggmi.UserName
		operator_ggmi, err := getGroupMemberInfo(ctx, groupidStr, operatoridStr, stream)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		operatorName = &operator_ggmi.UserName
	}
	return &response.Response_QQEvent_GroupUnmute{
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
func (qegah *qqevent_groupAddH) matchType() (response.QQEventType_GroupAddType, error) {
	j := new(response_d.JSON_qqEventType)
	if err := json.Unmarshal(qegah.buf, j); err != nil {
		return -1, err
	}
	if j.NoticeType == nil {
		return -1, errors.New("NoticeType为nil")
	}
	if j.SubType == nil {
		j.SubType = new(string)
	}
	switch mt, st := *j.NoticeType, *j.SubType; {
	default:
		return -1, fmt.Errorf("QQ事件群增加类型无匹配; NoticeType: %v, SubType: %v", mt, st)
	case mt == "group_increase":
		return response.QQEventType_GroupAddType_QQEventType_GroupAddType_Direct, nil
	}
}

func (qegah *qqevent_groupAddH) GroupAdd() (*response.Response_QQEvent_GroupAdd, error) {
	t, err := qegah.matchType()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	qegah.d.Type = &t
	switch t {
	case response.QQEventType_GroupAddType_QQEventType_GroupAddType_Direct:
		d, err := qegah.direct()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qegah.d.Direct = d
	case response.QQEventType_GroupAddType_QQEventType_GroupAddType_Invite:
		i, err := qegah.invite()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qegah.d.Invite = i
	}
	return &qegah.d, nil
}

func (qegah *qqevent_groupAddH) direct() (*response.Response_QQEvent_GroupAdd_Direct, error) {
	j := new(response_d.JSON_qqEvent_groupAdd_direct)
	if err := json.Unmarshal(qegah.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	var (
		joinerName   *string
		approverName *string
	)
	if qegah.extra {
		connectorClient := connector.NewConnectorClient(define.GRPCClient)
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		useridStr := strconv.FormatInt(j.UserId, 10)
		operatoridStr := strconv.FormatInt(j.OperatorId, 10)
		ctx, cancel := context.WithTimeout(define.ConnectorCtx, time.Second*15)
		defer cancel()
		stream, err := connectorClient.Read(ctx, &connector.Empty{})
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		user_ggmi, err := getGroupMemberInfo(ctx, groupidStr, useridStr, stream)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		operator_ggmi, err := getGroupMemberInfo(ctx, groupidStr, operatoridStr, stream)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		joinerName = &user_ggmi.UserName
		approverName = &operator_ggmi.UserName
	}
	return &response.Response_QQEvent_GroupAdd_Direct{
		JoinerId:     strconv.FormatInt(j.UserId, 10),
		JoinerName:   joinerName,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		ApproverId:   strconv.FormatInt(j.OperatorId, 10),
		ApproverName: approverName,
	}, nil
}

func (qegah *qqevent_groupAddH) invite() (*response.Response_QQEvent_GroupAdd_Invite, error) {
	return nil, errors.ErrUnsupported
}

func (qegrh *qqevent_groupRemoveH) GroupRemove() (*response.Response_QQEvent_GroupRemove, error) {
	t, err := qegrh.matchType()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	qegrh.d.Type = &t
	switch t {
	case response.QQEventType_GroupRemoveType_QQEventType_GroupRemoveType_Kick:
		k, err := qegrh.kick()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qegrh.d.Kick = k
	case response.QQEventType_GroupRemoveType_QQEventType_GroupRemoveType_Manual:
		m, err := qegrh.manual()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qegrh.d.Manual = m
	}
	return &qegrh.d, nil
}

func (qegrh *qqevent_groupRemoveH) matchType() (response.QQEventType_GroupRemoveType, error) {
	j := new(response_d.JSON_qqEventType)
	if err := json.Unmarshal(qegrh.buf, j); err != nil {
		logger.Println(err)
		return -1, err
	}
	if j.NoticeType == nil {
		return -1, errors.New("NoticeType为nil")
	}
	if j.SubType == nil {
		j.SubType = new(string)
	}
	switch mt, st := *j.NoticeType, *j.SubType; {
	default:
		return -1, fmt.Errorf("QQ事件群减少类型无匹配; NoticeType: %v, SubType: %v", mt, st)
	case mt == "group_decrease" && st == "leave":
		return response.QQEventType_GroupRemoveType_QQEventType_GroupRemoveType_Manual, nil
	case mt == "group_decrease" && st == "kick":
		return response.QQEventType_GroupRemoveType_QQEventType_GroupRemoveType_Kick, nil
	}
}

func (qegrh *qqevent_groupRemoveH) manual() (*response.Response_QQEvent_GroupRemove_Manual, error) {
	j := new(response_d.JSON_qqEvent_groupRemove_manual)
	if err := json.Unmarshal(qegrh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	var quiterName *string
	if qegrh.extra {
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		useridStr := strconv.FormatInt(j.UserId, 10)
		ctx, cancel := context.WithTimeout(define.ConnectorCtx, time.Second*15)
		defer cancel()
		ggmi, err := getGroupMemberInfo(ctx, groupidStr, useridStr, nil)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		quiterName = &ggmi.UserName
	}
	return &response.Response_QQEvent_GroupRemove_Manual{
		QuiterId:   strconv.FormatInt(j.UserId, 10),
		QuiterName: quiterName,
		GroupId:    strconv.FormatInt(j.GroupId, 10),
		Timestamp:  j.Timestamp,
		BotId:      strconv.FormatInt(j.SelfId, 10),
	}, nil
}

func (qegrh *qqevent_groupRemoveH) kick() (*response.Response_QQEvent_GroupRemove_Kick, error) {
	j := new(response_d.JSON_qqEvent_groupRemove_kick)
	if err := json.Unmarshal(qegrh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	var (
		quiterName   *string
		operatorName *string
	)
	if qegrh.extra {
		connectorClient := connector.NewConnectorClient(define.GRPCClient)
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		quiteridStr := strconv.FormatInt(j.UserId, 10)
		operatoridStr := strconv.FormatInt(j.OperatorId, 10)
		ctx, cancel := context.WithTimeout(define.ConnectorCtx, time.Second*15)
		defer cancel()
		stream, err := connectorClient.Read(ctx, &connector.Empty{})
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		quiter_ggmi, err := getGroupMemberInfo(ctx, groupidStr, quiteridStr, stream)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		operator_ggmi, err := getGroupMemberInfo(ctx, groupidStr, operatoridStr, stream)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		quiterName = &quiter_ggmi.UserName
		operatorName = &operator_ggmi.UserName
	}
	return &response.Response_QQEvent_GroupRemove_Kick{
		TargetId:     strconv.FormatInt(j.UserId, 10),
		TargetName:   quiterName,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		OperatorId:   strconv.FormatInt(j.OperatorId, 10),
		OperatorName: operatorName,
	}, nil
}

func (qemrh *qqevent_messageRecallH) MessageRecall() (*response.Response_QQEvent_MessageRecall, error) {
	t, err := qemrh.matchType()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	qemrh.d.Type = &t
	switch t {
	case response.QQEventType_MessageRecallType_QQEventType_MessageRecallType_Group:
		g, err := qemrh.group()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qemrh.d.Group = g
	case response.QQEventType_MessageRecallType_QQEventType_MessageRecallType_Private:
		p, err := qemrh.private()
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		qemrh.d.Private = p
	}
	return &qemrh.d, nil
}

func (qemrh *qqevent_messageRecallH) matchType() (response.QQEventType_MessageRecallType, error) {
	j := new(response_d.JSON_qqEventType)
	if err := json.Unmarshal(qemrh.buf, j); err != nil {
		logger.Println(err)
		return -1, err
	}
	if j.NoticeType == nil {
		return -1, errors.New("NoticeType为nil")
	}
	if j.SubType == nil {
		j.SubType = new(string)
	}
	switch mt, st := *j.NoticeType, *j.SubType; {
	default:
		return -1, fmt.Errorf("QQ事件类型消息撤回无匹配; NoticeType: %v, SubType: %v", mt, st)
	case mt == "group_recall":
		return response.QQEventType_MessageRecallType_QQEventType_MessageRecallType_Group, nil
	case mt == "friend_recall":
		return response.QQEventType_MessageRecallType_QQEventType_MessageRecallType_Private, nil
	}
}

func (qemrh *qqevent_messageRecallH) group() (*response.Response_QQEvent_MessageRecall_Group, error) {
	j := new(response_d.JSON_qqEvent_messageRecall_group)
	if err := json.Unmarshal(qemrh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	var (
		targetName   *string
		operatorName *string
	)
	if qemrh.extra {
		connectorClient := connector.NewConnectorClient(define.GRPCClient)
		groupidStr := strconv.FormatInt(j.GroupId, 10)
		ctx, cancel := context.WithTimeout(define.ConnectorCtx, time.Second*15)
		defer cancel()
		// 同一人
		if j.UserId == j.OperatorId {
			useridStr := strconv.FormatInt(j.UserId, 10)
			ggmi, err := getGroupMemberInfo(ctx, groupidStr, useridStr, nil)
			if err != nil {
				logger.Println(err)
				return nil, err
			}
			targetName = &ggmi.UserName
		} else { //不同人
			useridStr := strconv.FormatInt(j.UserId, 10)
			operatoridStr := strconv.FormatInt(j.OperatorId, 10)
			stream, err := connectorClient.Read(ctx, &connector.Empty{})
			if err != nil {
				logger.Println(err)
				return nil, err
			}
			user_ggmi, err := getGroupMemberInfo(ctx, groupidStr, useridStr, stream)
			if err != nil {
				logger.Println(err)
				return nil, err
			}
			targetName = &user_ggmi.UserName
			operator_ggmi, err := getGroupMemberInfo(ctx, groupidStr, operatoridStr, stream)
			if err != nil {
				logger.Println(err)
				return nil, err
			}
			operatorName = &operator_ggmi.UserName
		}
	}
	return &response.Response_QQEvent_MessageRecall_Group{
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

func (qemrh *qqevent_messageRecallH) private() (*response.Response_QQEvent_MessageRecall_Private, error) {
	j := new(response_d.JSON_qqEvent_messageRecall_private)
	if err := json.Unmarshal(qemrh.buf, j); err != nil {
		logger.Println(err)
		return nil, err
	}
	var recallerName *string
	ctx, cancel := context.WithTimeout(define.ConnectorCtx, time.Second*15)
	defer cancel()
	if qemrh.extra {
		useridStr := strconv.FormatInt(j.UserId, 10)
		gfi, err := getFriendInfo(ctx, useridStr, nil)
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		recallerName = &gfi.UserName
	}
	return &response.Response_QQEvent_MessageRecall_Private{
		RecallerId:   strconv.FormatInt(j.UserId, 10),
		RecallerName: recallerName,
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		MessageId:    strconv.FormatInt(j.MessageId, 10),
	}, nil
}

func sendCommand(ctx context.Context, buf []byte, requestEcho string, stream grpc.ServerStreamingClient[connector.ReadResponse]) (*cmdEventH, error) {
	connectorClient := connector.NewConnectorClient(define.GRPCClient)
	if stream == nil {
		x, err := connectorClient.Read(ctx, &connector.Empty{})
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		stream = x
	}
	writeCh := make(chan error, 1)
	readCh := make(chan any, 1)
	//写入
	go func() {
		if _, err := connectorClient.Write(ctx, &connector.WriteRequest{
			Buf: buf,
		}); err != nil {
			writeCh <- err
			logger.Println(err)
			return
		}
	}()
	//读取
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				readCh <- err
				logger.Println(err)
				return
			}
			respH := new(responseH)
			respH.buf = resp.Buf
			rt, err := respH.matchType()
			if err != nil {
				readCh <- err
				logger.Println(err)
				return
			}
			if rt != response.ResponseType_ResponseType_CmdEvent {
				continue
			}
			ceh := new(cmdEventH)
			ceh.buf = resp.Buf
			echo, err := ceh.Echo()
			if err != nil {
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
	case err := <-writeCh:
		return nil, err
	case x := <-readCh:
		switch x := x.(type) {
		case error:
			logger.Println(x)
			return nil, x
		case *cmdEventH:
			return x, nil
		default:
			return nil, errors.New("异常错误")
		}
	case <-ctx.Done():
		return nil, errors.New("获取额外用户信息超时")
	}
}

func getGroupMemberInfo(ctx context.Context, groupid, memberid string, stream grpc.ServerStreamingClient[connector.ReadResponse]) (*response.Response_CmdEvent_GetGroupMemberInfo, error) {
	requestEcho := strconv.FormatInt(rand.Int63(), 10)
	//获取写入内容
	buf, err := request.GetGroupMemberInfo(groupid, memberid, &requestEcho)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	ceh, err := sendCommand(ctx, buf, requestEcho, stream)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return ceh.GetGroupMemberInfo()
}

func getFriendInfo(ctx context.Context, friendid string, stream grpc.ServerStreamingClient[connector.ReadResponse]) (*response.Response_CmdEvent_GetFriendInfo, error) {
	connectorClient := connector.NewConnectorClient(define.GRPCClient)
	if stream == nil {
		x, err := connectorClient.Read(ctx, &connector.Empty{})
		if err != nil {
			logger.Println(err)
			return nil, err
		}
		stream = x
	}
	writeCh := make(chan error, 1)
	readCh := make(chan any, 1)
	requestEcho := strconv.FormatInt(rand.Int63(), 10)
	//获取写入内容
	buf, err := request.GetFriendInfo(friendid, &requestEcho)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	//写入
	go func() {
		if _, err := connectorClient.Write(ctx, &connector.WriteRequest{
			Buf: buf,
		}); err != nil {
			writeCh <- err
			logger.Println(err)
			return
		}
	}()
	//读取
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				readCh <- err
				logger.Println(err)
				return
			}
			respH := new(responseH)
			respH.buf = resp.Buf
			rt, err := respH.matchType()
			if err != nil {
				readCh <- err
				logger.Println(err)
				return
			}
			if rt != response.ResponseType_ResponseType_CmdEvent {
				continue
			}
			ceh := new(cmdEventH)
			ceh.buf = resp.Buf
			echo, err := ceh.Echo()
			if err != nil {
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
	case err := <-writeCh:
		logger.Println(err)
		return nil, err
	case x := <-readCh:
		switch x := x.(type) {
		case error:
			logger.Println(x)
			return nil, x
		case *cmdEventH:
			return x.GetFriendInfo()
		default:
			return nil, errors.New("异常错误")
		}
	case <-ctx.Done():
		return nil, errors.New("获取额外用户信息超时")
	}
}
