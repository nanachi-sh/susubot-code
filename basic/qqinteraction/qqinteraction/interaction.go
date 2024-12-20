package qqinteraction

import (
	"context"
	"regexp"

	"github.com/nanachi-sh/susubot-code/basic/qqinteraction/define"
	"github.com/nanachi-sh/susubot-code/basic/qqinteraction/log"
	"github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/connector"
	response_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/handler/response"
)

var logger = log.Get()

func Start() {
	stream, err := define.ConnectorC.Read(context.Background(), &connector.Empty{})
	if err != nil {
		logger.Fatalln(err)
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			logger.Fatalln(err)
		}
		respum, err := define.Handler_ResponseC.Unmarshal(define.HandlerCtx, &response_pb.UnmarshalRequest{
			Buf:            resp.Buf,
			ExtraInfo:      true,
			IgnoreCmdEvent: true,
		})
		if err != nil {
			logger.Fatalln(err)
		}
		switch *respum.Type {
		case response_pb.ResponseType_ResponseType_CmdEvent:
			continue
		case response_pb.ResponseType_ResponseType_Message:
			respum.Message
		case response_pb.ResponseType_ResponseType_QQEvent:
		}
	}
}

func message_match(text string) pluginType {
	switch {
	case randomanimal_match(text) != randomanimal_Unknown:
		return pluginType_RandomAnimal
	case randomfortune_match(text) != randomfortune_Unknown:
		return pluginType_RandomFortune
	case twoonone_match(text) != twoonone_Unknown:
		return pluginType_TwoOnOne
	default:
		return pluginType_Unknown
	}
}

type pluginType int

const (
	pluginType_Unknown pluginType = iota
	pluginType_RandomAnimal
	pluginType_RandomFortune
	pluginType_TwoOnOne
)

func randomanimal(message *response_pb.Response_Message) {
	if *message.Type != response_pb.MessageType_MessageType_Group {
		return
	}
	group := message.Group
	text := ""

	if text == "" {
		return
	}

}

func getText(mcs []*response_pb.MessageChainObject) string {
	for _, v := range mcs {
		if *v.Type == response_pb.MessageChainType_MessageChainType_Text {
			return v.Text.Text
			break
		}
	}
}

type randomanimalAction int

const (
	randomanimal_Unknown randomanimalAction = iota
	randomanimal_GetCat
	randomanimal_GetDog
	randomanimal_GetFox
	randomanimal_GetDuck
	randomanimal_GetChicken_CXK
)

func randomanimal_match(text string) randomanimalAction {
	switch text {
	case "来只猫", "来只猫猫", "来只小猫", "来只猫咪", "来只优蓝猫":
		return randomanimal_GetCat
	case "来只狗", "来只狗狗", "来只修狗", "来只小狗", "来只狗子", "来只狗老板":
		return randomanimal_GetDog
	case "来只狐狸", "来只狐", "来只狐狐", "来只小狐狸", "来只苏苏狐":
		return randomanimal_GetFox
	case "来只鸭子", "来只鸭", "来只鸭鸭":
		return randomanimal_GetDuck
	case "来只鸡", "来只坤坤", "来只坤", "来只只因":
		return randomanimal_GetChicken_CXK
	default:
		return randomanimal_Unknown
	}
}

func randomfortune() {

}

type randomfortuneAction int

const (
	randomfortune_Unknown randomfortuneAction = iota
	randomfortune_GetFortune
)

func randomfortune_match(text string) randomfortuneAction {
	switch text {
	default:
		return randomfortune_Unknown
	case "#抽签":
		return randomfortune_GetFortune
	}
}

func twoonone() {

}

type twoononeAction int

const (
	twoonone_Unknown twoononeAction = iota
	twoonone_CreateAccount
	twoonone_GetAccount
	twoonone_CreateRoom
	twoonone_GetDaliyCoin
	twoonone_RobLandowner_Rob
	twoonone_RobLandowner_NoRob
	twoonone_ExitRoom
	twoonone_StartRoom
	twoonone_GetRooms
	twoonone_SendCard_NoSend
	twoonone_SendCard_Send
	twoonone_JoinRoom
)

func twoonone_match(text string) twoononeAction {
	switch {
	case "开号":
		return twoonone_CreateAccount
	case "个人信息":
		return twoonone_GetAccount
	case "开桌":
		return twoonone_CreateRoom
	case "领豆子":
		return twoonone_GetDaliyCoin
	case "抢地主", "抢", "我抢":
		return twoonone_RobLandowner_Rob
	case "不抢", "不抢地主":
		return twoonone_RobLandowner_NoRob
	case "下桌":
		return twoonone_ExitRoom
	case "桌列表":
		return twoonone_GetRooms
	case "不要", "要不起", "不出":
		return twoonone_SendCard_NoSend
	case "发牌":
		return twoonone_StartRoom
	}
	if ok, _ := regexp.MatchString(`\A(!|！)([3456789jqkaJQKA2]|10|大王|小王)+`, text); ok {
		return twoonone_SendCard_Send
	}
	if ok, _ := regexp.MatchString(`\A上桌([(0-9)]| ){3,}`, text); ok {
		return twoonone_JoinRoom
	}
	return twoonone_Unknown
}
