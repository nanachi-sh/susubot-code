package jwt

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/nanachi-sh/susubot-code/basic/jwt/define"
	"github.com/nanachi-sh/susubot-code/basic/jwt/jwt/jwt"
	"github.com/nanachi-sh/susubot-code/basic/jwt/log"
	database_pb "github.com/nanachi-sh/susubot-code/basic/jwt/protos/database"
	jwt_pb "github.com/nanachi-sh/susubot-code/basic/jwt/protos/jwt"
	uno "github.com/nanachi-sh/susubot-code/basic/jwt/protos/qqverifier"
)

var logger = log.Get()

const (
	uno_access_nearExpired = 0.1 //访问JWT仅剩百分之10时可更新
)

func Uno_Sign(req *jwt_pb.Uno_SignRequest) *jwt_pb.Uno_SignResponse {
	if pwd := req.GetPassword(); pwd != nil { //通过密码来获取JWT
		if pwd.Id == "" || pwd.Password == "" {
			return &jwt_pb.Uno_SignResponse{
				Body: &jwt_pb.Uno_SignResponse_Err{
					Err: jwt_pb.Errors_ValueError,
				},
			}
		}
		resp, err := define.DatabaseC.Uno_VerifyUser(define.DatabaseCtx, &database_pb.Uno_VerifyUserRequest{
			Id:       pwd.Id,
			Password: pwd.Password,
		})
		if err != nil {
			logger.Println(err)
			return &jwt_pb.Uno_SignResponse{
				Body: &jwt_pb.Uno_SignResponse_Err{
					Err: jwt_pb.Errors_Undefined,
				},
			}
		}
		if serr := resp.GetErr(); serr != database_pb.Errors_EMPTY {
			switch serr {
			default:
				logger.Printf("未处理错误：%v\n", serr.String())
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_Undefined,
					},
				}
			case database_pb.Errors_UserNoExist:
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_UserNoExist,
					},
				}
			case database_pb.Errors_Undefined:
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_Undefined,
					},
				}
			case database_pb.Errors_UserPasswordWrong:
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_UserPasswordWrong,
					},
				}
			}
		} else if ok := resp.GetOk(); ok {
			resp, err := define.DatabaseC.Uno_GetUser(define.DatabaseCtx, &database_pb.Uno_GetUserRequest{
				Id: pwd.Id,
			})
			if err != nil {
				logger.Println(err)
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_Undefined,
					},
				}
			}
			if serr := resp.GetErr(); serr != database_pb.Errors_EMPTY {
				switch serr {
				default:
					logger.Printf("未处理错误：%v\n", serr.String())
					return &jwt_pb.Uno_SignResponse{
						Body: &jwt_pb.Uno_SignResponse_Err{
							Err: jwt_pb.Errors_Undefined,
						},
					}
				}
			} else if ui := resp.GetUserinfo(); ui != nil {
				access_jwt, err := jwt.Uno_Marshal_AccessJWT(jwt.Uno_UserInfo{
					Id:        ui.Id,
					Name:      ui.Name,
					WinCount:  int(ui.WinCount),
					LoseCount: int(ui.LoseCount),
				})
				if err != nil {
					logger.Println(err)
					return &jwt_pb.Uno_SignResponse{
						Body: &jwt_pb.Uno_SignResponse_Err{
							Err: jwt_pb.Errors_Undefined,
						},
					}
				}
				refresh_jwt, err := jwt.Uno_Marshal_RefreshJWT(ui.Id)
				if err != nil {
					logger.Println(err)
					return &jwt_pb.Uno_SignResponse{
						Body: &jwt_pb.Uno_SignResponse_Err{
							Err: jwt_pb.Errors_Undefined,
						},
					}
				}
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Password{
						Password: &jwt_pb.Uno_SignResponse_FromPassword{
							RefreshJWT: refresh_jwt,
							AccessJWT:  access_jwt,
						},
					},
				}
			} else {
				logger.Println("异常错误")
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_Undefined,
					},
				}
			}
		}
	} else if fromjwt := req.GetJWT(); fromjwt != nil { //通过刷新和访问JWT来获取新的访问JWT
		if fromjwt.AccessJWT == "" || fromjwt.RefreshJWT == "" {
			return &jwt_pb.Uno_SignResponse{
				Body: &jwt_pb.Uno_SignResponse_Err{
					Err: jwt_pb.Errors_ValueError,
				},
			}
		}
		if err := jwt.Uno_Verify(fromjwt.AccessJWT); err != nil {
			switch err {
			default:
				logger.Printf("未处理错误: %v", err)
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_Undefined,
					},
				}
			case jwtgo.ErrTokenExpired:
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_JWTExpired,
					},
				}
			}
		}
		accessBody, err := jwt.Uno_Unmarshal_AccessJWT(fromjwt.AccessJWT)
		if err != nil {
			logger.Println(err)
			return &jwt_pb.Uno_SignResponse{
				Body: &jwt_pb.Uno_SignResponse_Err{Err: jwt_pb.Errors_Undefined},
			}
		}
		iat := accessBody.IssuedAt
		exp := accessBody.ExpiresAt
		// 计算有效时间
		dur := exp.Sub(iat.Time)
		// 计算至少需达到的时间，至少为 dur - dur/near
		after := iat.Add(dur - time.Duration(float64(dur)/uno_access_nearExpired))
		// 判断
		if after.Sub(time.Now()) <= 0 { //已达到更新阈值
			// 验证更新JWT
			if err := jwt.Uno_Verify(fromjwt.RefreshJWT); err != nil {
				switch err {
				default:
					logger.Printf("未处理错误: %v", err)
					return &jwt_pb.Uno_SignResponse{
						Body: &jwt_pb.Uno_SignResponse_Err{
							Err: jwt_pb.Errors_Undefined,
						},
					}
				case jwtgo.ErrTokenExpired:
					return &jwt_pb.Uno_SignResponse{
						Body: &jwt_pb.Uno_SignResponse_Err{
							Err: jwt_pb.Errors_JWTExpired,
						},
					}
				}
			}
			refreshBody, err := jwt.Uno_Unmarshal_RefreshJWT(fromjwt.RefreshJWT)
			if err != nil {
				logger.Println(err)
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{Err: jwt_pb.Errors_Undefined},
				}
			}
			if accessBody.Id != refreshBody.Id {
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{Err: jwt_pb.Errors_JWTInconformity},
				}
			}
			resp, err := define.DatabaseC.Uno_GetUser(define.DatabaseCtx, &database_pb.Uno_GetUserRequest{
				Id: refreshBody.Id,
			})
			if err != nil {
				logger.Println(err)
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{Err: jwt_pb.Errors_Undefined},
				}
			}
			var ui *database_pb.Uno_UserInfo
			if serr := resp.GetErr(); serr != database_pb.Errors_EMPTY {
				switch serr {
				default:
					logger.Printf("未处理错误: %v\n", serr)
					return &jwt_pb.Uno_SignResponse{
						Body: &jwt_pb.Uno_SignResponse_Err{Err: jwt_pb.Errors_Undefined},
					}
				}
			} else if v := resp.GetUserinfo(); v != nil {
				ui = v
			} else {
				logger.Println("异常错误")
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_Undefined,
					},
				}
			}
			newAccessJWT, err := jwt.Uno_Marshal_AccessJWT(jwt.Uno_UserInfo{
				Id:        ui.Id,
				Name:      ui.Name,
				WinCount:  int(ui.WinCount),
				LoseCount: int(ui.LoseCount),
			})
			if err != nil {
				logger.Println("异常错误")
				return &jwt_pb.Uno_SignResponse{
					Body: &jwt_pb.Uno_SignResponse_Err{
						Err: jwt_pb.Errors_Undefined,
					},
				}
			}
			return &jwt_pb.Uno_SignResponse{
				Body: &jwt_pb.Uno_SignResponse_JWT{
					JWT: &jwt_pb.Uno_SignResponse_FromRefreshJWT{
						AccessJWT: newAccessJWT,
					},
				},
			}
		}
	} else {
		return &jwt_pb.Uno_SignResponse{
			Body: &jwt_pb.Uno_SignResponse_Err{
				Err: jwt_pb.Errors_ValueError,
			},
		}
	}
	return &jwt_pb.Uno_SignResponse{
		Body: &jwt_pb.Uno_SignResponse_Err{Err: jwt_pb.Errors_Undefined},
	}
}

type source int

const (
	unknown source = iota
	qq
)

func uno_userSource(id string) (string, source) {
	if len(id) > 2 && strings.ToLower(id[:2]) == "qq" {
		id = "qq" + id[2:]
		return id, qq
	} else {
		return "", unknown
	}
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
		ToString()
}

func Uno_Register(req *jwt_pb.Uno_RegisterRequest) *jwt_pb.Uno_RegisterResponse {
	if req.Id == "" || req.Name == "" || req.Password == "" || req.VerifyHash == "" {
		return &jwt_pb.Uno_RegisterResponse{
			Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_ValueError},
		}
	}
	id, src := uno_userSource(req.Id)
	if src == unknown {
		return &jwt_pb.Uno_RegisterResponse{
			Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_ValueError},
		}
	}
	req.Id = id
	// 解密密码
	passwordPlain := uno_decryptPassword(req.Password)
	// 确认验证结果
	switch src {
	default:
		return &jwt_pb.Uno_RegisterResponse{
			Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_UserUnknownSource},
		}
	case qq:
		resp, err := define.QQVerifierC.Verified(define.QQVerifierCtx, &uno.VerifiedRequest{
			VerifyHash: req.VerifyHash,
		})
		if err != nil {
			logger.Println(err)
			return &jwt_pb.Uno_RegisterResponse{
				Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_Undefined},
			}
		}
		if resp.Err != nil {
			return &jwt_pb.Uno_RegisterResponse{
				Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_UserVerifyError},
			}
		}
		if resp.Result != nil {
			if *resp.Result != uno.Result_Verified {
				return &jwt_pb.Uno_RegisterResponse{
					Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_UserVerifyError},
				}
			}
		}
		if resp.VarifyId != req.Id {
			return &jwt_pb.Uno_RegisterResponse{
				Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_UserVerifyError},
			}
		}
	}
	// 尝试注册
	resp, err := define.DatabaseC.Uno_CreateUser(define.DatabaseCtx, &database_pb.Uno_CreateUserRequest{
		Id:       req.Id,
		Name:     req.Name,
		Password: passwordPlain,
	})
	if err != nil {
		logger.Println(err)
		return &jwt_pb.Uno_RegisterResponse{
			Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_Undefined},
		}
	}
	if serr := resp.GetErr(); serr != database_pb.Errors_EMPTY {
		switch serr {
		default:
			logger.Printf("未处理错误: %v\n", serr.String())
			return &jwt_pb.Uno_RegisterResponse{
				Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_Undefined},
			}
		case database_pb.Errors_UserExist:
			return &jwt_pb.Uno_RegisterResponse{
				Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_UserExist},
			}
		}
	} else if ok := resp.GetOk(); ok {
		access_jwt, err := jwt.Uno_Marshal_AccessJWT(jwt.Uno_UserInfo{
			Id:        req.Id,
			Name:      req.Name,
			WinCount:  0,
			LoseCount: 0,
		})
		if err != nil {
			logger.Println(err)
			return &jwt_pb.Uno_RegisterResponse{
				Body: &jwt_pb.Uno_RegisterResponse_Err{
					Err: jwt_pb.Errors_Undefined,
				},
			}
		}
		refresh_jwt, err := jwt.Uno_Marshal_RefreshJWT(req.Id)
		if err != nil {
			logger.Println(err)
			return &jwt_pb.Uno_RegisterResponse{
				Body: &jwt_pb.Uno_RegisterResponse_Err{
					Err: jwt_pb.Errors_Undefined,
				},
			}
		}
		return &jwt_pb.Uno_RegisterResponse{
			Body: &jwt_pb.Uno_RegisterResponse_JWT{
				JWT: &jwt_pb.Uno_RegisterResponse_Response{
					RefreshJWT: access_jwt,
					AccessJWT:  refresh_jwt,
				},
			},
		}
	} else {
		logger.Println("异常错误")
		return &jwt_pb.Uno_RegisterResponse{
			Body: &jwt_pb.Uno_RegisterResponse_Err{Err: jwt_pb.Errors_Undefined},
		}
	}
}
