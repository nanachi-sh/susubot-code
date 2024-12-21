package qqinteraction

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/qqinteraction/define"
	"github.com/nanachi-sh/susubot-code/basic/qqinteraction/log"
	connector_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/connector"
	request_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/handler/request"
	response_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/handler/response"
	randomanimal_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/randomanimal"
	randomfortune_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/randomfortune"
	twoonone_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/twoonone"
)

var logger = log.Get()

var (
	twoonone_rooms        = make(map[string]*roomSI) //id To room
	twoonone_player2room  = make(map[string]*roomSI) //
	twoonone_playerStatus = make(map[string]struct{})
)

type roomSI struct {
	hash           string
	id             string
	landownerCards []twoonone_pb.Card
}

func Start() {
	stream, err := define.ConnectorC.Read(define.ConnectorCtx, &connector_pb.Empty{})
	if err != nil {
		logger.Fatalln(err)
	}
	rs, err := define.TwoOnOneC.GetRooms(define.TwoOnOneCtx, &twoonone_pb.Empty{})
	if err != nil {
		logger.Println(err)
	} else {
		for _, v := range rs.Rooms {
			var r *roomSI
			for {
				id := randomString(3, OnlyNumber)
				if _, ok := twoonone_rooms[id]; ok {
					continue
				} else {
					var loCards []twoonone_pb.Card
					if len(v.LandownerCards) > 0 {
						loCards = v.LandownerCards
					}
					r = &roomSI{
						hash:           v.Hash,
						id:             id,
						landownerCards: loCards,
					}
					twoonone_rooms[id] = r
					break
				}
			}
			for _, v2 := range v.Players {
				twoonone_player2room[v2.AccountInfo.Id] = r
			}
		}
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			logger.Fatalln(err)
		}
		go func() {
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
				return
			case response_pb.ResponseType_ResponseType_Message:
				var mcs []*response_pb.MessageChainObject
				message := respum.Message
				switch *message.Type {
				case response_pb.MessageType_MessageType_Private:
					mcs = message.Private.MessageChain
				case response_pb.MessageType_MessageType_Group:
					group := message.Group
					if !matchWhiteList(group.GroupId) {
						return
					}
					mcs = group.MessageChain
				}
				text := getText(mcs)
				if text == "" {
					return
				}
				switch message_match(text) {
				case pluginType_RandomAnimal:
					randomanimal(message, text)
				case pluginType_RandomFortune:
					randomfortune(message, text)
				case pluginType_TwoOnOne:
					twoonone(message, text)
				}
			case response_pb.ResponseType_ResponseType_QQEvent:
			}
		}()
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

func randomanimal(message *response_pb.Response_Message, text string) {
	if *message.Type != response_pb.MessageType_MessageType_Group {
		return
	}
	group := message.Group
	action := randomanimal_match(text)
	var resp *randomanimal_pb.BasicResponse
	switch action {
	case randomanimal_GetCat:
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: "猫猫正在赶来的路上",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
		x, err := define.RandomAnimalC.GetCat(define.RandomAnimalCtx, &randomanimal_pb.BasicRequest{
			AutoUpload: true,
		})
		if err != nil {
			logger.Println(err)
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: "猫猫跑到半路跑丢了，可能是出错或者超时，再试一次？",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp = x
	case randomanimal_GetFox:
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: "小狐狸正在赶来的路上",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
		x, err := define.RandomAnimalC.GetFox(define.RandomAnimalCtx, &randomanimal_pb.BasicRequest{
			AutoUpload: true,
		})
		if err != nil {
			logger.Println(err)
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: "小狐狸跑到半路跑丢了，可能是出错或者超时，再试一次？",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp = x
	case randomanimal_GetDog:
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: "狗子正在赶来的路上",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
		x, err := define.RandomAnimalC.GetDog(define.RandomAnimalCtx, &randomanimal_pb.BasicRequest{
			AutoUpload: true,
		})
		if err != nil {
			logger.Println(err)
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: "狗子跑到半路跑丢了，可能是出错或者超时，再试一次？",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp = x
	case randomanimal_GetChicken_CXK:
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: "坤坤正在赶来的路上",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
		x, err := define.RandomAnimalC.GetChiken_CXK(define.RandomAnimalCtx, &randomanimal_pb.BasicRequest{
			AutoUpload: true,
		})
		if err != nil {
			logger.Println(err)
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: "坤坤跑到半路篮球丢了，可能是出错或者超时或者是未添加坤坤图片，再试一次？",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp = x
	case randomanimal_GetDuck:
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: "鸭正在赶来的路上",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
		x, err := define.RandomAnimalC.GetDuck(define.RandomAnimalCtx, &randomanimal_pb.BasicRequest{
			AutoUpload: true,
		})
		if err != nil {
			logger.Println(err)
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: "鸭跑到半路跑丢了，可能是出错或者超时，再试一次？",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp = x
	}
	u := fmt.Sprintf("%v%v", define.ExternalURL, resp.Response.URLPath)
	switch resp.Type {
	case randomanimal_pb.Type_Image:
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Image,
				Image: &request_pb.MessageChain_Image{
					URL: &u,
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case randomanimal_pb.Type_Video:
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Video,
				Video: &request_pb.MessageChain_Video{
					URL: &u,
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	default:
		logger.Println("非预期类型")
		return
	}
}

func sendMessageToGroup(groupid string, mcs []*request_pb.MessageChainObject) error {
	req, err := define.Handler_RequestC.SendGroupMessage(define.HandlerCtx, &request_pb.SendGroupMessageRequest{
		GroupId:      groupid,
		MessageChain: mcs,
	})
	if err != nil {
		return err
	}
	if _, err := define.ConnectorC.Write(define.ConnectorCtx, &connector_pb.WriteRequest{
		Buf: req.Buf,
	}); err != nil {
		return err
	}
	return nil
}

func sendMessageToFriend(friendid string, mcs []*request_pb.MessageChainObject) error {
	req, err := define.Handler_RequestC.SendFriendMessage(define.HandlerCtx, &request_pb.SendFriendMessageRequest{
		FriendId:     friendid,
		MessageChain: mcs,
	})
	if err != nil {
		return err
	}
	if _, err := define.ConnectorC.Write(define.ConnectorCtx, &connector_pb.WriteRequest{
		Buf: req.Buf,
	}); err != nil {
		return err
	}
	return nil
}

func matchWhiteList(groupid string) bool {
	if len(define.Conf.WhiteList) == 0 {
		return true
	}
	for _, v := range define.Conf.WhiteList {
		if groupid == v {
			return true
		}
	}
	return false
}

func getText(mcs []*response_pb.MessageChainObject) string {
	for _, v := range mcs {
		if v.Type == response_pb.MessageChainType_MessageChainType_Text {
			return v.Text.Text
		}
	}
	return ""
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

func randomfortune(message *response_pb.Response_Message, text string) {
	if *message.Type != response_pb.MessageType_MessageType_Group {
		return
	}
	group := message.Group
	action := randomfortune_match(text)
	var resp *randomfortune_pb.BasicResponse
	switch action {
	case randomfortune_GetFortune:
		x, err := define.RandomFortuneC.GetFortune(define.RandomFortuneCtx, &randomfortune_pb.BasicRequest{
			ReturnMethod: randomfortune_pb.BasicRequest_Hash,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		resp = x
	}
	u := fmt.Sprintf("%v%v", define.ExternalURL, resp.Response.URLPath)
	if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
		&request_pb.MessageChainObject{
			Type: request_pb.MessageChainType_MessageChainType_At,
			At: &request_pb.MessageChain_At{
				TargetId: group.SenderId,
			},
		},
		&request_pb.MessageChainObject{
			Type: request_pb.MessageChainType_MessageChainType_Image,
			Image: &request_pb.MessageChain_Image{
				URL: &u,
			},
		},
	}); err != nil {
		logger.Println(err)
		return
	}
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

func twoonone(message *response_pb.Response_Message, text string) {
	if *message.Type != response_pb.MessageType_MessageType_Group {
		return
	}
	group := message.Group
	senderid := group.SenderId
	sendername := ""
	if group.SenderName != nil {
		sendername = *group.SenderName
	}
	action := twoonone_match(text)
	if action == twoonone_JoinORExit {
		if _, ok := twoonone_playerStatus[senderid]; ok {
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 退出斗地主",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			delete(twoonone_playerStatus, senderid)
		} else {
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 进入斗地主",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			twoonone_playerStatus[senderid] = struct{}{}
		}
		return
	}
	if _, ok := twoonone_playerStatus[senderid]; !ok {
		return
	}
	text = twoonone_adjust(action, text)
	switch action {
	case twoonone_CreateAccount:
		if sendername == "" {
			logger.Println("发送者名字为空")
			return
		}
		resp, err := define.TwoOnOneC.CreateAccount(define.TwoOnOneCtx, &twoonone_pb.CreateAccountRequest{
			PlayerId:   senderid,
			PlayerName: sendername,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理错误类型：%v\n", e.String())
			case twoonone_pb.Errors_PlayerAccountExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 开号失败，账号已存在",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: senderid,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 开号成功，已自动领取双倍每日豆子",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_GetAccount:
		resp, err := define.TwoOnOneC.GetAccount(define.TwoOnOneCtx, &twoonone_pb.GetAccountRequest{
			PlayerId: senderid,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理异常：%v\n", e.String())
			case twoonone_pb.Errors_PlayerNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 获取失败，你还未开号",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		ai := resp.Info.AccountInfo
		playCount := ai.WinCount + ai.LoseCount
		winChance := ""
		if playCount != 0 {
			winChance = strconv.FormatFloat((float64(ai.WinCount)/float64(playCount))*100, 'g', 4, 64)
		} else {
			winChance = "0"
		}
		accountInfo := fmt.Sprintf("获取成功，你现在有 %v 个豆子，总共进行了 %v 场游戏，获胜 %v 场，失败 %v 场，胜率 %v%%", ai.Coin, playCount, ai.WinCount, ai.LoseCount, winChance)
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: senderid,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " " + accountInfo,
				},
			},
		}); err != nil {
			panic(err)
		}
	case twoonone_CreateRoom:
		resp, err := define.TwoOnOneC.CreateRoom(define.TwoOnOneCtx, &twoonone_pb.CreateRoomRequest{
			BasicCoin:       200,
			InitialMultiple: 1,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		hash := resp.RoomHash
		id := ""
		for {
			id = randomString(3, OnlyNumber)
			if _, ok := twoonone_rooms[id]; ok {
				continue
			} else {
				twoonone_rooms[id] = &roomSI{
					hash:           hash,
					id:             id,
					landownerCards: []twoonone_pb.Card{},
				}
				break
			}
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: senderid,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 开桌成功，桌id为： " + id,
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_GetDaliyCoin:
		resp, err := define.TwoOnOneC.GetDailyCoin(define.TwoOnOneCtx, &twoonone_pb.GetDailyCoinRequest{
			PlayerId: senderid,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理错误类型：%v\n", e.String())
			case twoonone_pb.Errors_PlayerNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 领取失败，你还未开号",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerAlreadyGetDailyCoin:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 领取失败，你已领取今日豆子",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: senderid,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 领取成功，喜提500豆子",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_RobLandowner_Rob:
		r := twoonone_player2room[senderid]
		if r == nil {
			logger.Println("异常错误")
			return
		}
		resp, err := define.TwoOnOneC.RobLandownerAction(define.TwoOnOneCtx, &twoonone_pb.RobLandownerActionRequest{
			PlayerId: senderid,
			Action:   twoonone_pb.RobLandownerActions_Rob,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理错误类型：%v\n", e.String())
			case twoonone_pb.Errors_PlayerNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 抢地主失败，你还未开号",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_RoomNoRobLandownering:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 抢地主失败，你所在桌不处于抢地主阶段",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoOperatorNow:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 抢地主失败，还未轮到你抢地主",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoExistAnyRoom:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 抢地主失败，你未在任意桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		if resp.MultipleNotice {
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("当前倍率 %v 倍", *resp.Multiple),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
		if resp.IntoSendingCard {
			if err := sendMessageToFriend(group.SenderId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("你的手牌为：%v", cardToCardHuman(resp.NextOperator.TableInfo.Cards)),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf(" %v(%v)当上了地主", resp.NextOperator.AccountInfo.Name, resp.NextOperator.AccountInfo.Id),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("地主牌是%v", cardToCardHuman(r.landownerCards)),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: resp.NextOperator.AccountInfo.Id,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 请出牌",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: resp.NextOperator.AccountInfo.Id,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 轮到你抢地主了",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_RobLandowner_NoRob:
		r := twoonone_player2room[senderid]
		if r == nil {
			logger.Println("异常错误")
			return
		}
		resp, err := define.TwoOnOneC.RobLandownerAction(define.TwoOnOneCtx, &twoonone_pb.RobLandownerActionRequest{
			PlayerId: senderid,
			Action:   twoonone_pb.RobLandownerActions_NoRob,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理错误类型：%v\n", e.String())
			case twoonone_pb.Errors_PlayerNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 不抢地主失败，你还未开号",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_RoomNoRobLandownering:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 不抢地主失败，你所在桌不处于抢地主阶段",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoOperatorNow:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 不抢地主失败，还未轮到你抢地主",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoExistAnyRoom:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 不抢地主失败，你未在任意桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		if resp.IntoSendingCard {
			if err := sendMessageToFriend(group.SenderId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("你的手牌为：%v", cardToCardHuman(resp.NextOperator.TableInfo.Cards)),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf(" %v(%v)当上了地主", resp.NextOperator.AccountInfo.Name, resp.NextOperator.AccountInfo.Id),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("地主牌是%v", cardToCardHuman(r.landownerCards)),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: resp.NextOperator.AccountInfo.Id,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 请出牌",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		} else {
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: resp.NextOperator.AccountInfo.Id,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 轮到你抢地主了",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
	case twoonone_ExitRoom:
		resp, err := define.TwoOnOneC.ExitRoom(define.TwoOnOneCtx, &twoonone_pb.ExitRoomRequest{
			PlayerId: senderid,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理异常：%v\n", e.String())
			case twoonone_pb.Errors_PlayerNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 下桌失败，你还未开号",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_RoomStarted:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 下桌失败，你所在桌已开始游戏",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoExistAnyRoom:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 下桌失败，你未在任意桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
		}
		delete(twoonone_player2room, senderid)
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: senderid,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 下桌成功，当前桌内玩家：\n" + playersToStr(resp.RoomPlayers),
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_StartRoom:
		r := twoonone_player2room[senderid]
		if r == nil {
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 发牌失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.TwoOnOneC.StartRoom(define.TwoOnOneCtx, &twoonone_pb.StartRoomRequest{
			PlayerId: senderid,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理异常：%v\n", e.String())
			case twoonone_pb.Errors_PlayerNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 发牌失败，你还未开号",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_RoomStarted:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 发牌失败，你所在桌已开始游戏",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_RoomPlayerNoFull:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 发牌失败，你所在桌玩家数未满",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoExistAnyRoom:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 发牌失败，你未在任意桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		r.landownerCards = resp.LandownerCards
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: resp.NextOperator.AccountInfo.Id,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 轮到你抢地主了",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_GetRooms:
		roomnames := ""
		i := 0
		for _, v := range twoonone_rooms {
			roomnames += fmt.Sprintf("%v.%v", i+1, v.id)
			if i != len(twoonone_rooms)-1 {
				roomnames += "\n"
			}
			i++
		}
		if len(twoonone_rooms) == 0 {
			roomnames += "空"
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: senderid,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 桌列表：\n" + roomnames,
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_SendCard_NoSend:
		resp, err := define.TwoOnOneC.SendCardAction(define.TwoOnOneCtx, &twoonone_pb.SendCardRequest{
			PlayerId:  senderid,
			Action:    twoonone_pb.SendCardActions_NoSend,
			SendCards: []twoonone_pb.Card{},
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理异常：%v\n", e.String())
			case twoonone_pb.Errors_PlayerNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 不出牌失败，你还未开号",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_RoomNoSendingCards:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 不出牌失败，你所在桌还未进入出牌阶段",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoExistAnyRoom:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 不出牌失败，你未在任意桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoOperatorNow:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 不出牌失败，还未轮到你出牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerIsOnlySendCarder:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 不出牌失败，你是唯一可以出牌的玩家",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: resp.NextOperator.AccountInfo.Id,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 轮到你出牌了",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_SendCard_Send:
		cardStr := strings.ReplaceAll(text, "10", "X")
		cardStr = strings.ReplaceAll(cardStr, "小王", "Y")
		cardStr = strings.ReplaceAll(cardStr, "大王", "Z")
		cardStrS := strings.Split(cardStr, "")
		card := []twoonone_pb.Card{}
		for _, v := range cardStrS {
			card = append(card, cardStr2Card[v])
		}
		resp, err := define.TwoOnOneC.SendCardAction(define.TwoOnOneCtx, &twoonone_pb.SendCardRequest{
			PlayerId:  senderid,
			Action:    twoonone_pb.SendCardActions_Send,
			SendCards: card,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理错误类型：%v\n", e.String())
			case twoonone_pb.Errors_SendCardUnknown:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 出牌失败，无法匹配你的牌型",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerCardNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 出牌失败，你的手牌不足",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_SendCardSizeLELastCard:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 出牌失败，你出的牌比上一张牌小",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_SendCardTypeNELastCard:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 出牌失败，你的牌型不为上一张牌型",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_SendCardContinousNELastCard:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 出牌失败，你的牌连续数与上一张牌不同",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 出牌失败，你还未开号",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_RoomNoSendingCards:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 出牌失败，你所在桌还未进入出牌阶段",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoExistAnyRoom:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 出牌失败，你未在任意桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoOperatorNow:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 出牌失败，还未轮到你出牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		if resp.GameFinish {
			for {
				rhash := ""
				for k, v := range twoonone_player2room {
					if k == senderid {
						rhash = v.hash
						// 删除桌
						delete(twoonone_rooms, v.id)
						break
					}
				}
				for k, v := range twoonone_player2room {
					// 删除在桌内的玩家
					if v.hash == rhash {
						delete(twoonone_player2room, k)
					}
				}
				break
			}
			if resp.GameFinishE.Spring {
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: fmt.Sprintf("春天！当前倍率 %v 倍", resp.GameFinishE.Multiple),
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			winis := ""
			switch resp.GameFinishE.Winner {
			case twoonone_pb.Role_Landowner:
				winis = "地主"
			case twoonone_pb.Role_Farmer:
				winis = "农民"
			}
			lo := resp.GameFinishE.Landowner
			f1 := resp.GameFinishE.Farmer1
			f2 := resp.GameFinishE.Farmer2
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("游戏结束，%v获胜\n%v(%v)剩余牌：%v\n%v(%v)剩余牌：%v\n%v(%v)剩余牌：%v\n", winis, lo.AccountInfo.Name, lo.AccountInfo.Id, cardToCardHuman(lo.TableInfo.Cards), f1.AccountInfo.Name, f1.AccountInfo.Id, cardToCardHuman(f1.TableInfo.Cards), f2.AccountInfo.Name, f2.AccountInfo.Id, cardToCardHuman(f2.TableInfo.Cards)),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		if resp.SenderCardNumberNotice {
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("%v(%v)就剩 %v 张牌了", sendername, senderid, *resp.SenderCardNumber),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
		if resp.SenderCardTypeNotice {
			ctStr := ""
			ct := resp.SenderCardTypeNoticeE.SenderCardType
			multiple := resp.SenderCardTypeNoticeE.Multiple
			switch ct {
			case twoonone_pb.CardType_KingBomb:
				ctStr = "王炸！"
			case twoonone_pb.CardType_Bomb:
				ctStr = "炸弹！"
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("%v当前倍率 %v 倍", ctStr, multiple),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: fmt.Sprintf(" %v(%v)出了%v", sendername, senderid, cardToCardHuman(card)),
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
		if err := sendMessageToFriend(group.SenderId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: fmt.Sprintf("你的手牌为：%v", cardToCardHuman(resp.NextOperator.TableInfo.Cards)),
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: resp.NextOperator.AccountInfo.Id,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 轮到你出牌了",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_JoinRoom:
		echo := randomString(6, OnlyNumber)
		req, err := define.Handler_RequestC.GetFriendList(define.HandlerCtx, &request_pb.BasicRequest{
			Echo: &echo,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		go func() {
			if _, err := define.ConnectorC.Write(define.ConnectorCtx, &connector_pb.WriteRequest{
				Buf: req.Buf,
			}); err != nil {
				logger.Println(err)
				return
			}
		}()
		s, err := define.ConnectorC.Read(define.ConnectorCtx, &connector_pb.Empty{})
		if err != nil {
			logger.Println(err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		go func() {
			defer cancel()
			ctx = context.WithValue(ctx, myBool{}, false)
			for {
				recv, err := s.Recv()
				if err != nil {
					logger.Println(err)
					return
				}
				resp, err := define.Handler_ResponseC.Unmarshal(define.HandlerCtx, &response_pb.UnmarshalRequest{
					Buf:          recv.Buf,
					Type:         response_pb.ResponseType_ResponseType_CmdEvent.Enum(),
					CmdEventType: response_pb.CmdEventType_CmdEventType_GetFriendList.Enum(),
				})
				if err != nil {
					continue
				}
				if resp.CmdEvent.Echo != echo {
					continue
				}
				gfl := resp.CmdEvent.GetFriendList
				if !gfl.OK {
					logger.Printf("获取好友列表失败, retcode: %v\n", *gfl.Retcode)
					return
				}
				for _, v := range gfl.Friends {
					if v.UserId == senderid {
						ctx = context.WithValue(ctx, myBool{}, true)
						return
					}
				}
			}
		}()
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			return
		}
		ok := ctx.Value(myBool{}).(bool)
		if !ok {
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 上桌失败，请先添加机器人为好友",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		id := text
		r := twoonone_rooms[id]
		if r == nil {
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 上桌失败，该桌不存在",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.TwoOnOneC.JoinRoom(define.TwoOnOneCtx, &twoonone_pb.JoinRoomRequest{
			RoomHash: r.hash,
			PlayerId: senderid,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理错误类型：%v\n", e.String())
			case twoonone_pb.Errors_RoomFull:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 上桌失败，房间已满",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_RoomExistPlayer:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 上桌失败，你已在桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_PlayerNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 上桌失败，你还未开号",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case twoonone_pb.Errors_RoomNoExist:
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 上桌失败，未找到该桌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		twoonone_player2room[senderid] = r
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: senderid,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 上桌成功，当前桌内玩家：\n" + playersToStr(resp.RoomPlayers),
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case twoonone_GetRoom:
		if text == "" {
			r := twoonone_player2room[senderid]
			if r == nil {
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 获取桌信息失败，你未在任意桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
				return
			}
			resp, err := define.TwoOnOneC.GetRoom(define.TwoOnOneCtx, &twoonone_pb.GetRoomRequest{
				RoomHash: &r.hash,
			})
			if err != nil {
				logger.Println(err)
				return
			}
			ri := resp.Info
			stageStr := ""
			switch ri.Stage {
			case twoonone_pb.RoomStage_WaitingStart:
				stageStr = "等待开始"
			case twoonone_pb.RoomStage_RobLandownering:
				stageStr = "抢地主中"
			case twoonone_pb.RoomStage_SendingCards:
				stageStr = "出牌中"
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf(` 获取桌信息成功，你所在桌信息如下：
						id：%v
						哈希：%v
						底分：%v
						倍率：%v
						游戏状态：%v
						玩家列表：%v`, r.id, ri.Hash, ri.BasicCoin, ri.Multiple, stageStr, "\n"+playersToStr(ri.Players)),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		} else {
			r := twoonone_rooms[text]
			if r == nil {
				if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_At,
						At: &request_pb.MessageChain_At{
							TargetId: senderid,
						},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: " 获取桌信息失败，未找到该桌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
				return
			}
			resp, err := define.TwoOnOneC.GetRoom(define.TwoOnOneCtx, &twoonone_pb.GetRoomRequest{
				RoomHash: &r.hash,
			})
			if err != nil {
				logger.Println(err)
				return
			}
			ri := resp.Info
			stageStr := ""
			switch ri.Stage {
			case twoonone_pb.RoomStage_WaitingStart:
				stageStr = "等待开始"
			case twoonone_pb.RoomStage_RobLandownering:
				stageStr = "抢地主中"
			case twoonone_pb.RoomStage_SendingCards:
				stageStr = "出牌中"
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf(` 获取桌信息成功， %v 桌信息如下：
						id：%v
						哈希：%v
						底分：%v
						倍率：%v
						游戏状态：%v
						玩家列表：%v`, r.id, r.id, ri.Hash, ri.BasicCoin, ri.Multiple, stageStr, "\n"+playersToStr(ri.Players)),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
	}
}

func playersToStr(x []*twoonone_pb.PlayerInfo) string {
	str := ""
	for i, v := range x {
		str += fmt.Sprintf("%v(%v)", v.AccountInfo.Name, v.AccountInfo.Id)
		if i != len(x)-1 {
			str += "\n"
		}
	}
	if str == "" {
		str = "空"
	}
	return str
}

type myBool struct{}

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
	twoonone_JoinORExit
	twoonone_GetRoom
)

func twoonone_adjust(action twoononeAction, text string) string {
	switch action {
	case twoonone_JoinRoom:
		str := strings.TrimPrefix(text, "上桌")
		return strings.TrimSpace(str)
	case twoonone_SendCard_Send:
		str := strings.TrimPrefix(text, "!")
		str = strings.TrimPrefix(str, "！")
		str = strings.TrimSpace(str)
		return strings.ToUpper(str)
	case twoonone_GetRoom:
		str := strings.TrimPrefix(text, "桌信息")
		return strings.TrimSpace(str)
	default:
		return text
	}
}

func twoonone_match(text string) twoononeAction {
	switch text {
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
	case "#斗地主":
		return twoonone_JoinORExit
	}
	if ok, _ := regexp.MatchString(`\A(!|！)([3456789jqkaJQKA2]|10|大王|小王)+`, text); ok {
		return twoonone_SendCard_Send
	}
	if ok, _ := regexp.MatchString(`\A上桌([(0-9)]| ){3,}`, text); ok {
		return twoonone_JoinRoom
	}
	if ok, _ := regexp.MatchString(`\A桌信息.*`, text); ok {
		return twoonone_GetRoom
	}
	return twoonone_Unknown
}

type dictionary int

const (
	OnlyNumber dictionary = iota
)

func randomString(length int, dic dictionary) string {
	d := ""
	switch dic {
	case OnlyNumber:
		d = "0123456789"
	}
	var builder strings.Builder
	for n := 0; n < length; n++ {
		builder.Write([]byte{d[rand.Intn(len(d))]})
	}
	return builder.String()
}

func cardToCardHuman(x []twoonone_pb.Card) string {
	// 升序
	sort.Slice(x, func(i, j int) bool {
		return x[i] < x[j]
	})
	cardH := ""
	for _, v := range x {
		cardH += fmt.Sprintf("[%v]", card2cardStr[v])
	}
	return cardH
}

func init() {
	cardStr2Card = map[string]twoonone_pb.Card{
		"3": 0,
		"4": 1,
		"5": 2,
		"6": 3,
		"7": 4,
		"8": 5,
		"9": 6,
		"X": 7,
		"J": 8,
		"Q": 9,
		"K": 10,
		"A": 11,
		"2": 12,
		"Y": 13,
		"Z": 14,
	}
	card2cardStr = map[twoonone_pb.Card]string{
		0:  "3",
		1:  "4",
		2:  "5",
		3:  "6",
		4:  "7",
		5:  "8",
		6:  "9",
		7:  "10",
		8:  "J",
		9:  "Q",
		10: "K",
		11: "A",
		12: "2",
		13: "小王",
		14: "大王",
	}
}

var (
	cardStr2Card map[string]twoonone_pb.Card
	card2cardStr map[twoonone_pb.Card]string
)
