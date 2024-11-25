// 机器人核心响应处理
package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/nanachi-sh/susubot-code/handler/protos/handler"
	"github.com/nanachi-sh/susubot-code/handler/response/define"
)

type responseH struct {
	buf   []byte
	rtype handler.ResponseType

	cet *handler.CmdEventType
}

type botEventH struct {
	d   handler.R_BotEvent
	buf []byte
}

type (
	cmdEventH struct {
		d   handler.R_CmdEvent
		buf []byte
	}

	cmdEvent_GetMessageH struct {
		j   *define.JSON_cmdEvent_GetMessage
		d   handler.CE_GetMessage
		buf []byte
	}
)

type messageH struct {
	d   handler.R_Message
	buf []byte
}

type (
	qqeventH struct {
		d   handler.R_QQEvent
		buf []byte
	}

	qqevent_groupAddH struct {
		d   handler.QQE_GroupAdd
		buf []byte
	}

	qqevent_groupRemoveH struct {
		d   handler.QQE_GroupRemove
		buf []byte
	}

	qqevent_messageRecallH struct {
		d   handler.QQE_MessageRecall
		buf []byte
	}
)

func New(buf []byte, rt *handler.ResponseType, cet *handler.CmdEventType) (*responseH, error) {
	if rt == nil {
		ret, err := matchType(buf)
		if err != nil {
			return nil, err
		}
		rt = &ret
	}
	rh := &responseH{
		buf:   buf,
		rtype: *rt,
		cet:   cet,
	}
	return rh, nil
}

func unmarshalMessageChain(mc []*define.MessageChain) ([]*handler.MessageChainObject, error) {
	ret := []*handler.MessageChainObject{}
	for _, json := range mc {
		switch json.Type {
		case "text":
			d, ok := json.Data["text"]
			if !ok {
				return nil, errors.New("text不存在")
			}
			switch d := d.(type) {
			case string:
				ret = append(ret, &handler.MessageChainObject{
					Type: handler.MessageChainType_MCT_Text.Enum(),
					Text: &handler.MC_Text{
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
				ret = append(ret, &handler.MessageChainObject{
					Type: handler.MessageChainType_MCT_At.Enum(),
					At: &handler.MC_At{
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
				ret = append(ret, &handler.MessageChainObject{
					Type: handler.MessageChainType_MCT_Reply.Enum(),
					Reply: &handler.MC_Reply{
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
				ret = append(ret, &handler.MessageChainObject{
					Type: handler.MessageChainType_MCT_Image.Enum(),
					Image: &handler.MC_Image{
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
				ret = append(ret, &handler.MessageChainObject{
					Type: handler.MessageChainType_MCT_Voice.Enum(),
					Voice: &handler.MC_Voice{
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
				ret = append(ret, &handler.MessageChainObject{
					Type: handler.MessageChainType_MCT_Video.Enum(),
					Video: &handler.MC_Video{
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

func matchType(buf []byte) (handler.ResponseType, error) {
	j := new(define.JSON_responseType)
	if err := json.Unmarshal(buf, j); err != nil {
		return -1, err
	}
	if j.PostType == nil {
		return handler.ResponseType_RT_CmdEvent, nil
	} else {
		switch pt := *j.PostType; pt {
		case "message":
			return handler.ResponseType_RT_Message, nil
		case "notice":
			return handler.ResponseType_RT_QQEvent, nil
		case "meta_event":
			return handler.ResponseType_RT_BotEvent, nil
		default:
			return -1, fmt.Errorf("响应事件类型无匹配; PostType: %v", pt)
		}
	}
}

func (rh *responseH) MarshalToResponse() (*handler.BotResponseUnmarshalResponse, error) {
	ret := new(handler.BotResponseUnmarshalResponse)
	ret.Type = &rh.rtype
	switch rh.rtype {
	case handler.ResponseType_RT_BotEvent:
		be, err := rh.BotEvent()
		if err != nil {
			return nil, err
		}
		ret.BotEvent = be
	case handler.ResponseType_RT_CmdEvent:
		ce, err := rh.CmdEvent()
		if err != nil {
			return nil, err
		}
		ret.CmdEvent = ce
	case handler.ResponseType_RT_Message:
		m, err := rh.Message()
		if err != nil {
			return nil, err
		}
		ret.Message = m
	case handler.ResponseType_RT_QQEvent:
		qqe, err := rh.QQEvent()
		if err != nil {
			return nil, err
		}
		ret.QQEvent = qqe
	}
	return ret, nil
}

func (rh *responseH) BotEvent() (*handler.R_BotEvent, error) {
	beh := new(botEventH)
	beh.buf = rh.buf
	t, err := beh.matchType()
	if err != nil {
		return nil, err
	}
	beh.d.Type = &t
	switch t {
	case handler.BotEventType_BET_Connected:
		ret, err := beh.Connected()
		if err != nil {
			return nil, err
		}
		beh.d.Connected = ret
	case handler.BotEventType_BET_HeartPacket:
		ret, err := beh.HeartPacket()
		if err != nil {
			return nil, err
		}
		beh.d.HeartPacket = ret
	default:
		return nil, fmt.Errorf("机器人响应类型无匹配; MetaEventType: %v", beh.d.Type)
	}
	return &beh.d, nil
}

func (beh *botEventH) matchType() (handler.BotEventType, error) {
	j := new(define.JSON_botEventType)
	if err := json.Unmarshal(beh.buf, j); err != nil {
		return -1, err
	}
	met, st := j.MetaEventType, j.SubType
	switch {
	case met == "lifecycle" && st == "connect":
		return handler.BotEventType_BET_Connected, nil
	case met == "heartbeat":
		return handler.BotEventType_BET_HeartPacket, nil
	default:
		return -1, fmt.Errorf("机器人响应类型无匹配; MetaEventType: %v, SubType: %v", met, st)
	}
}

func (beh *botEventH) Connected() (*handler.BE_Connected, error) {
	j := new(define.JSON_botEvent_Connected)
	if err := json.Unmarshal(beh.buf, j); err != nil {
		return nil, err
	}
	return &handler.BE_Connected{
		Timestamp: j.Timestamp,
		BotId:     strconv.FormatInt(j.SelfId, 10),
	}, nil
}

func (beh *botEventH) HeartPacket() (*handler.BE_HeartPacket, error) {
	j := new(define.JSON_botEvent_HeartPacket)
	if err := json.Unmarshal(beh.buf, j); err != nil {
		return nil, err
	}
	return &handler.BE_HeartPacket{
		Timestamp: j.Timestamp,
		BotId:     strconv.FormatInt(j.SelfId, 10),
		Interval:  j.Interval,
		Status: &handler.BE_HeartPacketStatus{
			Online: j.Status.Online,
			Good:   j.Status.Good,
		},
	}, nil
}

func (rh *responseH) CmdEvent() (*handler.R_CmdEvent, error) {
	ceh := new(cmdEventH)
	ceh.buf = rh.buf
	if rh.cet == nil {
		return nil, errors.New("未指定命令响应类型")
	}
	ceh.d.Type = rh.cet
	switch *rh.cet {
	case handler.CmdEventType_CET_GetFriendList:
		gfl, err := ceh.GetFriendList()
		if err != nil {
			return nil, err
		}
		ceh.d.GetFriendList = gfl
	case handler.CmdEventType_CET_GetGroupInfo:
		ggi, err := ceh.GetGroupInfo()
		if err != nil {
			return nil, err
		}
		ceh.d.GetGroupInfo = ggi
	case handler.CmdEventType_CET_GetGroupMemberInfo:
		ggmi, err := ceh.GetGroupMemberInfo()
		if err != nil {
			return nil, err
		}
		ceh.d.GetGroupMemberInfo = ggmi
	case handler.CmdEventType_CET_GetMessage:
		gm, err := ceh.GetMessage()
		if err != nil {
			return nil, err
		}
		ceh.d.GetMessage = gm
	}
	return &ceh.d, nil
}

func (ceh *cmdEventH) GetFriendList() (*handler.CE_GetFriendList, error) {
	j := new(define.JSON_cmdEvent_GetFriendList)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		return nil, err
	}
	ok := false
	if j.Status == define.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(handler.CE_GetFriendList)
	if ok {
		ret.OK = true
		ret.Retcode = nil
		friends := []*handler.CE_GetFriendList_FriendInfo{}
		for _, v := range j.Data {
			remark := new(string)
			if v.REmark == "" {
				remark = nil
			} else {
				remark = &v.REmark
			}
			friends = append(friends, &handler.CE_GetFriendList_FriendInfo{
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

func (ceh *cmdEventH) GetGroupInfo() (*handler.CE_GetGroupInfo, error) {
	j := new(define.JSON_cmdEvent_GetGroupInfo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		return nil, err
	}
	ok := false
	if j.Status == define.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(handler.CE_GetGroupInfo)
	if ok {
		ret = &handler.CE_GetGroupInfo{
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

func (ceh *cmdEventH) GetGroupMemberInfo() (*handler.CE_GetGroupMemberInfo, error) {
	j := new(define.JSON_cmdEvent_GetGroupMemberInfo)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		return nil, err
	}
	ok := false
	if j.Status == define.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(handler.CE_GetGroupMemberInfo)
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
		var role handler.GroupRole
		switch j.Member.Role {
		case define.Role_Member:
			role = handler.GroupRole_GR_Member
		case define.Role_Admin:
			role = handler.GroupRole_GR_Admin
		case define.Role_Owner:
			role = handler.GroupRole_GR_Owner
		}
		ret = &handler.CE_GetGroupMemberInfo{
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

func (ceh *cmdEventH) GetMessage() (*handler.CE_GetMessage, error) {
	j := new(define.JSON_cmdEvent_GetMessage)
	if err := json.Unmarshal(ceh.buf, j); err != nil {
		return nil, err
	}
	ok := false
	if j.Status == define.JSON_cmdEvent_Status_OK {
		ok = true
	}
	ret := new(handler.CE_GetMessage)
	if ok {
		cegmh := &cmdEvent_GetMessageH{
			j:   nil,
			d:   handler.CE_GetMessage{},
			buf: ceh.buf,
		}
		m, err := cegmh.Message()
		if err != nil {
			return nil, err
		}
		ret = &handler.CE_GetMessage{
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

func (cegmh *cmdEvent_GetMessageH) matchType() (handler.MessageType, error) {
	j := cegmh.j
	switch mt, st := j.Data.MessageType, j.Data.SubType; {
	case mt == "private" && st == "friend":
		return handler.MessageType_MT_Private, nil
	case mt == "group" && st == "normal":
		return handler.MessageType_MT_Group, nil
	default:
		return -1, fmt.Errorf("命令响应获取消息事件类型无匹配; MessageType: %v, SubType: %v", mt, st)
	}
}

func (cegmh *cmdEvent_GetMessageH) Message() (*handler.GM_Message, error) {
	if cegmh.j == nil {
		j := new(define.JSON_cmdEvent_GetMessage)
		if err := json.Unmarshal(cegmh.buf, j); err != nil {
			return nil, err
		}
		cegmh.j = j
	}
	m := new(handler.GM_Message)
	t, err := cegmh.matchType()
	if err != nil {
		return nil, err
	}
	m.Type = &t
	switch t {
	case handler.MessageType_MT_Group:
		g, err := cegmh.group()
		if err != nil {
			return nil, err
		}
		m.Group = g
	case handler.MessageType_MT_Private:
		p, err := cegmh.private()
		if err != nil {
			return nil, err
		}
		m.Private = p
	}
	return m, nil
}

func (cegmh *cmdEvent_GetMessageH) group() (*handler.GMM_Group, error) {
	j := cegmh.j
	sname := &j.Data.Sender.Nickname
	var srole, brole *handler.GroupRole
	if j.Data.Sender.Role != "" {
		switch j.Data.Sender.Role {
		case define.Role_Member:
			srole = handler.GroupRole_GR_Member.Enum()
		case define.Role_Admin:
			srole = handler.GroupRole_GR_Admin.Enum()
		case define.Role_Owner:
			srole = handler.GroupRole_GR_Owner.Enum()
		}
	}
	jmc := []*define.MessageChain{}
	for _, v := range j.Data.MessageChain {
		jmc = append(jmc, &define.MessageChain{
			Data: v.Data,
			Type: v.Type,
		})
	}
	mc, err := unmarshalMessageChain(jmc)
	if err != nil {
		return nil, err
	}
	return &handler.GMM_Group{
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

func (cegmh *cmdEvent_GetMessageH) private() (*handler.GMM_Private, error) {
	return nil, errors.ErrUnsupported
}

func (rh *responseH) Message() (*handler.R_Message, error) {
	mh := new(messageH)
	mh.buf = rh.buf
	t, err := mh.matchType()
	if err != nil {
		return nil, err
	}
	mh.d.Type = &t
	switch t {
	case handler.MessageType_MT_Group:
		g, err := mh.group()
		if err != nil {
			return nil, err
		}
		mh.d.Group = g
	case handler.MessageType_MT_Private:
		p, err := mh.private()
		if err != nil {
			return nil, err
		}
		mh.d.Private = p
	}
	return &mh.d, nil
}

func (mh *messageH) matchType() (handler.MessageType, error) {
	j := new(define.JSON_messageType)
	if err := json.Unmarshal(mh.buf, j); err != nil {
		return -1, err
	}
	if j.MessageType == nil || j.SubType == nil {
		return -1, errors.New("MessageType/SubType为nil")
	}
	switch mt, st := *j.MessageType, *j.SubType; {
	case mt == "private" && st == "friend":
		return handler.MessageType_MT_Private, nil
	case mt == "group" && st == "normal":
		return handler.MessageType_MT_Group, nil
	default:
		return -1, fmt.Errorf("消息事件类型无匹配; MessageType: %v, SubType: %v", mt, st)
	}
}

func (mh *messageH) group() (*handler.M_Group, error) {
	j := new(define.JSON_message_Group)
	if err := json.Unmarshal(mh.buf, j); err != nil {
		return nil, err
	}
	sname := &j.Sender.Nickname
	var srole, brole *handler.GroupRole
	if j.Sender.Role != "" {
		switch j.Sender.Role {
		case define.Role_Member:
			srole = handler.GroupRole_GR_Member.Enum()
		case define.Role_Admin:
			srole = handler.GroupRole_GR_Admin.Enum()
		case define.Role_Owner:
			srole = handler.GroupRole_GR_Owner.Enum()
		}
	}
	jmc := []*define.MessageChain{}
	for _, v := range j.MessageChain {
		jmc = append(jmc, &define.MessageChain{
			Data: v.Data,
			Type: v.Type,
		})
	}
	mc, err := unmarshalMessageChain(jmc)
	if err != nil {
		return nil, err
	}
	return &handler.M_Group{
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

func (mh *messageH) private() (*handler.M_Private, error) {
	j := new(define.JSON_message_Private)
	if err := json.Unmarshal(mh.buf, j); err != nil {
		return nil, err
	}
	sname := &j.Sender.Nickname
	jmc := []*define.MessageChain{}
	for _, v := range j.MessageChain {
		jmc = append(jmc, &define.MessageChain{
			Data: v.Data,
			Type: v.Type,
		})
	}
	mc, err := unmarshalMessageChain(jmc)
	if err != nil {
		return nil, err
	}
	return &handler.M_Private{
		SenderId:     strconv.FormatInt(j.Sender.UserId, 10),
		SenderName:   sname,
		MessageId:    strconv.FormatInt(j.MessageId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		MessageChain: mc,
	}, nil
}

func (rh *responseH) QQEvent() (*handler.R_QQEvent, error) {
	qeh := new(qqeventH)
	qeh.buf = rh.buf
	t, err := qeh.matchType()
	if err != nil {
		return nil, err
	}
	qeh.d.Type = &t
	switch t {
	case handler.QQEventType_QQET_GroupAdd:
		qegah := new(qqevent_groupAddH)
		qegah.buf = rh.buf
		ga, err := qegah.GroupAdd()
		if err != nil {
			return nil, err
		}
		qeh.d.GroupAdd = ga
	case handler.QQEventType_QQET_GroupRemove:
		qegrh := new(qqevent_groupRemoveH)
		qegrh.buf = rh.buf
		gr, err := qegrh.GroupRemove()
		if err != nil {
			return nil, err
		}
		qeh.d.GroupRemove = gr
	case handler.QQEventType_QQET_GroupMute:
		gm, err := qeh.groupMute()
		if err != nil {
			return nil, err
		}
		qeh.d.GroupMute = gm
	case handler.QQEventType_QQET_GroupUnmute:
		gum, err := qeh.groupUnmute()
		if err != nil {
			return nil, err
		}
		qeh.d.GroupUnmute = gum
	case handler.QQEventType_QQET_MessageRecall:
		qemrh := new(qqevent_messageRecallH)
		qemrh.buf = rh.buf
		mr, err := qemrh.MessageRecall()
		if err != nil {
			return nil, err
		}
		qeh.d.MessageRecall = mr
	}
	return &qeh.d, nil
}

func (qeh *qqeventH) matchType() (handler.QQEventType, error) {
	j := new(define.JSON_qqEventType)
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
		return handler.QQEventType_QQET_GroupAdd, nil
	case mt == "group_decrease":
		return handler.QQEventType_QQET_GroupRemove, nil
	case mt == "group_ban" && st == "ban":
		return handler.QQEventType_QQET_GroupMute, nil
	case mt == "group_ban" && st == "lift_ban":
		return handler.QQEventType_QQET_GroupUnmute, nil
	case mt == "group_recall", mt == "friend_recall":
		return handler.QQEventType_QQET_MessageRecall, nil
	}
}

func (qeh *qqeventH) groupMute() (*handler.QQE_GroupMute, error) {
	j := new(define.JSON_qqEvent_groupMute)
	if err := json.Unmarshal(qeh.buf, j); err != nil {
		return nil, err
	}
	return &handler.QQE_GroupMute{
		TargetId:     strconv.FormatInt(j.UserId, 10),
		TargetName:   nil,
		Timestamp:    j.Timestamp,
		OperatorId:   strconv.FormatInt(j.OperatorId, 10),
		OperatorName: nil,
		Duration:     int32(j.Duration),
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		BotId:        strconv.FormatInt(j.SelfId, 10),
	}, nil
}

func (qeh *qqeventH) groupUnmute() (*handler.QQE_GroupUnmute, error) {
	j := new(define.JSON_qqEvent_groupUnmute)
	if err := json.Unmarshal(qeh.buf, j); err != nil {
		return nil, err
	}
	return &handler.QQE_GroupUnmute{
		TargetId:     strconv.FormatInt(j.UserId, 10),
		TargetName:   nil,
		Timestamp:    j.Timestamp,
		OperatorId:   strconv.FormatInt(j.OperatorId, 10),
		OperatorName: nil,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		BotId:        strconv.FormatInt(j.SelfId, 10),
	}, nil
}

// WIP
func (qegah *qqevent_groupAddH) matchType() (handler.QQE_GroupAddType, error) {
	j := new(define.JSON_qqEventType)
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
		return handler.QQE_GroupAddType_GAT_Direct, nil
	}
}

func (qegah *qqevent_groupAddH) GroupAdd() (*handler.QQE_GroupAdd, error) {
	t, err := qegah.matchType()
	if err != nil {
		return nil, err
	}
	qegah.d.Type = &t
	switch t {
	case handler.QQE_GroupAddType_GAT_Direct:
		d, err := qegah.direct()
		if err != nil {
			return nil, err
		}
		qegah.d.Direct = d
	case handler.QQE_GroupAddType_GAT_Invite:
		i, err := qegah.invite()
		if err != nil {
			return nil, err
		}
		qegah.d.Invite = i
	}
	return &qegah.d, nil
}

func (qegah *qqevent_groupAddH) direct() (*handler.GA_Direct, error) {
	j := new(define.JSON_qqEvent_groupAdd_direct)
	if err := json.Unmarshal(qegah.buf, j); err != nil {
		return nil, err
	}
	return &handler.GA_Direct{
		JoinerId:     strconv.FormatInt(j.UserId, 10),
		JoinerName:   nil,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		ApproverId:   strconv.FormatInt(j.OperatorId, 10),
		ApproverName: nil,
	}, nil
}

func (qegah *qqevent_groupAddH) invite() (*handler.GA_Invite, error) {
	return nil, errors.ErrUnsupported
}

func (qegrh *qqevent_groupRemoveH) GroupRemove() (*handler.QQE_GroupRemove, error) {
	t, err := qegrh.matchType()
	if err != nil {
		return nil, err
	}
	qegrh.d.Type = &t
	switch t {
	case handler.QQE_GroupRemoveType_GRT_Kick:
		k, err := qegrh.kick()
		if err != nil {
			return nil, err
		}
		qegrh.d.Kick = k
	case handler.QQE_GroupRemoveType_GRT_Manual:
		m, err := qegrh.manual()
		if err != nil {
			return nil, err
		}
		qegrh.d.Manual = m
	}
	return &qegrh.d, nil
}

func (qegrh *qqevent_groupRemoveH) matchType() (handler.QQE_GroupRemoveType, error) {
	j := new(define.JSON_qqEventType)
	if err := json.Unmarshal(qegrh.buf, j); err != nil {
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
		return handler.QQE_GroupRemoveType_GRT_Manual, nil
	case mt == "group_decrease" && st == "kick":
		return handler.QQE_GroupRemoveType_GRT_Kick, nil
	}
}

func (qegrh *qqevent_groupRemoveH) manual() (*handler.GR_Manual, error) {
	j := new(define.JSON_qqEvent_groupRemove_manual)
	if err := json.Unmarshal(qegrh.buf, j); err != nil {
		return nil, err
	}
	return &handler.GR_Manual{
		QuiterId:   strconv.FormatInt(j.UserId, 10),
		QuiterName: nil,
		GroupId:    strconv.FormatInt(j.GroupId, 10),
		Timestamp:  j.Timestamp,
		BotId:      strconv.FormatInt(j.SelfId, 10),
	}, nil
}

func (qegrh *qqevent_groupRemoveH) kick() (*handler.GR_Kick, error) {
	j := new(define.JSON_qqEvent_groupRemove_kick)
	if err := json.Unmarshal(qegrh.buf, j); err != nil {
		return nil, err
	}
	return &handler.GR_Kick{
		TargetId:     strconv.FormatInt(j.UserId, 10),
		TargetName:   nil,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		OperatorId:   strconv.FormatInt(j.OperatorId, 10),
		OperatorName: nil,
	}, nil
}

func (qemrh *qqevent_messageRecallH) MessageRecall() (*handler.QQE_MessageRecall, error) {
	t, err := qemrh.matchType()
	if err != nil {
		return nil, err
	}
	qemrh.d.Type = &t
	switch t {
	case handler.QQE_MessageRecallType_MRT_Group:
		g, err := qemrh.group()
		if err != nil {
			return nil, err
		}
		qemrh.d.Group = g
	case handler.QQE_MessageRecallType_MRT_Private:
		p, err := qemrh.private()
		if err != nil {
			return nil, err
		}
		qemrh.d.Private = p
	}
	return &qemrh.d, nil
}

func (qemrh *qqevent_messageRecallH) matchType() (handler.QQE_MessageRecallType, error) {
	j := new(define.JSON_qqEventType)
	if err := json.Unmarshal(qemrh.buf, j); err != nil {
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
		return handler.QQE_MessageRecallType_MRT_Group, nil
	case mt == "friend_recall":
		return handler.QQE_MessageRecallType_MRT_Private, nil
	}
}

func (qemrh *qqevent_messageRecallH) group() (*handler.MR_Group, error) {
	j := new(define.JSON_qqEvent_messageRecall_group)
	if err := json.Unmarshal(qemrh.buf, j); err != nil {
		return nil, err
	}
	return &handler.MR_Group{
		TargetId:     strconv.FormatInt(j.UserId, 10),
		TargetName:   nil,
		Timestamp:    j.Timestamp,
		OperatorId:   strconv.FormatInt(j.OperatorId, 10),
		OperatorName: nil,
		GroupId:      strconv.FormatInt(j.GroupId, 10),
		BotId:        strconv.FormatInt(j.SelfId, 10),
		MessageId:    strconv.FormatInt(j.MessageId, 10),
	}, nil
}

func (qemrh *qqevent_messageRecallH) private() (*handler.MR_Private, error) {
	j := new(define.JSON_qqEvent_messageRecall_private)
	if err := json.Unmarshal(qemrh.buf, j); err != nil {
		return nil, err
	}
	return &handler.MR_Private{
		RecallerId:   strconv.FormatInt(j.UserId, 10),
		RecallerName: nil,
		Timestamp:    j.Timestamp,
		BotId:        strconv.FormatInt(j.SelfId, 10),
		MessageId:    strconv.FormatInt(j.MessageId, 10),
	}, nil
}
