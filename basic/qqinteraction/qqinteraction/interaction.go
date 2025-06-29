package qqinteraction

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/draw"

	"github.com/nanachi-sh/susubot-code/basic/qqinteraction/define"
	"github.com/nanachi-sh/susubot-code/basic/qqinteraction/log"
	connector_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/connector"
	request_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/handler/request"
	response_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/handler/response"
	randomanimal_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/randomanimal"
	randomfortune_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/randomfortune"
	twoonone_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/twoonone"
	uno_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/uno"
	"github.com/nfnt/resize"
	"google.golang.org/grpc/metadata"
)

var logger = log.Get()

var (
	twoonone_rooms        = make(map[string]*roomSI) //id To room
	twoonone_player2room  = make(map[string]*roomSI) //
	twoonone_playerStatus = make(map[string]struct{})

	uno_rooms        = make(map[string]*roomUNO)
	uno_player2room  = make(map[string]*roomUNO)
	uno_playerinfo   = make(map[string]*playerinfoUNO)
	uno_playerStatus = make(map[string]struct{})
	uno_privilegeCtx context.Context
)

const (
	uno_imageDir = "/config/unoCardImages"
)

type roomSI struct {
	hash           string
	id             string
	landownerCards []twoonone_pb.Card
}

type roomUNO struct {
	groupid string
	hash    string
	id      string
}

type playerinfoUNO struct {
	roomHash  string
	hash      string
	playerCtx context.Context
}

func Start() {
	stream, err := define.ConnectorC.Read(define.ConnectorCtx, &connector_pb.Empty{})
	if err != nil {
		logger.Fatalln(err)
	}
	if rs, err := define.TwoOnOneC.GetRooms(define.TwoOnOneCtx, &twoonone_pb.Empty{}); err != nil {
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
	if rs, err := define.UnoC.GetRooms(uno_privilegeCtx, &uno_pb.Empty{}); err != nil {
		logger.Println(err)
	} else {
		for _, v := range rs.Rooms {
			var r *roomUNO
			for {
				id := randomString(3, OnlyNumber)
				if _, ok := uno_rooms[id]; ok {
					continue
				} else {
					r = &roomUNO{
						hash: v.Hash,
						id:   id,
					}
					uno_rooms[id] = r
					break
				}
			}
			for _, v2 := range v.Players {
				uno_player2room[v2.Id] = r
			}
		}
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			logger.Fatalln(err)
		}
		fmt.Println(string(resp.Buf))
		go func() {
			respum, err := define.Handler_ResponseC.Unmarshal(define.HandlerCtx, &response_pb.UnmarshalRequest{
				Buf:            resp.Buf,
				ExtraInfo:      true,
				IgnoreCmdEvent: true,
			})
			if err != nil {
				logger.Println(err)
				return
			}
			respu := respum.GetResponse()
			if respu == nil {
				return
			}
			switch respu.Type {
			case response_pb.ResponseType_ResponseType_CmdEvent:
				return
			case response_pb.ResponseType_ResponseType_Message:
				var mcs []*response_pb.MessageChainObject
				message := respu.GetMessage()
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
				matchs, ok := message_match(text)
				if ok {
					for _, v := range matchs {
						switch v {
						case pluginType_RandomAnimal:
							go randomanimal(message, text)
						case pluginType_RandomFortune:
							go randomfortune(message, text)
						case pluginType_TwoOnOne:
							go twoonone(message, text)
						case pluginType_Uno:
							go uno(message, text)
						case pluginType_TEST:
							go test(message, text)
						}
					}
				}
			case response_pb.ResponseType_ResponseType_QQEvent:
			}
		}()
	}
}

func message_match(text string) ([]pluginType, bool) {
	ret := []pluginType{}
	// if randomanimal_match(text) != randomanimal_Unknown {
	// 	ret = append(ret, pluginType_RandomAnimal)
	// }
	// if randomfortune_match(text) != randomfortune_Unknown {
	// 	ret = append(ret, pluginType_RandomFortune)
	// }
	// if twoonone_match(text) != twoonone_Unknown {
	// 	ret = append(ret, pluginType_TwoOnOne)
	// }
	// if uno_match(text) != uno_Unknown {
	// 	ret = append(ret, pluginType_Uno)
	// }
	ret = append(ret, pluginType_TEST)
	if len(ret) == 0 {
		return nil, false
	}
	return ret, true
}

type pluginType int

const (
	pluginType_Unknown pluginType = iota
	pluginType_RandomAnimal
	pluginType_RandomFortune
	pluginType_TwoOnOne
	pluginType_Uno
	pluginType_TEST
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
			MemberId:     &group.SenderId,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		resp = x
	}
	if resp.AlreadyGetFortune {
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_At,
				At: &request_pb.MessageChain_At{
					TargetId: group.SenderId,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: " 你今天已经求过签了",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
		return
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
			return
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
			return
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
						Text: fmt.Sprintf("你的手牌为：%v", twoonone_cardToCardHuman(resp.NextOperator.TableInfo.Cards)),
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
						Text: fmt.Sprintf("%v(%v)当上了地主", resp.NextOperator.AccountInfo.Name, resp.NextOperator.AccountInfo.Id),
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
						Text: fmt.Sprintf("地主牌是%v", twoonone_cardToCardHuman(r.landownerCards)),
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
						Text: fmt.Sprintf("你的手牌为：%v", twoonone_cardToCardHuman(resp.NextOperator.TableInfo.Cards)),
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
						Text: fmt.Sprintf("%v(%v)当上了地主", resp.NextOperator.AccountInfo.Name, resp.NextOperator.AccountInfo.Id),
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
						Text: fmt.Sprintf("地主牌是%v", twoonone_cardToCardHuman(r.landownerCards)),
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
			return
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
					Text: " 下桌成功，当前桌内玩家：\n" + twoonone_playersToStr(resp.RoomPlayers),
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
						Text: fmt.Sprintf("游戏结束，%v获胜\n%v(%v)剩余牌：%v\n%v(%v)剩余牌：%v\n%v(%v)剩余牌：%v\n", winis, lo.AccountInfo.Name, lo.AccountInfo.Id, twoonone_cardToCardHuman(lo.TableInfo.Cards), f1.AccountInfo.Name, f1.AccountInfo.Id, twoonone_cardToCardHuman(f1.TableInfo.Cards), f2.AccountInfo.Name, f2.AccountInfo.Id, twoonone_cardToCardHuman(f2.TableInfo.Cards)),
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
					Text: fmt.Sprintf(" %v(%v)出了%v", sendername, senderid, twoonone_cardToCardHuman(card)),
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
					Text: fmt.Sprintf("你的手牌为：%v", twoonone_cardToCardHuman(resp.NextOperator.TableInfo.Cards)),
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
				ce := resp.GetResponse().GetCmdEvent()
				if ce.Echo != echo {
					continue
				}
				gfl := ce.GetFriendList
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
					Text: " 上桌成功，当前桌内玩家：\n" + twoonone_playersToStr(resp.RoomPlayers),
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
						玩家列表：%v`, r.id, ri.Hash, ri.BasicCoin, ri.Multiple, stageStr, "\n"+twoonone_playersToStr(ri.Players)),
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
						玩家列表：%v`, r.id, r.id, ri.Hash, ri.BasicCoin, ri.Multiple, stageStr, "\n"+twoonone_playersToStr(ri.Players)),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
	}
}

func twoonone_playersToStr(x []*twoonone_pb.PlayerInfo) string {
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

func twoonone_cardToCardHuman(x []twoonone_pb.Card) string {
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
	cs, err := http.ParseCookie(fmt.Sprintf("account_hash=%v", define.PrivilegeUserHash))
	if err != nil {
		logger.Fatalln(err)
	}
	md, ok := metadata.FromOutgoingContext(define.UnoCtx)
	if ok {
		md.Append("Cookie", fmt.Sprintf("%v=%v", cs[0].Name, cs[0].Value))
		uno_privilegeCtx = metadata.NewOutgoingContext(define.UnoCtx, md)
	} else {
		md := metadata.New(map[string]string{})
		md.Append("Cookie", fmt.Sprintf("%v=%v", cs[0].Name, cs[0].Value))
		uno_privilegeCtx = metadata.NewOutgoingContext(define.UnoCtx, md)
	}
}

var (
	cardStr2Card map[string]twoonone_pb.Card
	card2cardStr map[twoonone_pb.Card]string
)

func uno_match(text string) unoAction {
	switch text {
	case "开桌":
		return uno_CreateRoom
	case "下桌":
		return uno_ExitRoom
	case "桌列表":
		return uno_GetRooms
	case "不要", "要不起", "不出":
		return uno_SendCard_NoSend
	case "开始游戏":
		return uno_StartRoom
	case "#UNO", "#乌诺", "#uno", "#Uno":
		return uno_JoinORExit
	case "挑战":
		return uno_Challenge
	case "摸牌", "摸卡", "抽牌", "抽卡":
		return uno_DrawCard
	case "最后一张", "上一张", "上一张牌":
		return uno_GetLastCard
	}
	switch strings.ToUpper(text) {
	case "UNO!", "UNO！":
		return uno_CallUNO
	}
	if ok, _ := regexp.MatchString(`(?i)^(?:!|！)+(([RGYB][0-9])?|([RGYB](dt|draw two|wild draw four|wdf|reverse|rev|re|wild|skip|\+2|\+4))?){1,1}$`, text); ok {
		return uno_SendCard_Send
	}
	if ok, _ := regexp.MatchString(`\A上桌([(0-9)]| ){3,}`, text); ok {
		return uno_JoinRoom
	}
	if ok, _ := regexp.MatchString(`\A桌信息([0-9 ]{3,})?`, text); ok {
		return uno_GetRoom
	}
	if ok, _ := regexp.MatchString(`(?i)没喊uno`, text); ok {
		return uno_IndicateUNO
	}
	return uno_Unknown
}

func uno_adjust(action unoAction, text string) string {
	switch action {
	case uno_JoinRoom:
		str := strings.TrimPrefix(text, "上桌")
		return strings.TrimSpace(str)
	case uno_SendCard_Send:
		str := strings.TrimPrefix(text, "!")
		str = strings.TrimPrefix(str, "！")
		str = strings.TrimSpace(str)
		return strings.ToUpper(str)
	case uno_GetRoom:
		str := strings.TrimPrefix(text, "桌信息")
		return strings.TrimSpace(str)
	default:
		return text
	}
}

func (unoR *roomUNO) listenEvent() {
	stream, err := define.UnoC.RoomEvent(uno_privilegeCtx, &uno_pb.RoomEventRequest{
		RoomHash:   unoR.hash,
		PlayerHash: define.PrivilegeUserHash,
	})
	if err != nil {
		logger.Println(err)
		return
	}
	for close := false; close; {
		resp, err := stream.Recv()
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.DrawCard_IntoSendCard != nil {
			e := resp.DrawCard_IntoSendCard
			if err := sendMessageToGroup(unoR.groupid, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("%v(%v)当上了庄家", e.Banker.Name, e.Banker.Id),
					},
				},
			}); err != nil {
				logger.Println(err)
				continue
			}
			for _, v := range e.Players {
				resp, err := define.UnoC.GetPlayer(uno_privilegeCtx, &uno_pb.GetPlayerRequest{
					RoomHash: unoR.hash,
					PlayerId: v.Id,
				})
				if err != nil {
					logger.Println(err)
					continue
				}
				if resp.Err != nil {
					switch e := *resp.Err; e {
					default:
						logger.Printf("未处理错误类型：%v\n", e.String())
					}
					continue
				}
				p := resp.Extra
				cs := []uno_pb.Card{}
				for _, v := range p.PlayerRoomInfo.Cards {
					cs = append(cs, *v)
				}
				img, err := uno_generateCardsImage(cs, uno_defaultColumn)
				if err != nil {
					logger.Println(err)
					return
				}
				buf, err := image2Buf(img)
				if err != nil {
					logger.Println(err)
					return
				}
				if err := sendMessageToFriend(p.PlayerAccountInfo.Id, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Image,
						Image: &request_pb.MessageChain_Image{
							Buf: buf,
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			img, err := uno_getCardImage(*e.LeadCard, nil)
			if err != nil {
				logger.Println(err)
				return
			}
			buf, err := image2Buf(img)
			if err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToGroup(unoR.groupid, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Image,
					Image: &request_pb.MessageChain_Image{
						Buf: buf,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: "引牌为",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToGroup(unoR.groupid, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: e.Banker.Id,
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
			return
		}
		if resp.DrawCard_Skipped != nil {
			e := resp.DrawCard_Skipped
			if err := sendMessageToGroup(unoR.groupid, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: e.NextOperator.Id,
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
		}
		if resp.GameFinish != nil {
			close = true
			for {
				rhash := ""
				delete(uno_rooms, unoR.hash)
				for k, v := range uno_player2room {
					// 删除在桌内的玩家
					if v.hash == rhash {
						delete(uno_player2room, k)
					}
				}
				break
			}
			e := resp.GameFinish
			winner := e.Winner
			if err := sendMessageToGroup(unoR.groupid, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: fmt.Sprintf("游戏结束，%v(%v)获胜", winner.PlayerAccountInfo.Name, winner.PlayerAccountInfo.Id),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			for _, v := range e.Players {
				if v.PlayerAccountInfo.Id == winner.PlayerAccountInfo.Id {
					continue
				}
				cs := []uno_pb.Card{}
				for _, v := range v.PlayerRoomInfo.Cards {
					cs = append(cs, *v)
				}
				img, err := uno_generateCardsImage(cs, -1)
				if err != nil {
					logger.Println(err)
					return
				}
				buf, err := image2Buf(img)
				if err != nil {
					logger.Println(err)
					return
				}
				if err := sendMessageToGroup(unoR.groupid, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type:  request_pb.MessageChainType_MessageChainType_Image,
						Image: &request_pb.MessageChain_Image{Buf: buf},
					},
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Text,
						Text: &request_pb.MessageChain_Text{
							Text: fmt.Sprintf("%v(%v)剩余牌", v.PlayerAccountInfo.Name, v.PlayerAccountInfo.Id),
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
	}
}

func GIFs2Str(g GIFs) string {
	switch g {
	case cxk:
		return "蔡徐坤"
	case rua:
		return "RUARUARUA！"
	}
	return ""
}

func test(message *response_pb.Response_Message, text string) {
	if message.Group == nil {
		return
	}
	var is GIFs
	fmt.Println(text)
	if regexp.MustCompile(`^(@.+)?蔡徐坤$`).MatchString(text) {
		is = cxk
	} else if regexp.MustCompile(`^(@.+)?rua$`).MatchString(text) {
		is = rua
	} else {
		return
	}
	genId := message.Group.SenderId
	for _, v := range message.Group.MessageChain {
		if v.At != nil {
			genId = v.At.TargetId
		}
	}
	fmt.Printf("接收，生成目标 %s，目标ID %s\n", GIFs2Str(is), genId)
	t := time.Now()
	buf, err := genGIF(genId, is)
	if err != nil {
		panic(err)
	}
	fmt.Printf("生成完毕，耗时: %s\n", time.Since(t).String())
	resph, err := define.Handler_RequestC.SendGroupMessage(define.ConnectorCtx, &request_pb.SendGroupMessageRequest{
		GroupId: message.Group.GroupId,
		MessageChain: []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Image,
				Image: &request_pb.MessageChain_Image{
					Buf: buf,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	req := resph.GetBuf()
	_, err = define.ConnectorC.Write(define.ConnectorCtx, &connector_pb.WriteRequest{
		Buf: req,
	})
	if err != nil {
		panic(err)
	}
}

type GIFs int

const (
	cxk GIFs = iota
	rua
)

func genGIF(id string, is GIFs) ([]byte, error) {
	dirN := strconv.FormatInt(rand.Int63(), 10)
	if err := os.Mkdir(dirN, 0755); err != nil {
		return nil, err
	}
	avatar, err := loadAvatar(id)
	if err != nil {
		return nil, err
	}
	var info *info
	switch is {
	case cxk:
		targets, i, err := loadTarget("cxk")
		if err != nil {
			return nil, err
		}
		info = i
		wg := new(sync.WaitGroup)
		for _, target := range targets {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w := target.rect.Max.X - target.rect.Min.X
				h := target.rect.Max.Y - target.rect.Min.Y
				avatarCircle := image.NewRGBA(avatar.Bounds())
				draw.DrawMask(avatarCircle, avatarCircle.Rect, avatar, image.Point{}, &circle{image.Pt(avatar.Bounds().Max.X, avatar.Bounds().Max.Y), avatar.Bounds().Dx() / 2}, image.Point{avatar.Bounds().Max.X / 2, avatar.Bounds().Max.Y / 2}, draw.Over)
				avatar = resize.Resize(uint(w), uint(h), avatarCircle, resize.Bilinear)
				bg := image.NewRGBA(target.image.Bounds())
				draw.Copy(bg, image.Point{}, target.image, target.image.Bounds(), draw.Src, nil)
				draw.Draw(bg, target.rect, avatar, image.Point{}, draw.Over)
				f, err := os.Create(fmt.Sprintf("%s/%d.png", dirN, target.frame))
				if err != nil {
					panic(err)
				}
				if err := png.Encode(f, bg); err != nil {
					panic(err)
				}
			}()
		}
		wg.Wait()
	case rua:
		targets, i, err := loadTarget("rua")
		if err != nil {
			return nil, err
		}
		info = i
		wg := new(sync.WaitGroup)
		for _, target := range targets {
			wg.Add(1)
			go func() {
				defer wg.Done()
				bg := image.NewRGBA(target.image.Bounds())
				draw.Copy(bg, image.Point{}, target.image, bg.Rect, draw.Src, nil)
				w := target.rect.Max.X - target.rect.Min.X
				h := target.rect.Max.Y - target.rect.Min.Y
				avatarCircle := image.NewRGBA(avatar.Bounds())
				draw.DrawMask(avatarCircle, avatarCircle.Rect, avatar, image.Point{}, &circle{image.Pt(avatar.Bounds().Max.X, avatar.Bounds().Max.Y), avatar.Bounds().Dx() / 2}, image.Point{avatar.Bounds().Max.X / 2, avatar.Bounds().Max.Y / 2}, draw.Over)
				avatar = resize.Resize(uint(w), uint(h), avatarCircle, resize.Bilinear)
				draw.DrawMask(bg, target.rect, avatar, image.Point{}, &NoExistMask{bg}, image.Point{target.rect.Min.X, target.rect.Min.Y}, draw.Over)
				f, err := os.Create(fmt.Sprintf("%s/%d.png", dirN, target.frame))
				if err != nil {
					panic(err)
				}
				if err := png.Encode(f, bg); err != nil {
					panic(err)
				}
			}()
		}
		wg.Wait()
	default:
		return nil, nil
	}
	palettegen := exec.Command("ffmpeg", "-loglevel", "quiet", "-i", dirN+`/%d.png`, "-vf", "palettegen", "-y", fmt.Sprintf("%s/palette.png", dirN))
	gifgen := exec.Command("ffmpeg", "-loglevel", "quiet", "-r", fmt.Sprintf("%d", info.delay), "-i", dirN+`/%d.png`, "-i", fmt.Sprintf("%s/palette.png", dirN), "-lavfi", "paletteuse", "-y", fmt.Sprintf("%s/output.gif", dirN))
	if err := palettegen.Run(); err != nil {
		return nil, err
	}
	if err := gifgen.Run(); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(fmt.Sprintf("%s/output.gif", dirN))
	if err != nil {
		return nil, err
	}
	if err := os.RemoveAll(dirN); err != nil {
		return nil, err
	}
	return buf, nil
}

type NoExistMask struct {
	img image.Image
}

func (e *NoExistMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (e *NoExistMask) Bounds() image.Rectangle {
	return e.img.Bounds()
}

func (e *NoExistMask) At(x, y int) color.Color {
	c := e.img.At(x, y)
	_, _, _, a := c.RGBA()
	if a == 0 {
		return color.RGBA{0, 0, 0, 255}
	} else {
		return color.Alpha{0}
	}
}

type target struct {
	frame int
	image image.Image
	rect  image.Rectangle
}

type info struct {
	delay int
}

func loadTarget(dir string) ([]*target, *info, error) {
	m := map[int]*target{}
	info := &info{}
	entrys, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	for _, entry := range entrys {
		path := fmt.Sprintf("%s/%s", dir, entry.Name())
		if entry.Name() == "info" {
			buf, err := os.ReadFile(path)
			if err != nil {
				return nil, nil, err
			}
			spl := strings.Split(string(buf), "\n")
			delayS := spl[0]
			delay, err := strconv.Atoi(delayS)
			if err != nil {
				return nil, nil, err
			}
			info.delay = delay
			continue
		}
		spl := strings.Split(entry.Name(), ".")
		frame := spl[0]
		format := spl[1]
		frameI, err := strconv.ParseInt(frame, 10, 32)
		if err != nil {
			return nil, nil, err
		}
		if _, ok := m[int(frameI)]; !ok {
			m[int(frameI)] = &target{frame: int(frameI)}
		}
		d := m[int(frameI)]
		buf, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, err
		}
		switch format {
		case "png":
			img, err := png.Decode(bytes.NewBuffer(buf))
			if err != nil {
				return nil, nil, err
			}
			d.image = img
		case "jpg":
			img, err := jpeg.Decode(bytes.NewBuffer(buf))
			if err != nil {
				return nil, nil, err
			}
			d.image = img
		case "json":
			j := make(map[string]any)
			if err := json.Unmarshal(buf, &j); err != nil {
				return nil, nil, err
			}
			p := j["shapes"].([]any)[0].(map[string]any)["points"].([]any)
			xy1 := p[1].([]any)
			xy2 := p[0].([]any)
			x1, y1 := xy1[0].(float64), xy1[1].(float64)
			x2, y2 := xy2[0].(float64), xy2[1].(float64)
			d.rect = image.Rect(
				int(x1),
				int(y1),
				int(x2),
				int(y2),
			)
		}
	}
	targets := []*target{}
	for frame := 1; ; frame++ {
		d, ok := m[frame]
		if !ok {
			break
		}
		targets = append(targets, d)
	}
	return targets, info, nil
}

// func loadTarget(dir string) ([]*target, error) {
// 	m := map[int]*target{}
// 	entrys, err := os.ReadDir(dir)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, entry := range entrys {
// 		path := fmt.Sprintf("%s/%s", dir, entry.Name())
// 		spl := strings.Split(entry.Name(), ".")
// 		format := spl[1]
// 		spl = strings.Split(spl[0], "_")
// 		frame := spl[1]
// 		frameI, err := strconv.ParseInt(frame, 10, 32)
// 		if err != nil {
// 			panic(err)
// 		}
// 		if _, ok := m[int(frameI)]; !ok {
// 			m[int(frameI)] = &target{frame: int(frameI)}
// 		}
// 		d := m[int(frameI)]
// 		buf, err := os.ReadFile(path)
// 		if err != nil {
// 			panic(err)
// 		}
// 		if format == "jpg" {
// 			img, err := jpeg.Decode(bytes.NewBuffer(buf))
// 			if err != nil {
// 				panic(err)
// 			}
// 			d.image = img
// 		} else if format == "json" {
// 			j := make(map[string]any)
// 			if err := json.Unmarshal(buf, &j); err != nil {
// 				panic(err)
// 			}
// 			p := j["shapes"].([]any)[0].(map[string]any)["points"].([]any)
// 			xy1 := p[1].([]any)
// 			xy2 := p[0].([]any)
// 			x1, y1 := xy1[0].(float64), xy1[1].(float64)
// 			x2, y2 := xy2[0].(float64), xy2[1].(float64)
// 			d.rect = image.Rect(
// 				int(x1),
// 				int(y1),
// 				int(x2),
// 				int(y2),
// 			)
// 		}
// 	}
// 	targets := []*target{}
// 	for frame := 1; ; frame++ {
// 		d, ok := m[frame]
// 		if !ok {
// 			break
// 		}
// 		targets = append(targets, d)
// 	}
// 	return targets, nil
// }

// func loadAvatar(path string) (image.Image, error) {
// 	buf, err := os.ReadFile(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return jpeg.Decode(bytes.NewBuffer(buf))
// }

func loadAvatar(id string) (image.Image, error) {
	resp, err := http.Get(fmt.Sprintf("https://q2.qlogo.cn/headimg_dl?dst_uin=%s&spec=5", id))
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(resp.Body)
}

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.RGBA{0, 0, 0, 255}
	}
	return color.Alpha{0}
}

func uno(message *response_pb.Response_Message, text string) {
	if *message.Type != response_pb.MessageType_MessageType_Group {
		return
	}
	group := message.Group
	senderid := group.SenderId
	sendername := ""
	if group.SenderName != nil {
		sendername = *group.SenderName
	}
	action := uno_match(text)
	if action == uno_JoinORExit {
		if _, ok := uno_playerStatus[senderid]; ok {
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
						Text: " 退出UNO",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			delete(uno_playerStatus, senderid)
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
						Text: " 进入UNO",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			uno_playerStatus[senderid] = struct{}{}
		}
		return
	}
	if _, ok := uno_playerStatus[senderid]; !ok {
		return
	}
	text = uno_adjust(action, text)
	switch action {
	case uno_CreateRoom:
		resp, err := define.UnoC.CreateRoom(uno_privilegeCtx, &uno_pb.Empty{})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理异常：%v\n", e.String())
			case uno_pb.Errors_NoFoundAccountHash, uno_pb.Errors_NoValidAccountHash:
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
							Text: " 开桌失败，特权用户哈希设置有误，请联系苏苏配置",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		hash := resp.RoomHash
		id := ""
		for {
			id = randomString(3, OnlyNumber)
			if _, ok := uno_rooms[id]; ok {
				continue
			} else {
				uno_rooms[id] = &roomUNO{
					groupid: group.GroupId,
					hash:    hash,
					id:      id,
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
	case uno_ExitRoom:
		pi, ok := uno_playerinfo[senderid]
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
						Text: " 下桌失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.UnoC.ExitRoom(pi.playerCtx, &uno_pb.ExitRoomRequest{
			PlayerId: senderid,
			RoomHash: pi.roomHash,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理异常：%v\n", e.String())
			case uno_pb.Errors_RoomStarted:
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
			case uno_pb.Errors_PlayerNoExistAnyRoom:
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
			return
		}
		delete(uno_playerinfo, senderid)
		delete(uno_player2room, senderid)
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
					Text: " 下桌成功，当前桌内玩家：\n" + uno_playersToStr(resp.Players),
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case uno_StartRoom:
		pi, ok := uno_playerinfo[senderid]
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
						Text: " 下桌失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.UnoC.StartRoom(pi.playerCtx, &uno_pb.StartRoomRequest{
			RoomHash: pi.roomHash,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理异常：%v\n", e.String())
			case uno_pb.Errors_RoomStarted:
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
							Text: " 开始游戏失败，你所在桌已开始游戏",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_RoomNoReachPlayers:
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
							Text: " 开始游戏失败，你所在桌玩家数不足，最低需2人",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerNoExistAnyRoom:
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
							Text: " 开始游戏失败，你未在任意桌内",
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
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: "进入选庄家环节，请各位玩家抽牌",
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case uno_GetRooms:
		roomnames := ""
		i := 0
		for _, v := range uno_rooms {
			roomnames += fmt.Sprintf("%v.%v", i+1, v.id)
			if i != len(uno_rooms)-1 {
				roomnames += "\n"
			}
			i++
		}
		if len(uno_rooms) == 0 {
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
	case uno_SendCard_NoSend:
		pi, ok := uno_playerinfo[senderid]
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
						Text: " 下桌失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.UnoC.NoSendCard(pi.playerCtx, &uno_pb.NoSendCardRequest{
			PlayerId: senderid,
			RoomHash: pi.roomHash,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理异常：%v\n", e.String())
			case uno_pb.Errors_PlayerNoDrawCard:
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
							Text: " 不出牌失败，你还未摸牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_RoomNoSendingCard:
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
			case uno_pb.Errors_PlayerNoExistAnyRoom:
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
			case uno_pb.Errors_PlayerNoOperatorNow:
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
			case uno_pb.Errors_PlayerCannotNoSendCard:
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
							Text: " 不出牌失败，你不可以不出牌，可能是受到Draw two等牌，但还未摸罚牌",
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
					TargetId: resp.NextOperator.Id,
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
	case uno_SendCard_Send:
		card, ok := uno_cardStr2Card(text)
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
						Text: " 出牌失败，字符串转牌失败",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		pi, ok := uno_playerinfo[senderid]
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
						Text: " 下桌失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.UnoC.SendCard(pi.playerCtx, &uno_pb.SendCardRequest{
			PlayerId: senderid,
			RoomHash: pi.roomHash,
			SendCard: &card,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理错误类型：%v\n", e.String())
			case uno_pb.Errors_BlackCardNoSpecifiedColor:
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
							Text: " 出牌失败，黑牌未指定颜色(异常错误，不应出现)",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerCannotSendCardFromHandCard:
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
							Text: " 出牌失败，你已摸牌，只能出摸到的牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerCannotSendCard:
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
							Text: " 出牌失败，你不可以出牌，可能是因为受到Draw two等牌效果",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerCardNoExist:
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
			case uno_pb.Errors_SendCardColorORNumberNELastCard:
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
							Text: " 出牌失败，你的牌颜色或数字都与上一张牌不符",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_RoomNoSendingCard:
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
			case uno_pb.Errors_PlayerNoExistAnyRoom:
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
			case uno_pb.Errors_PlayerNoOperatorNow:
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
		extra := ""
		if card.FeatureCard != nil {
			switch card.FeatureCard.FeatureCard {
			case uno_pb.FeatureCards_Wild, uno_pb.FeatureCards_WildDrawFour:
				switch card.FeatureCard.Color {
				case uno_pb.CardColor_Red:
					extra = "黑牌，转变为红色"
				case uno_pb.CardColor_Green:
					extra = "黑牌，转变为绿色"
				case uno_pb.CardColor_Blue:
					extra = "黑牌，转变为蓝色"
				case uno_pb.CardColor_Yellow:
					extra = "黑牌，转变为黄色"
				}
			}
		}
		img, err := uno_getCardImage(card, nil)
		if err != nil {
			logger.Println(err)
			return
		}
		buf, err := image2Buf(img)
		if err != nil {
			logger.Println(err)
			return
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Image,
				Image: &request_pb.MessageChain_Image{
					Buf: buf,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: fmt.Sprintf("%v(%v)出了%v，他还剩下 %v 张牌", sendername, senderid, extra, len(resp.SenderCards)),
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
		cs := []uno_pb.Card{}
		for _, v := range resp.SenderCards {
			cs = append(cs, *v)
		}
		img, err = uno_generateCardsImage(cs, uno_defaultColumn)
		if err != nil {
			logger.Println(err)
			return
		}
		buf, err = image2Buf(img)
		if err != nil {
			logger.Println(err)
			return
		}
		if err := sendMessageToFriend(group.SenderId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Image,
				Image: &request_pb.MessageChain_Image{
					Buf: buf,
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
					TargetId: resp.NextOperator.Id,
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
	case uno_JoinRoom:
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
				ce := resp.GetResponse().GetCmdEvent()
				if ce.Echo != echo {
					continue
				}
				gfl := ce.GetFriendList
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
		r := uno_rooms[id]
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
		resp, err := define.UnoC.JoinRoom(uno_privilegeCtx, &uno_pb.JoinRoomRequest{
			PlayerInfo: &uno_pb.PlayerAccountInfo{
				Id:   senderid,
				Name: sendername,
			},
			RoomHash: r.hash,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		md, ok := metadata.FromOutgoingContext(define.UnoCtx)
		if !ok {
			md = metadata.New(map[string]string{})
		}
		md.Append("Cookie", fmt.Sprintf("player_hash=%v", resp.VerifyHash))
		pctx := metadata.NewOutgoingContext(define.UnoCtx, md)
		uno_playerinfo[senderid] = &playerinfoUNO{
			roomHash:  r.hash,
			hash:      resp.VerifyHash,
			playerCtx: pctx,
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理错误类型：%v\n", e.String())
			case uno_pb.Errors_RoomFull:
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
			case uno_pb.Errors_RoomExistPlayer:
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
			case uno_pb.Errors_RoomNoExist:
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
		uno_player2room[senderid] = r
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
					Text: " 上桌成功，当前桌内玩家：\n" + uno_playersToStr(resp.Players),
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case uno_GetRoom:
		if text == "" {
			r := uno_player2room[senderid]
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
			info, err := uno_getRoom(r.id, r.hash)
			if err != nil {
				logger.Println(err)
				return
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
						Text: fmt.Sprintf(" 获取桌信息成功，你所在桌信息如下：\n%v", info),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		} else {
			r := uno_rooms[text]
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
			info, err := uno_getRoom(r.id, r.hash)
			if err != nil {
				logger.Println(err)
				return
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
						Text: fmt.Sprintf(" 获取桌信息成功，%v 桌信息如下：\n%v", r.id, info),
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
	case uno_DrawCard:
		pi, ok := uno_playerinfo[senderid]
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
						Text: " 下桌失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.UnoC.DrawCard(pi.playerCtx, &uno_pb.DrawCardRequest{
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
			case uno_pb.Errors_PlayerCannotDrawCard:
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
							Text: " 抽牌失败，你不能抽牌，可能是受到Skip等牌的影响",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_RoomNoStart:
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
							Text: " 抽牌失败，你所在桌还未开始游戏",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerAlreadyDrawCard:
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
							Text: " 抽牌失败，你已抽过牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerNoDrawCard:
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
							Text: " 抽牌失败，你不能抽牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerNoOperatorNow:
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
							Text: " 抽牌失败，还未轮到你出牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerNoExistAnyRoom:
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
							Text: " 抽牌失败，你不在任意一个桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		switch {
		case resp.ElectingBanker != nil:
			eb := resp.ElectingBanker
			img, err := uno_getCardImage(*eb.ElectBankerCard, nil)
			if err != nil {
				logger.Println(err)
				return
			}
			buf, err := image2Buf(img)
			if err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Image,
					Image: &request_pb.MessageChain_Image{
						Buf: buf,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_At,
					At: &request_pb.MessageChain_At{
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 你抽到了",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		case resp.SendingCard != nil:
			sc := resp.SendingCard
			cs := []uno_pb.Card{}
			for _, v := range sc.PlayerCard {
				cs = append(cs, *v)
			}
			if sc.DrawCard != nil {
				cs = append(cs, *sc.DrawCard)
			}
			img, err := uno_generateCardsImage(cs, uno_defaultColumn)
			if err != nil {
				logger.Println(err)
				return
			}
			buf, err := image2Buf(img)
			if err != nil {
				logger.Println(err)
				return
			}
			if err := sendMessageToFriend(group.SenderId, []*request_pb.MessageChainObject{
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Image,
					Image: &request_pb.MessageChain_Image{
						Buf: buf,
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
	case uno_CallUNO:
		pi, ok := uno_playerinfo[senderid]
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
						Text: " 报UNO失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.UnoC.CallUNO(pi.playerCtx, &uno_pb.CallUNORequest{
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
			case uno_pb.Errors_PlayerNoExistAnyRoom:
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
							Text: " 报UNO失败，你不在任意一个桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerAlreadyCallUNO:
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
							Text: " 报UNO失败，你已报过UNO",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerCannotCallUNO:
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
							Text: " 报UNO失败，你的手牌数不为2，已被系统罚摸两张牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
				cs := []uno_pb.Card{}
				for _, v := range resp.PlayerCard {
					cs = append(cs, *v)
				}
				img, err := uno_generateCardsImage(cs, uno_defaultColumn)
				if err != nil {
					logger.Println(err)
					return
				}
				buf, err := image2Buf(img)
				if err != nil {
					logger.Println(err)
					return
				}
				if err := sendMessageToFriend(group.SenderId, []*request_pb.MessageChainObject{
					&request_pb.MessageChainObject{
						Type: request_pb.MessageChainType_MessageChainType_Image,
						Image: &request_pb.MessageChain_Image{
							Buf: buf,
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerNoOperatorNow:
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
							Text: " 报UNO失败，还未轮到你出牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_RoomNoSendingCard:
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
							Text: " 报UNO失败，房间还未进入出牌阶段",
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
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: fmt.Sprintf("%v(%v)报了UNO！", sendername, senderid),
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	case uno_Challenge:
		pi, ok := uno_playerinfo[senderid]
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
						Text: " 下桌失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.UnoC.Challenge(pi.playerCtx, &uno_pb.ChallengeRequest{
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
			case uno_pb.Errors_PlayerNoOperatorNow:
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
							Text: " 挑战失败，还未轮到你出牌",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_CannotChallenge:
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
							Text: " 挑战失败，上一张牌不为Wild draw four",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_Challenged:
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
							Text: " 挑战失败，你已挑战过了",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_RoomNoSendingCard:
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
							Text: " 挑战失败，房间还未进入出牌阶段",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		if resp.IsWin {
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
						Text: " 挑战成功，上一位玩家已被系统罚摸四张",
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
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 挑战失败，上一位玩家不存在符合条件的牌，你将被额外罚摸2张牌",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
	case uno_IndicateUNO:
		atTarget := ""
		for _, v := range group.MessageChain {
			if v.Type == response_pb.MessageChainType_MessageChainType_At {
				atTarget = v.At.TargetId
				break
			}
		}
		if atTarget == "" {
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
						Text: " 指出UNO失败，请at你要指出未报UNO的玩家",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		pi, ok := uno_playerinfo[senderid]
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
						Text: " 下桌失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.UnoC.IndicateUNO(pi.playerCtx, &uno_pb.IndicateUNORequest{
			TargetId: atTarget,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if resp.Err != nil {
			switch e := *resp.Err; e {
			default:
				logger.Printf("未处理错误类型：%v\n", e.String())
			case uno_pb.Errors_PlayerIsOperatorNow:
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
							Text: " 指出UNO失败，指定玩家为当前操作者",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerNoExistAnyRoom:
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
							Text: " 指出UNO失败，你不在任意一个桌内",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_RoomNoSendingCard:
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
							Text: " 指出UNO失败，房间还未进入出牌阶段",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerAlreadyCallUNO:
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
							Text: " 指出UNO失败，玩家已报过UNO",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			case uno_pb.Errors_PlayerCannotCallUNO:
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
							Text: " 指出UNO失败，玩家还未达到可报UNO的条件",
						},
					},
				}); err != nil {
					logger.Println(err)
					return
				}
			}
			return
		}
		if resp.IndicateSuccessed {
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
						Text: " 指出UNO成功，玩家已被系统罚摸两张牌",
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
						TargetId: senderid,
					},
				},
				&request_pb.MessageChainObject{
					Type: request_pb.MessageChainType_MessageChainType_Text,
					Text: &request_pb.MessageChain_Text{
						Text: " 指出UNO失败，玩家已报UNO",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
		}
	case uno_GetLastCard:
		r := uno_player2room[senderid]
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
						Text: " 获取上一张牌失败，你未在任意桌内",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		resp, err := define.UnoC.GetRoom(uno_privilegeCtx, &uno_pb.GetRoomRequest{
			RoomHash: r.hash,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		if len(resp.Extra.CardPool) == 0 {
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
						Text: " 获取上一张牌失败，牌池为空",
					},
				},
			}); err != nil {
				logger.Println(err)
				return
			}
			return
		}
		card := resp.Extra.CardPool[len(resp.Extra.CardPool)-1]
		extra := ""
		if card.SendCard.FeatureCard != nil {
			switch card.SendCard.FeatureCard.FeatureCard {
			case uno_pb.FeatureCards_Wild, uno_pb.FeatureCards_WildDrawFour:
				switch card.SendCard.FeatureCard.Color {
				case uno_pb.CardColor_Red:
					extra = "黑牌，转变为红色"
				case uno_pb.CardColor_Green:
					extra = "黑牌，转变为绿色"
				case uno_pb.CardColor_Blue:
					extra = "黑牌，转变为蓝色"
				case uno_pb.CardColor_Yellow:
					extra = "黑牌，转变为黄色"
				}
			}
		}
		img, err := uno_getCardImage(*card.SendCard, nil)
		if err != nil {
			logger.Println(err)
			return
		}
		buf, err := image2Buf(img)
		if err != nil {
			logger.Println(err)
			return
		}
		if err := sendMessageToGroup(group.GroupId, []*request_pb.MessageChainObject{
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Image,
				Image: &request_pb.MessageChain_Image{
					Buf: buf,
				},
			},
			&request_pb.MessageChainObject{
				Type: request_pb.MessageChainType_MessageChainType_Text,
				Text: &request_pb.MessageChain_Text{
					Text: fmt.Sprintf("上一张牌为%v", extra),
				},
			},
		}); err != nil {
			logger.Println(err)
			return
		}
	}
}

func image2Buf(img image.Image) ([]byte, error) {
	var imgBuf bytes.Buffer
	if err := png.Encode(&imgBuf, img); err != nil {
		return nil, err
	}
	return imgBuf.Bytes(), nil
}

func uno_getRoom(id, hash string) (string, error) {
	if hash == "" {
		hash = uno_rooms[id].hash
	}
	resp, err := define.UnoC.GetRoom(uno_privilegeCtx, &uno_pb.GetRoomRequest{
		RoomHash: hash,
	})
	if err != nil {
		return "", err
	}
	ri := resp.Extra
	stageStr := ""
	switch ri.Stage {
	case uno_pb.Stage_WaitingStart:
		stageStr = "等待开始"
	case uno_pb.Stage_ElectingBanker:
		stageStr = "选庄家中"
	case uno_pb.Stage_SendingCard:
		stageStr = "出牌中"
	}
	ps := []*uno_pb.PlayerAccountInfo{}
	for _, v := range ri.Players {
		ps = append(ps, v.PlayerAccountInfo)
	}
	return fmt.Sprintf(`id：%v
	哈希：%v
	游戏状态：%v
	玩家列表：%v`, id, ri.Hash, stageStr, "\n"+uno_playersToStr(ps)), nil
}

func uno_playersToStr(x []*uno_pb.PlayerAccountInfo) string {
	str := ""
	for i, v := range x {
		str += fmt.Sprintf("%v(%v)", v.Name, v.Id)
		if i != len(x)-1 {
			str += "\n"
		}
	}
	if str == "" {
		str = "空"
	}
	return str
}

func uno_cardStr2Card(text string) (uno_pb.Card, bool) {
	text = strings.ToUpper(text)
	switch text {
	case "R0":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_Zero,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "R1":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_One,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "R2":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_Two,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "R3":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_Three,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "R4":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_Four,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "R5":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_Five,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "R6":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_Six,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "R7":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_Seven,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "R8":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_Eight,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "R9":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Red,
				Number: uno_pb.CardNumber_Nine,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "RSKIP":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Red,
				FeatureCard: uno_pb.FeatureCards_Skip,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "RDRAW TWO", "RDT", "R+2":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Red,
				FeatureCard: uno_pb.FeatureCards_DrawTwo,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "RREVERSE", "RREV", "RRE":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Red,
				FeatureCard: uno_pb.FeatureCards_Reverse,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "Y0":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_Zero,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "Y1":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_One,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "Y2":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_Two,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "Y3":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_Three,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "Y4":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_Four,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "Y5":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_Five,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "Y6":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_Six,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "Y7":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_Seven,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "Y8":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_Eight,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "Y9":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Yellow,
				Number: uno_pb.CardNumber_Nine,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "YSKIP":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Yellow,
				FeatureCard: uno_pb.FeatureCards_Skip,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "YDRAW TWO", "YDT", "Y+2":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Yellow,
				FeatureCard: uno_pb.FeatureCards_DrawTwo,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "YREVERSE", "YREV", "YRE":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Yellow,
				FeatureCard: uno_pb.FeatureCards_Reverse,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "B0":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_Zero,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "B1":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_One,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "B2":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_Two,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "B3":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_Three,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "B4":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_Four,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "B5":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_Five,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "B6":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_Six,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "B7":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_Seven,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "B8":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_Eight,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "B9":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Blue,
				Number: uno_pb.CardNumber_Nine,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "BSKIP":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Blue,
				FeatureCard: uno_pb.FeatureCards_Skip,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "BDRAW TWO", "BDT", "B+2":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Blue,
				FeatureCard: uno_pb.FeatureCards_DrawTwo,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "BREVERSE", "BRE", "BREV":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Blue,
				FeatureCard: uno_pb.FeatureCards_Reverse,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "G0":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_Zero,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "G1":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_One,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "G2":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_Two,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "G3":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_Three,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "G4":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_Four,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "G5":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_Five,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "G6":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_Six,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "G7":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_Seven,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "G8":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_Eight,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "G9":
		return uno_pb.Card{
			NormalCard: &uno_pb.NormalCard{
				Color:  uno_pb.CardColor_Green,
				Number: uno_pb.CardNumber_Nine,
			},
			Type: uno_pb.CardType_Normal,
		}, true
	case "GSKIP":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Green,
				FeatureCard: uno_pb.FeatureCards_Skip,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "GDRAW TWO", "GDT", "G+2":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Green,
				FeatureCard: uno_pb.FeatureCards_DrawTwo,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "GREVERSE", "GREV", "GRE":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Green,
				FeatureCard: uno_pb.FeatureCards_Reverse,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "RWILD":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Red,
				FeatureCard: uno_pb.FeatureCards_Wild,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "YWILD":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Yellow,
				FeatureCard: uno_pb.FeatureCards_Wild,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "GWILD":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Green,
				FeatureCard: uno_pb.FeatureCards_Wild,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "BWILD":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Blue,
				FeatureCard: uno_pb.FeatureCards_Wild,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "RWILD DRAW FOUR", "RWDF", "R+4":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Red,
				FeatureCard: uno_pb.FeatureCards_WildDrawFour,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "YWILD DRAW FOUR", "YWDF", "Y+4":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Yellow,
				FeatureCard: uno_pb.FeatureCards_WildDrawFour,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "GWILD DRAW FOUR", "GWDF", "G+4":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Green,
				FeatureCard: uno_pb.FeatureCards_WildDrawFour,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	case "BWILD DRAW FOUR", "BWDF", "B+4":
		return uno_pb.Card{
			FeatureCard: &uno_pb.FeatureCard{
				Color:       uno_pb.CardColor_Blue,
				FeatureCard: uno_pb.FeatureCards_WildDrawFour,
			},
			Type: uno_pb.CardType_Feature,
		}, true
	default:
		return uno_pb.Card{}, false
	}
}

type uno_cardImages struct {
	R0, R1, R2, R3, R4, R5, R6, R7, R8, R9 image.Image
	Y0, Y1, Y2, Y3, Y4, Y5, Y6, Y7, Y8, Y9 image.Image
	B0, B1, B2, B3, B4, B5, B6, B7, B8, B9 image.Image
	G0, G1, G2, G3, G4, G5, G6, G7, G8, G9 image.Image
	RSkip, BSkip, YSkip, GSkip             image.Image
	RDrawTwo, BDrawTwo, YDrawTwo, GDrawTwo image.Image
	RReverse, BReverse, YReverse, GReverse image.Image
	WildDrawTwo, Wild                      image.Image
}

var uno_cardImagesPool = sync.Pool{
	New: func() any {
		x, err := uno_getCardImagesFromDisk()
		if err != nil {
			return err
		}
		return x
	},
}

func uno_getCardImagesFromDisk() (*uno_cardImages, error) {
	ret := new(uno_cardImages)
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R0.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R0 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R1.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R1 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R2.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R2 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R3.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R3 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R4.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R4 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R5.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R5 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R6.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R6 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R7.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R7 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R8.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R8 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "R9.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.R9 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "RSkip.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.RSkip = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "RDraw two.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.RDrawTwo = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "RReverse.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.RReverse = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y0.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y0 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y1.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y1 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y2.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y2 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y3.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y3 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y4.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y4 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y5.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y5 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y6.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y6 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y7.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y7 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y8.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y8 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Y9.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Y9 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "YSkip.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.YSkip = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "YReverse.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.YReverse = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "YDraw two.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.YDrawTwo = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B0.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B0 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B1.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B1 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B2.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B2 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B3.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B3 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B4.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B4 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B5.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B5 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B6.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B6 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B7.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B7 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B8.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B8 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "B9.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.B9 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "BSkip.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.BSkip = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "BReverse.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.BReverse = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "BDraw two.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.BDrawTwo = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G0.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G0 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G1.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G1 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G2.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G2 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G3.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G3 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G4.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G4 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G5.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G5 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G6.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G6 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G7.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G7 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G8.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G8 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "G9.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.G9 = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "GReverse.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.GReverse = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "GSkip.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.GSkip = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "GDraw two.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.GDrawTwo = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Wild draw four.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.WildDrawTwo = img
	}
	if buf, err := os.ReadFile(uno_imageDir + "/" + "Wild.png"); err != nil {
		return nil, err
	} else {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret.Wild = img
	}
	return ret, nil
}

const (
	uno_defaultColumn = 4
)

func uno_getCardImagesFromPool() (*uno_cardImages, error) {
	x := uno_cardImagesPool.Get()
	switch x := x.(type) {
	case *uno_cardImages:
		uno_cardImagesPool.Put(x)
		return x, nil
	case error:
		return nil, x
	default:
		return nil, errors.New("unexcepted type")
	}
}

func uno_getCardImage(card uno_pb.Card, images *uno_cardImages) (image.Image, error) {
	if images == nil {
		x, err := uno_getCardImagesFromPool()
		if err != nil {
			return nil, err
		}
		images = x
	}
	switch card.Type {
	case uno_pb.CardType_Normal:
		switch card.NormalCard.Number {
		case uno_pb.CardNumber_Zero:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R0, nil
			case uno_pb.CardColor_Yellow:
				return images.Y0, nil
			case uno_pb.CardColor_Blue:
				return images.B0, nil
			case uno_pb.CardColor_Green:
				return images.G0, nil
			}
		case uno_pb.CardNumber_One:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R1, nil
			case uno_pb.CardColor_Yellow:
				return images.Y1, nil
			case uno_pb.CardColor_Blue:
				return images.B1, nil
			case uno_pb.CardColor_Green:
				return images.G1, nil
			}
		case uno_pb.CardNumber_Two:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R2, nil
			case uno_pb.CardColor_Yellow:
				return images.Y2, nil
			case uno_pb.CardColor_Blue:
				return images.B2, nil
			case uno_pb.CardColor_Green:
				return images.G2, nil
			}
		case uno_pb.CardNumber_Three:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R3, nil
			case uno_pb.CardColor_Yellow:
				return images.Y3, nil
			case uno_pb.CardColor_Blue:
				return images.B3, nil
			case uno_pb.CardColor_Green:
				return images.G3, nil
			}
		case uno_pb.CardNumber_Four:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R4, nil
			case uno_pb.CardColor_Yellow:
				return images.Y4, nil
			case uno_pb.CardColor_Blue:
				return images.B4, nil
			case uno_pb.CardColor_Green:
				return images.G4, nil
			}
		case uno_pb.CardNumber_Five:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R5, nil
			case uno_pb.CardColor_Yellow:
				return images.Y5, nil
			case uno_pb.CardColor_Blue:
				return images.B5, nil
			case uno_pb.CardColor_Green:
				return images.G5, nil
			}
		case uno_pb.CardNumber_Six:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R6, nil
			case uno_pb.CardColor_Yellow:
				return images.Y6, nil
			case uno_pb.CardColor_Blue:
				return images.B6, nil
			case uno_pb.CardColor_Green:
				return images.G6, nil
			}
		case uno_pb.CardNumber_Seven:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R7, nil
			case uno_pb.CardColor_Yellow:
				return images.Y7, nil
			case uno_pb.CardColor_Blue:
				return images.B7, nil
			case uno_pb.CardColor_Green:
				return images.G7, nil
			}
		case uno_pb.CardNumber_Eight:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R8, nil
			case uno_pb.CardColor_Yellow:
				return images.Y8, nil
			case uno_pb.CardColor_Blue:
				return images.B8, nil
			case uno_pb.CardColor_Green:
				return images.G8, nil
			}
		case uno_pb.CardNumber_Nine:
			switch card.NormalCard.Color {
			case uno_pb.CardColor_Red:
				return images.R9, nil
			case uno_pb.CardColor_Yellow:
				return images.Y9, nil
			case uno_pb.CardColor_Blue:
				return images.B9, nil
			case uno_pb.CardColor_Green:
				return images.G9, nil
			}
		}
	case uno_pb.CardType_Feature:
		switch card.FeatureCard.FeatureCard {
		case uno_pb.FeatureCards_Skip:
			switch card.FeatureCard.Color {
			case uno_pb.CardColor_Red:
				return images.RSkip, nil
			case uno_pb.CardColor_Yellow:
				return images.YSkip, nil
			case uno_pb.CardColor_Blue:
				return images.BSkip, nil
			case uno_pb.CardColor_Green:
				return images.GSkip, nil
			}
		case uno_pb.FeatureCards_Reverse:
			switch card.FeatureCard.Color {
			case uno_pb.CardColor_Red:
				return images.RReverse, nil
			case uno_pb.CardColor_Yellow:
				return images.YReverse, nil
			case uno_pb.CardColor_Blue:
				return images.BReverse, nil
			case uno_pb.CardColor_Green:
				return images.GReverse, nil
			}
		case uno_pb.FeatureCards_DrawTwo:
			switch card.FeatureCard.Color {
			case uno_pb.CardColor_Red:
				return images.RDrawTwo, nil
			case uno_pb.CardColor_Yellow:
				return images.YDrawTwo, nil
			case uno_pb.CardColor_Blue:
				return images.BDrawTwo, nil
			case uno_pb.CardColor_Green:
				return images.GDrawTwo, nil
			}
		case uno_pb.FeatureCards_Wild:
			return images.Wild, nil
		case uno_pb.FeatureCards_WildDrawFour:
			return images.WildDrawTwo, nil
		}
	}
	return nil, errors.New("无匹配牌")
}

const (
	uno_cardWidth  = 240
	uno_cardHeight = 360
)

func uno_generateCardsImage(cards []uno_pb.Card, maxColumn int) (image.Image, error) {
	// 计算底图大小
	if maxColumn <= 0 {
		maxColumn = 10000
	}
	bgX, bgY := 0, 0
	if len(cards) <= maxColumn {
		bgX = uno_cardWidth * len(cards)
		bgY = uno_cardHeight
	} else {
		bgX = uno_cardWidth * maxColumn
		column := 0
		for x := float64(uno_cardWidth * len(cards)); ; {
			if nx := x - float64(uno_cardWidth*maxColumn); x > 0 {
				x = nx
				column++
			} else {
				break
			}
		}
		bgY = uno_cardHeight * column
	}
	bg := image.NewRGBA(image.Rect(0, 0, bgX, bgY))
	images, err := uno_getCardImagesFromPool()
	if err != nil {
		return nil, err
	}
	// 塞入图片
FOROUT:
	for y, n := 0, 0; ; {
		for x := 0; ; {
			cardimg, err := uno_getCardImage(cards[n], images)
			if err != nil {
				return nil, err
			}
			draw.Draw(bg, image.Rect(x, y, x+uno_cardWidth, y+uno_cardHeight), cardimg, cardimg.Bounds().Min, draw.Over)
			n++
			if n > len(cards)-1 {
				break FOROUT
			}
			x += uno_cardWidth
			if uno_cardWidth*maxColumn == x {
				break
			}
		}
		y += uno_cardHeight
	}
	return bg, nil
}

type unoAction int

const (
	uno_Unknown unoAction = iota
	uno_CreateRoom
	uno_ExitRoom
	uno_StartRoom
	uno_GetRooms
	uno_SendCard_NoSend
	uno_SendCard_Send
	uno_JoinRoom
	uno_JoinORExit
	uno_GetRoom
	uno_DrawCard
	uno_CallUNO
	uno_Challenge
	uno_IndicateUNO
	uno_GetLastCard
)
