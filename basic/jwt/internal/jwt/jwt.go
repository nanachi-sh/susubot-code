package jwt

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/nanachi-sh/susubot-code/basic/jwt/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/jwt/internal/jwt/db"
	unomodel "github.com/nanachi-sh/susubot-code/basic/jwt/internal/model/uno"
	jwt_pb "github.com/nanachi-sh/susubot-code/basic/jwt/pkg/protos/jwt"
	qqverifier_pb "github.com/nanachi-sh/susubot-code/basic/jwt/pkg/protos/qqverifier"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	uno_access_nearExpired = 0.1 //访问JWT仅剩百分之10时可更新
)

type Request struct {
	logger logx.Logger
}

func NewRequest(l logx.Logger) *Request {
	return &Request{
		logger: l,
	}
}

func (r *Request) Uno_Sign(in *jwt_pb.Uno_SignRequest) (*jwt_pb.Uno_SignResponse, error) {
	if fromjwt := in.GetJWT(); fromjwt != nil {
		if fromjwt.AccessJWT == "" || fromjwt.RefreshJWT == "" {
			return &jwt_pb.Uno_SignResponse{}, status.Error(codes.InvalidArgument, "")
		}
		resp, serr := uno_SignFromRefreshJWT(r.logger, fromjwt)
		if serr != nil {
			return &jwt_pb.Uno_SignResponse{
				Body: &jwt_pb.Uno_SignResponse_Err{Err: *serr},
			}, nil
		}
		return &jwt_pb.Uno_SignResponse{
			Body: &jwt_pb.Uno_SignResponse_JWT{JWT: resp},
		}, nil
	} else if frompwd := in.GetPassword(); frompwd != nil {
		if frompwd.Id == "" || frompwd.Password == "" {
			return &jwt_pb.Uno_SignResponse{}, status.Error(codes.InvalidArgument, "")
		}
		resp, serr := uno_SignFromPassword(r.logger, frompwd)
		if serr != nil {
			return &jwt_pb.Uno_SignResponse{
				Body: &jwt_pb.Uno_SignResponse_Err{Err: *serr},
			}, nil
		}
		return &jwt_pb.Uno_SignResponse{
			Body: &jwt_pb.Uno_SignResponse_Password{Password: resp},
		}, nil
	} else {
		return &jwt_pb.Uno_SignResponse{}, status.Error(codes.InvalidArgument, "")
	}
}

func (r *Request) Uno_Register(in *jwt_pb.Uno_RegisterRequest) (*jwt_pb.Uno_RegisterResponse, error) {
	if in.Id == "" || in.Name == "" || in.Password == "" || in.VerifyHash == "" {
		return &jwt_pb.Uno_RegisterResponse{}, status.Error(codes.InvalidArgument, "")
	}
	resp, serr := uno_Register(r.logger, in)
	if serr != nil {
		return &jwt_pb.Uno_RegisterResponse{
			Body: &jwt_pb.Uno_RegisterResponse_Err{Err: *serr},
		}, nil
	}
	return &jwt_pb.Uno_RegisterResponse{
		Body: &jwt_pb.Uno_RegisterResponse_JWT{JWT: resp},
	}, nil
}

func uno_SignFromPassword(logger logx.Logger, in *jwt_pb.Uno_SignRequest_FromPassword) (*jwt_pb.Uno_SignResponse_FromPassword, *jwt_pb.Errors) {
	ok, serr := db.Uno_VerifyUser(logger, in.Id, in.Password)
	if serr != nil {
		return nil, serr
	}
	u, serr := db.Uno_GetUser(logger, in.Id)
	if serr != nil {
		return nil, serr
	}
	access_jwt, ok := uno_Marshal_AccessJWT(logger, u)
	if !ok {
		logger.Error("生成访问JWT失败")
		return nil, jwt_pb.Errors_Undefined.Enum()
	}
	refresh_jwt, ok := uno_Marshal_RefreshJWT(logger, u.Id)
	if !ok {
		logger.Error("生成刷新JWT失败")
		return nil, jwt_pb.Errors_Undefined.Enum()
	}
	return &jwt_pb.Uno_SignResponse_FromPassword{
		RefreshJWT: refresh_jwt,
		AccessJWT:  access_jwt,
	}, nil
}

func uno_SignFromRefreshJWT(logger logx.Logger, in *jwt_pb.Uno_SignRequest_FromRefreshJWT) (*jwt_pb.Uno_SignResponse_FromRefreshJWT, *jwt_pb.Errors) {
	serr := uno_Verify(logger, in.AccessJWT)
	if serr != nil {
		return nil, serr
	}
	accessBody, serr := uno_Unmarshal_AccessJWT(logger, in.AccessJWT)
	if serr != nil {
		return nil, serr
	}
	if !uno_refresh_check(accessBody) {
		return nil, jwt_pb.Errors_Undefined.Enum()
	}
	// 验证更新JWT
	if serr := uno_Verify(logger, in.RefreshJWT); serr != nil {
		return nil, serr
	}
	refreshBody, serr := uno_Unmarshal_RefreshJWT(logger, in.RefreshJWT)
	if serr != nil {
		return nil, serr
	}
	if accessBody.Id != refreshBody.Id {
		return nil, jwt_pb.Errors_JWTInconformity.Enum()
	}
	u, serr := db.Uno_GetUser(logger, refreshBody.Id)
	if serr != nil {
		return nil, serr
	}
	newAccessJWT, ok := uno_Marshal_AccessJWT(logger, u)
	if !ok {
		logger.Error("生成访问JWT失败")
		return nil, jwt_pb.Errors_Undefined.Enum()
	}
	return &jwt_pb.Uno_SignResponse_FromRefreshJWT{
		AccessJWT: newAccessJWT,
	}, nil
}

func uno_refresh_check(accessJWT Uno_Sign_AccessJWT_Body) bool {
	iat := accessJWT.IssuedAt
	exp := accessJWT.ExpiresAt
	// 计算有效时间
	dur := exp.Sub(iat.Time)
	// 计算至少需达到的时间，至少为 dur - dur/near
	after := iat.Add(dur - time.Duration(float64(dur)*uno_access_nearExpired))
	// 判断
	return time.Until(after) <= 0 //已达到更新阈值
}

func uno_decryptPassword(pwd string) string {
	key := strconv.FormatFloat(math.Floor(float64(time.Now().UnixMilli())/30000), 'f', 0, 64)
	return crypto.
		FromBase64String(pwd).
		SetKey(key).
		SetIv(key).
		Des().
		CBC().
		PKCS7Padding().
		Decrypt().
		ToString()
}

func uno_Register(logger logx.Logger, in *jwt_pb.Uno_RegisterRequest) (*jwt_pb.Uno_RegisterResponse_Response, *jwt_pb.Errors) {
	// 解密密码
	passwordPlain := uno_decryptPassword(in.Password)
	if passwordPlain == "" {
		return nil, jwt_pb.Errors_Undefined.Enum()
	}
	// 确认验证结果
	resp, err := configs.Call_QQVerifier.Verified(configs.DefaultCtx, &qqverifier_pb.VerifiedRequest{
		VerifyHash: in.VerifyHash,
	})
	if err != nil {
		logger.Error(err)
		return nil, jwt_pb.Errors_Undefined.Enum()
	}
	if resp.Err != nil {
		switch serr := *resp.Err; serr {
		default:
			logger.Errorf("未处理错误: %s", serr.String())
			return nil, jwt_pb.Errors_Undefined.Enum()
		case qqverifier_pb.Errors_UnVerified, qqverifier_pb.Errors_Expired, qqverifier_pb.Errors_VerifyNoFound:
			return nil, jwt_pb.Errors_UserVerifyError.Enum()
		}
	}
	if resp.Result != nil {
		if *resp.Result != qqverifier_pb.Result_Verified {
			return nil, jwt_pb.Errors_UserVerifyError.Enum()
		}
	}
	if resp.VarifyId != in.Id {
		return nil, jwt_pb.Errors_UserVerifyError.Enum()
	}
	// 尝试注册
	if serr := db.Uno_CreateUser(logger, in.Id, in.Name, passwordPlain); serr != nil {
		return nil, serr
	}
	u, serr := db.Uno_GetUser(logger, in.Id)
	if serr != nil {
		return nil, serr
	}
	access_jwt, ok := uno_Marshal_AccessJWT(logger, u)
	if !ok {
		logger.Error("生成访问JWT失败")
		return nil, jwt_pb.Errors_Undefined.Enum()
	}
	refresh_jwt, ok := uno_Marshal_RefreshJWT(logger, u.Id)
	if !ok {
		logger.Error("生成刷新JWT失败")
		return nil, jwt_pb.Errors_Undefined.Enum()
	}
	return &jwt_pb.Uno_RegisterResponse_Response{
		RefreshJWT: refresh_jwt,
		AccessJWT:  access_jwt,
	}, nil
}

var (
	uno_application     = "uno"
	uno_access_expires  = time.Hour
	uno_access_audience = jwtgo.ClaimStrings{"Access"}

	uno_refresh_expires  = time.Hour * 24 * 14
	uno_refresh_audience = jwtgo.ClaimStrings{"Update"}
)

type (
	Uno_Sign_AccessJWT_Body struct {
		jwtgo.RegisteredClaims
		Name      string `json:"name"`
		Id        string `json:"id"`
		WinCount  int    `json:"winc"`
		LoseCount int    `json:"losec"`
	}
	Uno_Sign_RefreshJWT_Body struct {
		jwtgo.RegisteredClaims
		Id string `json:"id"`
	}
)

const (
	issuer = "root"
)

func uno_Signed(logger logx.Logger, token *jwtgo.Token) (string, bool) {
	signed, err := token.SignedString(configs.JWTKey)
	if err != nil {
		logger.Error(err)
		return "", false
	}
	return signed, true
}

func uno_Marshal_AccessJWT(logger logx.Logger, ui unomodel.Users) (string, bool) {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodES256, Uno_Sign_AccessJWT_Body{
		RegisteredClaims: jwtgo.RegisteredClaims{
			Issuer:    issuer,
			Subject:   uno_application,
			Audience:  uno_access_audience,
			ExpiresAt: &jwtgo.NumericDate{Time: time.Now().Add(uno_access_expires)},
			NotBefore: &jwtgo.NumericDate{Time: time.Now()},
			IssuedAt:  &jwtgo.NumericDate{Time: time.Now()},
			ID:        strconv.FormatInt(rand.Int63(), 10),
		},
		Name:      ui.Name,
		Id:        ui.Id,
		WinCount:  int(ui.WinCount),
		LoseCount: int(ui.LoseCount),
	})
	return uno_Signed(logger, token)
}

func uno_Marshal_RefreshJWT(logger logx.Logger, id string) (string, bool) {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodES256, Uno_Sign_RefreshJWT_Body{
		RegisteredClaims: jwtgo.RegisteredClaims{
			Issuer:    issuer,
			Subject:   uno_application,
			Audience:  uno_refresh_audience,
			ExpiresAt: &jwtgo.NumericDate{Time: time.Now().Add(uno_refresh_expires)},
			NotBefore: &jwtgo.NumericDate{Time: time.Now()},
			IssuedAt:  &jwtgo.NumericDate{Time: time.Now()},
			ID:        strconv.FormatInt(rand.Int63(), 10),
		},
		Id: id,
	})
	return uno_Signed(logger, token)
}

func uno_Verify(logger logx.Logger, jwtStr string) *jwt_pb.Errors {
	if _, err := jwtgo.Parse(jwtStr, func(t *jwtgo.Token) (interface{}, error) {
		return configs.JWTKey.Public(), nil
	}); err != nil {
		switch err {
		default:
			logger.Errorf("未处理错误: %s", err)
		case jwtgo.ErrTokenExpired:
			return jwt_pb.Errors_JWTExpired.Enum()
		case jwtgo.ErrTokenMalformed, jwtgo.ErrSignatureInvalid:
			return jwt_pb.Errors_JWTError.Enum()
		}
		return jwt_pb.Errors_Undefined.Enum()
	}
	return nil
}

func uno_Unmarshal_AccessJWT(logger logx.Logger, jwtStr string) (Uno_Sign_AccessJWT_Body, *jwt_pb.Errors) {
	token, err := jwtgo.ParseWithClaims(jwtStr, new(Uno_Sign_AccessJWT_Body), func(t *jwtgo.Token) (interface{}, error) {
		return configs.JWTKey.Public(), nil
	})
	if err != nil {
		logger.Error(err)
		return Uno_Sign_AccessJWT_Body{}, jwt_pb.Errors_Undefined.Enum()
	}
	if body, ok := token.Claims.(*Uno_Sign_AccessJWT_Body); !ok {
		return Uno_Sign_AccessJWT_Body{}, jwt_pb.Errors_Undefined.Enum()
	} else {
		return *body, nil
	}
}

func uno_Unmarshal_RefreshJWT(logger logx.Logger, jwtStr string) (Uno_Sign_RefreshJWT_Body, *jwt_pb.Errors) {
	token, err := jwtgo.ParseWithClaims(jwtStr, new(Uno_Sign_RefreshJWT_Body), func(t *jwtgo.Token) (interface{}, error) {
		return configs.JWTKey.Public(), nil
	})
	if err != nil {
		logger.Error(err)
		return Uno_Sign_RefreshJWT_Body{}, jwt_pb.Errors_Undefined.Enum()
	}
	if body, ok := token.Claims.(*Uno_Sign_RefreshJWT_Body); !ok {
		return Uno_Sign_RefreshJWT_Body{}, jwt_pb.Errors_Undefined.Enum()
	} else {
		return *body, nil
	}
}
