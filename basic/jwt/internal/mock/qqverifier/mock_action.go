package mock_qqverifierclient

import (
	"context"
	"errors"
	"fmt"

	qqverifierclient "github.com/nanachi-sh/susubot-code/basic/jwt/internal/caller/qqverifier"
	"github.com/nanachi-sh/susubot-code/basic/jwt/pkg/protos/qqverifier"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DefaultMock() *MockQqverifier {
	ctrl := gomock.NewController(nil)
	mock := NewMockQqverifier(ctrl)
	mock.EXPECT().Verified(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, in *qqverifierclient.VerifiedRequest, _ ...any) (*qqverifierclient.VerifiedResponse, error) {
		fmt.Println("Request")
		if in == nil {
			return nil, errors.New("")
		}
		if in.VerifyHash == "123456789" {
			return &qqverifier.VerifiedResponse{
				Result:   qqverifier.Result_Verified.Enum(),
				VarifyId: "100000",
			}, nil
		}
		return nil, status.Error(codes.Unavailable, "")
	})
	return mock
}
