package qqverifier

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/verifier/internal/configs"
	connector_pb "github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/connector"
	handler_request_pb "github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/handler/request"
	handler_response_pb "github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/handler/response"
	verifier_pb "github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/verifier"
	"github.com/twmb/murmur3"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	verifyList       []*verifyinfo
	basicExpiredTime = time.Unix(0, 2145925524123)
)

type verifyinfo struct {
	hash              string
	qqid              string
	code              string
	expiredTime       time.Time //过期时间
	intervalAfterTime time.Time //间隔结束时间
	verified          bool
}

func (vi *verifyinfo) Expired() bool {
	return !(vi.expiredTime.UnixNano() > time.Now().UnixNano()) || vi.expiredTime.Equal(basicExpiredTime)
}

func (vi *verifyinfo) MarkExpired() {
	vi.expiredTime = basicExpiredTime
}

func (vi *verifyinfo) Intervaling() bool {
	return vi.intervalAfterTime.UnixNano() > time.Now().UnixNano()
}

func hash() string {
	buf := make([]byte, 100)
	for n := 0; n != len(buf); n++ {
		buf[n] = byte(rand.Intn(256))
	}
	h1, h2 := murmur3.SeedSum128(rand.Uint64(), rand.Uint64(), buf)
	return fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}

func findVerifysFromQQId(id string) []*verifyinfo {
	ret := []*verifyinfo{}
	for _, v := range verifyList {
		if v.qqid == id {
			ret = append(ret, v)
		}
	}
	return ret
}

func findVerifyFromHash(hash string) (*verifyinfo, bool) {
	for _, v := range verifyList {
		if v.hash == hash {
			return v, true
		}
	}
	return nil, false
}

func NewVerify(logger logx.Logger, in *verifier_pb.QQ_NewVerifyRequest) (*verifier_pb.QQ_NewVerifyResponse_VerifyHash, *verifier_pb.Errors) {
	for _, v := range findVerifysFromQQId(in.QQID) {
		if !v.Intervaling() {
			continue
		}
		return nil, verifier_pb.Errors_Intervaling.Enum()
	}
	if in.Interval == 0 {
		in.Interval = 60 * 1000
	}
	if in.Expires == 0 {
		in.Expires = 300 * 1000
	}
	echo := randomString(6, Mixed)
	resp, err := configs.Call_Handler_Request.GetFriendList(configs.DefaultCtx, &handler_request_pb.BasicRequest{
		Echo: &echo,
	})
	if err != nil {
		logger.Error(err)
		return nil, verifier_pb.Errors_Undefined.Enum()
	}
	stream, err := configs.Call_Connector.Read(configs.DefaultCtx, &connector_pb.Empty{})
	if err != nil {
		logger.Error(err)
		return nil, verifier_pb.Errors_Undefined.Enum()
	}
	if _, err := configs.Call_Connector.Write(configs.DefaultCtx, &connector_pb.WriteRequest{
		Buf: resp.Buf,
	}); err != nil {
		logger.Error(err)
		return nil, verifier_pb.Errors_Undefined.Enum()
	}
	ch := make(chan *verifier_pb.Errors)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	go func() {
		defer close(ch)
		ok := false
	FOROUT:
		for n := 0; n < 10; n++ {
			select {
			case <-ctx.Done():
				return
			default:
			}
			readResp, err := stream.Recv()
			if err != nil {
				logger.Error(err)
				return
			}
			if serr := readResp.GetErr(); serr != connector_pb.Errors_EMPTY {
				switch serr {
				default:
					logger.Errorf("未处理错误: %s", serr.String())
				}
				return
			}
			buf := readResp.GetBuf()
			call_resp, err := configs.Call_Handler_Response.Unmarshal(configs.DefaultCtx, &handler_response_pb.UnmarshalRequest{
				Buf:          buf,
				Type:         handler_response_pb.ResponseType_ResponseType_CmdEvent.Enum(),
				CmdEventType: handler_response_pb.CmdEventType_CmdEventType_GetFriendList.Enum(),
			})
			if err != nil {
				continue
			}
			if serr := call_resp.GetErr(); serr != handler_response_pb.Errors_EMPTY {
				switch serr {
				default:
					logger.Errorf("未处理错误: %s", serr.String())
				}
				return
			}
			resp := call_resp.GetResponse()
			ce := resp.GetCmdEvent()
			if ce.Echo != echo {
				continue
			}
			gfl := ce.GetFriendList
			for _, v := range gfl.Friends {
				if v.UserId == in.QQID {
					ok = true
					break FOROUT
				}
			}
			break
		}
		if !ok {
			ch <- verifier_pb.Errors_NoFriend.Enum()
		} else {
			ch <- nil
		}
	}()
	select {
	case serr := <-ch:
		if serr != nil {
			return nil, serr
		}
	case <-ctx.Done():
		logger.Error("请求好友列表失败")
		return nil, verifier_pb.Errors_Undefined.Enum()
	}
	code := randomString(4, OnlyNumber)
	hash := hash()
	call_req_resp, err := configs.Call_Handler_Request.SendFriendMessage(configs.DefaultCtx, &handler_request_pb.SendFriendMessageRequest{
		FriendId: in.QQID,
		MessageChain: []*handler_request_pb.MessageChainObject{
			&handler_request_pb.MessageChainObject{
				Type: handler_request_pb.MessageChainType_MessageChainType_Text,
				Text: &handler_request_pb.MessageChain_Text{
					Text: fmt.Sprintf("你的验证码为: %s", code),
				},
			},
		},
	})
	if err != nil {
		logger.Error(err)
		return nil, verifier_pb.Errors_Undefined.Enum()
	}
	if _, err := configs.Call_Connector.Write(configs.DefaultCtx, &connector_pb.WriteRequest{Buf: call_req_resp.Buf}); err != nil {
		logger.Error(err)
		return nil, verifier_pb.Errors_Undefined.Enum()
	}
	verifyList = append(verifyList, &verifyinfo{
		hash:              hash,
		qqid:              in.QQID,
		code:              code,
		expiredTime:       time.Now().Add(time.Millisecond * time.Duration(in.Expires)),
		intervalAfterTime: time.Now().Add(time.Millisecond * time.Duration(in.Interval)),
		verified:          false,
	})
	return &verifier_pb.QQ_NewVerifyResponse_VerifyHash{
		VerifyHash: hash,
	}, nil
}

func Verify(req *verifier_pb.QQ_VerifyRequest) (*verifier_pb.QQ_VerifyResponse_Response, *verifier_pb.Errors) {
	vi, ok := findVerifyFromHash(req.VerifyHash)
	if !ok {
		return nil, verifier_pb.Errors_VerifyNoFound.Enum()
	}
	if vi.Expired() {
		return nil, verifier_pb.Errors_Expired.Enum()
	}
	if vi.verified {
		return nil, verifier_pb.Errors_ErrVerified.Enum()
	}
	if vi.code != req.VerifyCode {
		return nil, verifier_pb.Errors_CodeWrong.Enum()
	}
	vi.verified = true
	return &verifier_pb.QQ_VerifyResponse_Response{
		Result:   verifier_pb.Result_Verified,
		VarifyId: vi.qqid,
	}, nil
}

func Verified(req *verifier_pb.QQ_VerifiedRequest) (*verifier_pb.QQ_VerifiedResponse_Response, *verifier_pb.Errors) {
	vi, ok := findVerifyFromHash(req.VerifyHash)
	if !ok {
		return nil, verifier_pb.Errors_VerifyNoFound.Enum()
	}
	if vi.Expired() {
		return nil, verifier_pb.Errors_Expired.Enum()
	}
	if !vi.verified {
		return nil, verifier_pb.Errors_UnVerified.Enum()
	}
	return &verifier_pb.QQ_VerifiedResponse_Response{
		Result:   verifier_pb.Result_Verified,
		VarifyId: vi.qqid,
	}, nil
}

type dictionary int

const (
	OnlyNumber dictionary = iota
	OnlyString
	Mixed
)

func randomString(length int, dic dictionary) string {
	d := ""
	switch dic {
	case OnlyNumber:
		d = "0123456789"
	case OnlyString:
		d = "abcdefghijklmnopqestuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case Mixed:
		d = "0123456789abcdefghijklmnopqestuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	var builder strings.Builder
	for n := 0; n < length; n++ {
		builder.Write([]byte{d[rand.Intn(len(d))]})
	}
	return builder.String()
}
