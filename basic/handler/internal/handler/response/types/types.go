package types

type MessageChain struct {
	Data map[string]any `json:"data"`
	Type string         `json:"type"`
}

type JSON_responseType struct {
	PostType *string `json:"post_type"`
	Echo     *string `json:"echo"`
}

type (
	JSON_botEventType struct {
		SubType       string `json:"sub_type"`
		MetaEventType string `json:"meta_event_type"`
	}

	JSON_botEvent_Connected struct {
		Timestamp int64 `json:"time"`
		SelfId    int64 `json:"self_id"`
	}

	JSON_botEvent_HeartPacket struct {
		Timestamp int64                             `json:"time"`
		SelfId    int64                             `json:"self_id"`
		Status    *JSON_botEvent_HeartPacket_Status `json:"status"`
		Interval  int64                             `json:"interval"` //ms
	}
	JSON_botEvent_HeartPacket_Status struct {
		Online bool `json:"online"`
		Good   bool `json:"good"`
	}
)

const (
	JSON_cmdEvent_Status_OK     = "ok"
	JSON_cmdEvent_Status_Failed = "failed"
)

type (
	JSON_cmdEvent_Echo struct {
		Echo string `json:"echo"`
	}

	JSON_cmdEvent_GetFriendInfo struct {
		Status   string `json:"status"`
		Retcode  int    `json:"retcode"`
		UserId   int64  `json:"user_id"`
		NickName string `json:"nickname"`
		Remark   string `json:"remark"`
	}

	JSON_cmdEvent_GetFriendList struct {
		Status  string                              `json:"status"`
		Retcode int                                 `json:"retcode"`
		Message string                              `json:"message"`
		Wording string                              `json:"wording"`
		Echo    string                              `json:"echo"`
		Data    []*JSON_cmdEvent_GetFriendList_Data `json:"data"`
	}
	JSON_cmdEvent_GetFriendList_Data struct {
		User_Id  int64  `json:"user_id"`
		NickName string `json:"nickname"`
		REmark   string `json:"remark"`
		Sex      string `json:"sex"`
		Level    int    `json:"level"`
	}

	JSON_cmdEvent_GetGroupMemberInfo struct {
		Status  string                                 `json:"status"`
		Retcode int                                    `json:"retcode"`
		Member  *JSON_cmdEvent_GetGroupMemberInfo_Data `json:"data"`
	}
	JSON_cmdEvent_GetGroupMemberInfo_Data struct {
		UserId           int64  `json:"user_id"`
		GroupId          int64  `json:"group_id"`
		UserName         string `json:"nickname"`
		Title            string `json:"title"` //头衔
		Title_ExpireTime int64  `json:"title_expire_time"`
		Card             string `json:"card"`
		JoinTime         int64  `json:"join_time"`
		LastActiveTime   int64  `json:"last_active_time"`
		LastSentTime     int64  `json:"last_sent_time"`
		Role             string `json:"role"`
		Sex              string `json:"sex"`
	}

	JSON_cmdEvent_GetMessage struct {
		Status  string `json:"status"`
		Retcode int    `json:"retcode"`
		Data    *JSON_cmdEvent_GetMessage_Data
	}
	JSON_cmdEvent_GetMessage_Data struct {
		MessageChain []*JSON_cmdEvent_GetMessage_Data_MessageChain `json:"message"`
		MessageType  string                                        `json:"message_type"`
		SubType      string                                        `json:"sub_type"`
		Timestamp    int64                                         `json:"time"`
		MessageId    int64                                         `json:"message_id"`
		Sender       *JSON_cmdEvent_GetMessage_Data_Sender         `json:"sender"`
		GroupId      int64                                         `json:"group_id"`
		SelfId       int64                                         `json:"self_id"`
	}
	JSON_cmdEvent_GetMessage_Data_MessageChain struct {
		Data map[string]any `json:"data"`
		Type string         `json:"type"`
	}
	JSON_cmdEvent_GetMessage_Data_Sender struct {
		Nickname string `json:"nickname"`
		Card     string `json:"card"`
		Role     string `json:"role"`
		UserId   int64  `json:"user_id"`
	}

	JSON_cmdEvent_GetGroupInfo struct {
		Status  string                           `json:"status"`
		Retcode int                              `json:"retcode"`
		Group   *JSON_cmdEvent_GetGroupInfo_Data `json:"data"`
	}
	JSON_cmdEvent_GetGroupInfo_Data struct {
		GroupId   int64  `json:"group_id"`
		GroupName string `json:"group_name"`
		MemberMax int    `json:"max_member"`
		MemberNow int    `json:"member_num"`
	}
)

type (
	JSON_messageType struct {
		SubType     *string `json:"sub_type"`
		MessageType *string `json:"message_type"`
	}
	JSON_message_MessageChain struct {
		Data map[string]any `json:"data"`
		Type string         `json:"type"`
	}

	JSON_message_Group struct {
		MessageId    int64                        `json:"message_id"`
		PeerId       int64                        `json:"peer_id"`
		MessageChain []*JSON_message_MessageChain `json:"message"`
		Timestamp    int64                        `json:"time"`
		SubType      string                       `json:"sub_type"`
		GroupId      int64                        `json:"group_id"`
		UserId       int64                        `json:"user_id"`
		MainType     string                       `json:"message_type"`
		SelfId       int64                        `json:"self_id"`
		Sender       *JSON_message_Group_Sender   `json:"sender"`
	}
	JSON_message_Group_Sender struct {
		Nickname string `json:"nickname"`
		Card     string `json:"card"`
		Role     string `json:"role"`
		UserId   int64  `json:"user_id"`
	}

	JSON_message_Private struct {
		MessageChain []*JSON_message_MessageChain `json:"message"`
		Sender       *JSON_message_Private_Sender `json:"sender"`
		MessageId    int64                        `json:"message_id"`
		Timestamp    int64                        `json:"time"`
		SelfId       int64                        `json:"self_id"`
		MainType     string                       `json:"message_type"`
		SubType      string                       `json:"sub_type"`
	}
	JSON_message_Private_Sender struct {
		Nickname string `json:"nickname"`
		UserId   int64  `json:"user_id"`
	}
)

type (
	JSON_qqEventType struct {
		NoticeType *string `json:"notice_type"`
		SubType    *string `json:"sub_type"`
	}

	JSON_qqEvent_groupAdd_direct struct {
		OperatorId int64 `json:"operator_id"`
		SelfId     int64 `json:"self_id"`
		UserId     int64 `json:"user_id"`
		Timestamp  int64 `json:"time"`
		GroupId    int64 `json:"group_id"`
	}

	JSON_qqEvent_groupRemove_manual struct {
		SelfId    int64 `json:"self_id"`
		UserId    int64 `json:"user_id"`
		Timestamp int64 `json:"time"`
		GroupId   int64 `json:"group_id"`
	}
	JSON_qqEvent_groupRemove_kick struct {
		OperatorId int64 `json:"operator_id"`
		SelfId     int64 `json:"self_id"`
		UserId     int64 `json:"user_id"`
		Timestamp  int64 `json:"time"`
		GroupId    int64 `json:"group_id"`
	}
	JSON_qqEvent_groupMute struct {
		OperatorId int64 `json:"operator_id"`
		SelfId     int64 `json:"self_id"`
		UserId     int64 `json:"user_id"`
		Timestamp  int64 `json:"time"`
		GroupId    int64 `json:"group_id"`
		Duration   int64 `json:"duration"`
	}
	JSON_qqEvent_groupUnmute struct {
		OperatorId int64 `json:"operator_id"`
		SelfId     int64 `json:"self_id"`
		UserId     int64 `json:"user_id"`
		Timestamp  int64 `json:"time"`
		GroupId    int64 `json:"group_id"`
		Duration   int64 `json:"duration"`
	}
	JSON_qqEvent_messageRecall_group struct {
		OperatorId int64 `json:"operator_id"`
		SelfId     int64 `json:"self_id"`
		UserId     int64 `json:"user_id"`
		Timestamp  int64 `json:"time"`
		GroupId    int64 `json:"group_id"`
		MessageId  int64 `json:"message_id"`
	}
	JSON_qqEvent_messageRecall_private struct {
		SelfId    int64 `json:"self_id"`
		UserId    int64 `json:"user_id"`
		Timestamp int64 `json:"time"`
		MessageId int64 `json:"message_id"`
	}
)
