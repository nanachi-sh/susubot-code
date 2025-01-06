package qqverifier

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/qqverifier/define"
	"github.com/nanachi-sh/susubot-code/basic/qqverifier/protos/connector"
	"github.com/nanachi-sh/susubot-code/basic/qqverifier/protos/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/qqverifier/protos/handler/response"
	qqverifier_pb "github.com/nanachi-sh/susubot-code/basic/qqverifier/protos/qqverifier"
	"github.com/twmb/murmur3"
)

var (
	verifyList []*verifyinfo
)

type verifyinfo struct {
	hash              string
	qqid              string
	code              string
	expiredTime       time.Time //过期时间
	intervalAfterTime time.Time //间隔结束时间
	verified          bool
}

func hash() string {
	buf := make([]byte, 100)
	for n := 0; n != len(buf); n++ {
		buf[n] = byte(rand.Intn(256))
	}
	h1, h2 := murmur3.SeedSum128(rand.Uint64(), rand.Uint64(), buf)
	return fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}

func findVerifyFromQQId(id string) (*verifyinfo, bool) {
	for _, v := range verifyList {
		if v.qqid == id {
			return v, true
		}
	}
	return nil, false
}

func findVerifyFromHash(hash string) (*verifyinfo, bool) {
	for _, v := range verifyList {
		if v.hash == hash {
			return v, true
		}
	}
	return nil, false
}

func NewVerify(req *qqverifier_pb.NewVerifyRequest) (*qqverifier_pb.NewVerifyResponse, error) {
	if vi, ok := findVerifyFromQQId(req.QQID); ok {
		if vi.intervalAfterTime.UnixNano() > time.Now().UnixNano() {
			return &qqverifier_pb.NewVerifyResponse{
				Err: qqverifier_pb.Errors_Intervaling.Enum(),
			}, nil
		} else {
			if !vi.verified {
				vi.expiredTime = time.Unix(0, 0)
			}
		}
	}
	if req.QQID == "" {
		return nil, errors.New("QQID不能为空")
	}
	if req.Interval == 0 {
		req.Interval = 60 * 1000
	}
	if req.Expires == 0 {
		req.Interval = 300 * 1000
	}
	echo := randomString(6, Mixed)
	resp, err := define.RequestC.GetFriendList(define.HandlerCtx, &request.BasicRequest{
		Echo: &echo,
	})
	if err != nil {
		return nil, err
	}
	stream, err := define.ConnectorC.Read(define.ConnectorCtx, &connector.Empty{})
	if err != nil {
		return nil, err
	}
	if _, err := define.ConnectorC.Write(define.ConnectorCtx, &connector.WriteRequest{
		Buf: resp.Buf,
	}); err != nil {
		return nil, err
	}
	ok := false
FOROUT:
	for n := 0; n < 10; n++ {
		readResp, err := stream.Recv()
		if err != nil {
			return nil, err
		}
		resp, err := define.ResponseC.Unmarshal(define.HandlerCtx, &response.UnmarshalRequest{
			Buf:          readResp.Buf,
			Type:         response.ResponseType_ResponseType_CmdEvent.Enum(),
			CmdEventType: response.CmdEventType_CmdEventType_GetFriendList.Enum(),
		})
		if err != nil {
			continue
		}
		if resp.CmdEvent.Echo != echo {
			continue
		}
		gfl := resp.CmdEvent.GetFriendList
		for _, v := range gfl.Friends {
			if v.UserId == req.QQID {
				ok = true
				break FOROUT
			}
		}
	}
	if !ok {
		return &qqverifier_pb.NewVerifyResponse{
			Err: qqverifier_pb.Errors_NoFriend.Enum(),
		}, nil
	}
	code := randomString(4, OnlyNumber)
	hash := hash()
	reqResp, err := define.RequestC.SendFriendMessage(define.HandlerCtx, &request.SendFriendMessageRequest{
		FriendId: req.QQID,
		MessageChain: []*request.MessageChainObject{
			&request.MessageChainObject{
				Type: request.MessageChainType_MessageChainType_Text,
				Text: &request.MessageChain_Text{
					Text: fmt.Sprintf("你的验证码为: %v", code),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if _, err := define.ConnectorC.Write(define.ConnectorCtx, &connector.WriteRequest{Buf: reqResp.Buf}); err != nil {
		return nil, err
	}
	verifyList = append(verifyList, &verifyinfo{
		hash:              hash,
		qqid:              req.QQID,
		code:              code,
		expiredTime:       time.Now().Add(time.Millisecond * time.Duration(req.Expires)),
		intervalAfterTime: time.Now().Add(time.Millisecond * time.Duration(req.Interval)),
		verified:          false,
	})
	return &qqverifier_pb.NewVerifyResponse{
		VerifyHash: hash,
	}, nil
}

func Verify(req *qqverifier_pb.VerifyRequest) (*qqverifier_pb.VerifyResponse, error) {
	if req.VerifyHash == "" || req.VerifyCode == "" {
		return nil, errors.New("Hash或Code不能为空")
	}
	vi, ok := findVerifyFromHash(req.VerifyHash)
	if !ok {
		return &qqverifier_pb.VerifyResponse{
			Err: qqverifier_pb.Errors_VerifyNoFound.Enum(),
		}, nil
	}
	if vi.expiredTime.UnixNano() > time.Now().UnixNano() {
		return &qqverifier_pb.VerifyResponse{
			Err: qqverifier_pb.Errors_Expired.Enum(),
		}, nil
	}
	if vi.verified {
		return &qqverifier_pb.VerifyResponse{
			Err: qqverifier_pb.Errors_ErrVerified.Enum(),
		}, nil
	}
	if vi.code != req.VerifyCode {
		return &qqverifier_pb.VerifyResponse{
			Err: qqverifier_pb.Errors_CodeWrong.Enum(),
		}, nil
	}
	vi.verified = true
	return &qqverifier_pb.VerifyResponse{
		Result: qqverifier_pb.Result_Verified.Enum(),
	}, nil
}

func Verified(req *qqverifier_pb.VerifiedRequest) (*qqverifier_pb.VerifiedResponse, error) {
	if req.VerifyHash == "" {
		return nil, errors.New("Hash不能为空")
	}
	vi, ok := findVerifyFromHash(req.VerifyHash)
	if !ok {
		return &qqverifier_pb.VerifiedResponse{
			Err: qqverifier_pb.Errors_VerifyNoFound.Enum(),
		}, nil
	}
	if vi.expiredTime.UnixNano() > time.Now().UnixNano() {
		return &qqverifier_pb.VerifiedResponse{
			Err: qqverifier_pb.Errors_Expired.Enum(),
		}, nil
	}
	if !vi.verified {
		return &qqverifier_pb.VerifiedResponse{
			Err: qqverifier_pb.Errors_UnVerified.Enum(),
		}, nil
	}
	return &qqverifier_pb.VerifiedResponse{
		Result: qqverifier_pb.Result_Verified.Enum(),
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
