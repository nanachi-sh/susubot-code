package request

import (
	"encoding/json"
	"fmt"

	filewebclient "github.com/nanachi-sh/susubot-code/basic/handler/internal/caller/fileweb"
	"github.com/nanachi-sh/susubot-code/basic/handler/internal/configs"
	request_t "github.com/nanachi-sh/susubot-code/basic/handler/internal/handler/request/types"
	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/fileweb"
	request_pb "github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/request"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	cmd_sendGroupMessage   = "send_group_msg"
	cmd_sendFriendMessage  = "send_private_msg"
	cmd_recall             = "delete_msg"
	cmd_getMessage         = "get_msg"
	cmd_getGroupMemberInfo = "get_group_member_info"
	cmd_getGroupInfo       = "get_group_info"
	cmd_getFriendList      = "get_friend_list"
	cmd_getFriendInfo      = "get_stranger_info"
)

type Request struct {
	logger logx.Logger
}

func NewRequest(l logx.Logger) *Request {
	return &Request{
		logger: l,
	}
}

func (r *Request) GetFriendInfo(in *request_pb.GetFriendInfoRequest) (*request_pb.BasicResponse, error) {
	if in.FriendId == "" {
		return &request_pb.BasicResponse{}, status.Error(codes.InvalidArgument, "")
	}
	buf, ok := getFriendInfo(r.logger, in.FriendId, in.Echo)
	if !ok {
		return &request_pb.BasicResponse{}, status.Error(codes.Unknown, "")
	}
	return &request_pb.BasicResponse{
		Buf: buf,
	}, nil
}
func (r *Request) GetFriendList(in *request_pb.BasicRequest) (*request_pb.BasicResponse, error) {
	buf, ok := getFriendList(r.logger, in.Echo)
	if !ok {
		return &request_pb.BasicResponse{}, status.Error(codes.Unknown, "")
	}
	return &request_pb.BasicResponse{
		Buf: buf,
	}, nil
}
func (r *Request) GetGroupInfo(in *request_pb.GetGroupInfoRequest) (*request_pb.BasicResponse, error) {
	if in.GroupId == "" {
		return &request_pb.BasicResponse{}, status.Error(codes.InvalidArgument, "")
	}
	buf, ok := getGroupInfo(r.logger, in.GroupId, in.Echo)
	if !ok {
		return &request_pb.BasicResponse{}, status.Error(codes.Unknown, "")
	}
	return &request_pb.BasicResponse{
		Buf: buf,
	}, nil
}
func (r *Request) GetGroupMemberInfo(in *request_pb.GetGroupMemberInfoRequest) (*request_pb.BasicResponse, error) {
	if in.GroupId == "" || in.UserId == "" {
		return &request_pb.BasicResponse{}, status.Error(codes.InvalidArgument, "")
	}
	buf, ok := getGroupMemberInfo(r.logger, in.GroupId, in.UserId, in.Echo)
	if !ok {
		return &request_pb.BasicResponse{}, status.Error(codes.Unknown, "")
	}
	return &request_pb.BasicResponse{
		Buf: buf,
	}, nil
}
func (r *Request) GetMessage(in *request_pb.GetMessageRequest) (*request_pb.BasicResponse, error) {
	if in.MessageId == "" {
		return &request_pb.BasicResponse{}, status.Error(codes.InvalidArgument, "")
	}
	buf, ok := getMessage(r.logger, in.MessageId, in.Echo)
	if !ok {
		return &request_pb.BasicResponse{}, status.Error(codes.Unknown, "")
	}
	return &request_pb.BasicResponse{
		Buf: buf,
	}, nil
}
func (r *Request) MessageRecall(in *request_pb.MessageRecallRequest) (*request_pb.BasicResponse, error) {
	if in.MessageId == "" {
		return &request_pb.BasicResponse{}, status.Error(codes.InvalidArgument, "")
	}
	buf, ok := messageRecall(r.logger, in.MessageId, in.Echo)
	if !ok {
		return &request_pb.BasicResponse{}, status.Error(codes.Unknown, "")
	}
	return &request_pb.BasicResponse{
		Buf: buf,
	}, nil
}
func (r *Request) SendFriendMessage(in *request_pb.SendFriendMessageRequest) (*request_pb.BasicResponse, error) {
	if in.FriendId == "" || len(in.MessageChain) == 0 {
		return &request_pb.BasicResponse{}, status.Error(codes.InvalidArgument, "")
	}
	buf, ok := sendFriendMessage(r.logger, in.FriendId, in.MessageChain, in.Echo)
	if !ok {
		return &request_pb.BasicResponse{}, status.Error(codes.Unknown, "")
	}
	return &request_pb.BasicResponse{
		Buf: buf,
	}, nil
}
func (r *Request) SendGroupMessage(in *request_pb.SendGroupMessageRequest) (*request_pb.BasicResponse, error) {
	if in.GroupId == "" || len(in.MessageChain) == 0 {
		return &request_pb.BasicResponse{}, status.Error(codes.InvalidArgument, "")
	}
	buf, ok := sendGroupMessage(r.logger, in.GroupId, in.MessageChain, in.Echo)
	if !ok {
		return &request_pb.BasicResponse{}, status.Error(codes.Unknown, "")
	}
	return &request_pb.BasicResponse{
		Buf: buf,
	}, nil
}

func getFriendInfo(logger logx.Logger, friendid string, echo *string) ([]byte, bool) {
	req := new(request_t.Request)
	req.Action = cmd_getFriendInfo
	if echo != nil {
		req.Echo = *echo
	}
	req.Params = new(request_t.Request_Params)
	req.Params.UserId = friendid
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return buf, true
}

func getFriendList(logger logx.Logger, echo *string) ([]byte, bool) {
	req := new(request_t.Request)
	req.Action = cmd_getFriendList
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return buf, true
}

func getGroupInfo(logger logx.Logger, groupid string, echo *string) ([]byte, bool) {
	req := new(request_t.Request)
	req.Action = cmd_getGroupInfo
	req.Params = new(request_t.Request_Params)
	req.Params.GroupId = groupid
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return buf, true
}

func getGroupMemberInfo(logger logx.Logger, groupid, memberid string, echo *string) ([]byte, bool) {
	req := new(request_t.Request)
	req.Action = cmd_getGroupMemberInfo
	req.Params = new(request_t.Request_Params)
	req.Params.GroupId = groupid
	req.Params.UserId = memberid
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return buf, true
}

func messageRecall(logger logx.Logger, messageid string, echo *string) ([]byte, bool) {
	req := new(request_t.Request)
	req.Action = cmd_recall
	req.Params = new(request_t.Request_Params)
	req.Params.MessageId = messageid
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return buf, true
}

func sendGroupMessage(logger logx.Logger, groupid string, inMcs []*request_pb.MessageChainObject, echo *string) ([]byte, bool) {
	req := new(request_t.Request)
	req.Action = cmd_sendGroupMessage
	req.Params = new(request_t.Request_Params)
	req.Params.GroupId = groupid
	mcs, ok := marshalMessageChain(logger, inMcs)
	if !ok {
		return nil, false
	}
	var mcs_j []map[string]any
	for _, v := range mcs {
		d, err := json.Marshal(v)
		if err != nil {
			logger.Error(err)
			return nil, false
		}
		var m map[string]any
		if err := json.Unmarshal(d, &m); err != nil {
			logger.Error(err)
			return nil, false
		}
		mcs_j = append(mcs_j, m)
	}
	req.Params.Message = mcs_j
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return buf, true
}

func sendFriendMessage(logger logx.Logger, friendid string, inMcs []*request_pb.MessageChainObject, echo *string) ([]byte, bool) {
	req := new(request_t.Request)
	req.Action = cmd_sendFriendMessage
	req.Params = new(request_t.Request_Params)
	req.Params.UserId = friendid
	mcs, ok := marshalMessageChain(logger, inMcs)
	if !ok {
		return nil, false
	}
	var mcs_j []map[string]any
	for _, v := range mcs {
		d, err := json.Marshal(v)
		if err != nil {
			logger.Error(err)
			return nil, false
		}
		var m map[string]any
		if err := json.Unmarshal(d, &m); err != nil {
			logger.Error(err)
			return nil, false
		}
		mcs_j = append(mcs_j, m)
	}
	req.Params.Message = mcs_j
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return buf, true
}

func getMessage(logger logx.Logger, messageid string, echo *string) ([]byte, bool) {
	req := new(request_t.Request)
	req.Action = cmd_getMessage
	req.Params = new(request_t.Request_Params)
	req.Params.MessageId = messageid
	if echo != nil {
		req.Echo = *echo
	}
	buf, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return buf, true
}

func marshalMessageChain(logger logx.Logger, mc []*request_pb.MessageChainObject) ([]*request_t.MessageChain, bool) {
	ret := []*request_t.MessageChain{}
	for _, v := range mc {
		switch v.Type {
		case request_pb.MessageChainType_MessageChainType_Text:
			if v.Text == nil {
				logger.Error("消息链Text结构体为nil")
				return nil, false
			}
			ret = append(ret, &request_t.MessageChain{
				Data: map[string]any{
					"text": v.Text.Text,
				},
				Type: "text",
			})
		case request_pb.MessageChainType_MessageChainType_At:
			if v.At == nil {
				logger.Error("消息链At结构体为nil")
				return nil, false
			}
			ret = append(ret, &request_t.MessageChain{
				Data: map[string]any{
					"qq": v.At.TargetId,
				},
				Type: "at",
			})
		case request_pb.MessageChainType_MessageChainType_Reply:
			if v.Reply == nil {
				logger.Error("消息链Reply结构体为nil")
				return nil, false
			}
			ret = append(ret, &request_t.MessageChain{
				Data: map[string]any{
					"id": v.Reply.MessageId,
				},
				Type: "reply",
			})
		case request_pb.MessageChainType_MessageChainType_Image:
			image := v.Image
			if image == nil {
				logger.Error("消息链Image结构体为nil")
				return nil, false
			}
			u := ""
			if image.URL != nil {
				u = *image.URL
			} else if image.Buf != nil {
				resp, err := configs.Call_FileWeb.Upload(configs.DefaultCtx, &fileweb.UploadRequest{
					Buf: image.Buf,
				})
				if err != nil {
					logger.Error(err)
					return nil, false
				}
				if serr := resp.GetErr(); serr != fileweb.Errors_EMPTY {
					switch serr {
					default:
						logger.Errorf("未处理错误: %s", serr.String())
					}
					return nil, false
				}
				u = fmt.Sprintf("%s/%s", configs.ASSETS_URL, resp.GetHash())
			}
			ret = append(ret, &request_t.MessageChain{
				Data: map[string]any{
					"file": u,
				},
				Type: "image",
			})
		case request_pb.MessageChainType_MessageChainType_Voice: //Voice
			voice := v.Voice
			if voice == nil {
				logger.Error("消息链Voice结构体为nil")
				return nil, false
			}
			u := ""
			if voice.URL != nil {
				u = *voice.URL
			} else if voice.Buf != nil {
				resp, err := configs.Call_FileWeb.Upload(configs.DefaultCtx, &filewebclient.UploadRequest{
					Buf: voice.Buf,
				})
				if err != nil {
					logger.Error(err)
					return nil, false
				}
				if serr := resp.GetErr(); serr != fileweb.Errors_EMPTY {
					switch serr {
					default:
						logger.Errorf("未处理错误: %s", serr.String())
					}
					return nil, false
				}
				u = fmt.Sprintf("%s/%s", configs.ASSETS_URL, resp.GetHash())
			}
			ret = append(ret, &request_t.MessageChain{
				Data: map[string]any{
					"file": u,
				},
				Type: "record",
			})
		case request_pb.MessageChainType_MessageChainType_Video:
			video := v.Video
			if video == nil {
				logger.Error("消息链Video结构体为nil")
				return nil, false
			}
			u := ""
			if video.URL != nil {
				u = *video.URL
			} else if video.Buf != nil {
				resp, err := configs.Call_FileWeb.Upload(configs.DefaultCtx, &fileweb.UploadRequest{
					Buf: video.Buf,
				})
				if err != nil {
					logger.Error(err)
					return nil, false
				}
				if serr := resp.GetErr(); serr != fileweb.Errors_EMPTY {
					switch serr {
					default:
						logger.Errorf("未处理错误: %s", serr.String())
					}
					return nil, false
				}
				u = fmt.Sprintf("%s/%s", configs.ASSETS_URL, resp.GetHash())
			}
			ret = append(ret, &request_t.MessageChain{
				Data: map[string]any{
					"file": u,
				},
				Type: "video",
			})
		default:
		}
	}
	return ret, true
}
