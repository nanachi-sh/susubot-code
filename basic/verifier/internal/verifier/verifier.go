package verifier

import (
	qqverifier "github.com/nanachi-sh/susubot-code/basic/verifier/internal/verifier/qq"
	verifier_pb "github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/verifier"
	"github.com/zeromicro/go-zero/core/logx"
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

func (r *Request) QQ_NewVerify(in *verifier_pb.QQ_NewVerifyRequest) (*verifier_pb.QQ_NewVerifyResponse, error) {
	if in.QQID == "" {
		return &verifier_pb.QQ_NewVerifyResponse{}, status.Error(codes.InvalidArgument, "")
	}
	resp, serr := qqverifier.NewVerify(r.logger, in)
	if serr != nil {
		return &verifier_pb.QQ_NewVerifyResponse{Body: &verifier_pb.QQ_NewVerifyResponse_Err{Err: *serr}}, nil
	}
	return &verifier_pb.QQ_NewVerifyResponse{Body: resp}, nil
}

func (r *Request) QQ_Verified(in *verifier_pb.QQ_VerifiedRequest) (*verifier_pb.QQ_VerifiedResponse, error) {
	if in.VerifyHash == "" {
		return &verifier_pb.QQ_VerifiedResponse{}, status.Error(codes.InvalidArgument, "")
	}
	resp, serr := qqverifier.Verified(in)
	if serr != nil {
		return &verifier_pb.QQ_VerifiedResponse{Body: &verifier_pb.QQ_VerifiedResponse_Err{Err: *serr}}, nil
	}
	return &verifier_pb.QQ_VerifiedResponse{Body: &verifier_pb.QQ_VerifiedResponse_Resp{Resp: resp}}, nil
}

func (r *Request) QQ_Verify(in *verifier_pb.QQ_VerifyRequest) (*verifier_pb.QQ_VerifyResponse, error) {
	if in.VerifyCode == "" || in.VerifyHash == "" {
		return &verifier_pb.QQ_VerifyResponse{}, status.Error(codes.InvalidArgument, "")
	}
	resp, serr := qqverifier.Verify(in)
	if serr != nil {
		return &verifier_pb.QQ_VerifyResponse{Body: &verifier_pb.QQ_VerifyResponse_Err{Err: *serr}}, nil
	}
	return &verifier_pb.QQ_VerifyResponse{Body: &verifier_pb.QQ_VerifyResponse_Resp{Resp: resp}}, nil
}
