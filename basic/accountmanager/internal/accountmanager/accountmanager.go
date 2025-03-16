package accountmanager

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/model/applications"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/utils"
	accountmanager_pb "github.com/nanachi-sh/susubot-code/basic/accountmanager/pkg/protos/accountmanager"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/pkg/types"
	"github.com/zeromicro/go-zero/core/logx"
	smtp "gopkg.in/gomail.v2"
)

type Request struct {
	logger logx.Logger
}

var (
	smtpDialer = smtp.NewDialer("smtp.qq.com", 587, configs.SMTP_USERNAME, configs.SMTP_AUTHCODE)
)

func NewRequest(l logx.Logger) *Request {
	return &Request{
		logger: l,
	}
}

func (r *Request) UserRegister(req *types.UserRegisterRequest) (resp any, err error) {
	if !checkEmailValid(req.Email) {
		err = types.NewError(accountmanager_pb.Error_ERROR_INVALID_ARGUMENT, "邮箱格式错误，或为不支持的邮箱，支持邮箱如下：QQ邮箱，网易邮箱，Google邮箱")
		return
	}
	return userRegister(r.logger, req)
}

func (r *Request) VerifyCode() (resp any, err error) {
	return verifyCode(r.logger)
}

func (r *Request) UserVerifyCode_Email(req *types.UserVerifyCodeEmailRequest) (resp any, err error) {
	if !checkEmailValid(req.Email) {
		err = types.NewError(accountmanager_pb.Error_ERROR_INVALID_ARGUMENT, "邮箱格式错误，或为不支持的邮箱，支持邮箱如下：QQ邮箱，网易邮箱，Google邮箱")
		return
	}
	return userVerifyCode_Email(r.logger, req)
}

func checkEmailValid(email string) bool {
	ok, err := regexp.MatchString(`^\w{1,64}@(qq\.com|vip\.qq\.com|126\.com|163\.com|gmail\.com)$`, email)
	if err != nil {
		panic(err)
	}
	return ok
}

func userRegister(logger logx.Logger, req *types.UserRegisterRequest) (resp any, err error) {
	rKey := fmt.Sprintf("verifycode_email_%s", req.Email)
	email, err := configs.Redis.Hget(rKey, "email")
	if err != nil {
		logger.Error(err)
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
		return
	}
	if email != req.Email {
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "请求邮箱与数据库中不同，异常错误")
		return
	}
	code, err := configs.Redis.Hget(rKey, "verify_code")
	if err != nil {
		logger.Error(err)
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
		return
	}
	if req.VerifyCode != code {
		err = types.NewError(accountmanager_pb.Error_ERROR_EMAIL_VERIFYCODE_ANSWER_FAIL, "")
		return
	}
	id := ""
	{
		result, err := configs.LDAP.Search(ldap.NewSearchRequest(
			configs.LDAP_BASIC_DN,
			ldap.ScopeSingleLevel,
			ldap.DerefFindingBaseObj,
			10000, 5, false, "(objectClass=inetOrgPerson)",
			nil, nil,
		))
		if err != nil {
			logger.Error(err)
			err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
			return nil, err
		}
		uids := []int64{}
		for _, v := range result.Entries {
			uid, err := strconv.ParseInt(v.GetAttributeValue("uid"), 10, 0)
			if err != nil {
				continue
			}
			uids = append(uids, uid)
		}
		if len(uids) > 0 {
			sort.Slice(uids, func(i, j int) bool {
				return uids[i] > uids[j]
			})
			id = strconv.FormatInt(uids[0]+1, 10)
		} else {
			id = "1000000000"
		}
	}
	l_req := ldap.NewAddRequest(fmt.Sprintf("uid=%s,%s", id, configs.LDAP_BASIC_DN), nil)
	l_req.Attribute("objectClass", []string{"inetOrgPerson"})
	l_req.Attribute("cn", []string{utils.RandomString(10, utils.Dict_Mixed)})
	l_req.Attribute("sn", []string{utils.RandomString(10, utils.Dict_Mixed)})
	l_req.Attribute("displayName", []string{"用户_" + utils.RandomString(16, utils.Dict_Number)})
	l_req.Attribute("mail", []string{req.Email})
	passwordPlain := decryptPassword(req.Password)
	if passwordPlain == "" {
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "密码解密失败")
		return
	}
	ssha_pwd, err := SSHAEncoder{}.Encode([]byte(passwordPlain))
	if err != nil {
		logger.Error(err)
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
		return
	}
	l_req.Attribute("userPassword", []string{string(ssha_pwd)})
	if err = configs.LDAP.Add(l_req); err != nil {
		logger.Error(err)
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
		return
	}
	if _, err = configs.Model_Applications.Insert(context.Background(), &applications.Twoonone{
		Id:               id,
		Wincount:         0,
		Losecount:        0,
		LastGetdaliyTime: 0,
		Coin:             0,
	}); err != nil {
		logger.Error(err)
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
		if err := configs.LDAP.Del(ldap.NewDelRequest(fmt.Sprintf("uid=%s,%s", id, configs.LDAP_BASIC_DN), nil)); err != nil {
			logger.Error(err)
		}
		return
	}
	return &accountmanager_pb.UserRegisterResponse{}, nil
}

func userVerifyCode_Email(logger logx.Logger, req *types.UserVerifyCodeEmailRequest) (resp any, err error) {
	code := utils.RandomString(4, utils.Dict_Number)
	rKey := fmt.Sprintf("verifycode_email_%s", req.Email)
	if err = configs.Redis.Hset(rKey, "verify_code", code); err != nil {
		logger.Error(err)
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
		return
	}
	if err = configs.Redis.Hset(rKey, "email", req.Email); err != nil {
		logger.Error(err)
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
		return
	}
	if err = configs.Redis.Expire(rKey, 30*60); err != nil {
		logger.Error(err)
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
		return
	}
	msg := smtp.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("nobody <%s>", configs.SMTP_USERNAME))
	msg.SetHeader("To", req.Email)
	msg.SetHeader("Subject", "验证码")
	msg.SetBody("text/html", fmt.Sprintf("你的验证码为：%s 请在三十分钟内使用", code))
	if err = smtpDialer.DialAndSend(msg); err != nil {
		logger.Error(err)
		err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "内部错误")
		if _, err = configs.Redis.Del(rKey); err != nil {
			logger.Error(err)
		}
		return
	}
	return &accountmanager_pb.UserVerifyCodeEmailResponse{}, nil
}

func verifyCode(logger logx.Logger) (resp any, err error) {
	id, b64, err := func(logger logx.Logger) (id string, b64 string, err error) {
		id, b64, _, err = configs.Captcha.Generate()
		if err != nil {
			logger.Error(err)
			err = types.NewError(accountmanager_pb.Error_ERROR_UNDEFINED, "")
			return
		}
		return
	}(logger)
	if err != nil {
		return nil, err
	}
	return &accountmanager_pb.VerifyCodeResponse{
		B64:      b64,
		VerifyId: id,
	}, nil
}
