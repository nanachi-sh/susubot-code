package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/nanachi-sh/susubot-code/basic/jwt/define"
	"github.com/nanachi-sh/susubot-code/basic/jwt/log"
)

type Uno_UserInfo struct {
	Id, Name            string
	WinCount, LoseCount int
}

var (
	logger  = log.Get()
	ecdsaPK *ecdsa.PrivateKey

	uno_application     = "uno"
	uno_access_expires  = time.Hour
	uno_access_audience = jwt.ClaimStrings{"Access"}

	uno_refresh_expires  = time.Hour * 24 * 14
	uno_refresh_audience = jwt.ClaimStrings{"Update"}
)

type (
	Uno_Sign_AccessJWT_Body struct {
		jwt.RegisteredClaims
		Name      string `json:"name"`
		Id        string `json:"id"`
		WinCount  int    `json:"winc"`
		LoseCount int    `json:"losec"`
	}
	Uno_Sign_RefreshJWT_Body struct {
		jwt.RegisteredClaims
		Id string `json:"id"`
	}
)

const (
	issuer = "root"
)

func init() {
	buf, err := os.ReadFile(fmt.Sprintf("%v/jwt.key", define.CertDir))
	if err != nil {
		logger.Fatalln(err)
	}
	pblock, _ := pem.Decode(buf)
	if pblock == nil {
		logger.Fatalln("JWT私钥不正确")
	}
	pk, err := x509.ParseECPrivateKey(pblock.Bytes)
	if err != nil {
		logger.Fatalln(err)
	}
	ecdsaPK = pk
}

func uno_Sign(token *jwt.Token) (string, error) {
	return token.SignedString(ecdsaPK)
}

func Uno_Marshal_AccessJWT(ui Uno_UserInfo) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, Uno_Sign_AccessJWT_Body{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   uno_application,
			Audience:  uno_access_audience,
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(uno_access_expires)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        strconv.FormatInt(rand.Int63(), 10),
		},
		Name:      ui.Name,
		Id:        ui.Id,
		WinCount:  ui.WinCount,
		LoseCount: ui.LoseCount,
	})
	return uno_Sign(token)
}

func Uno_Marshal_RefreshJWT(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, Uno_Sign_RefreshJWT_Body{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   uno_application,
			Audience:  uno_refresh_audience,
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(uno_refresh_expires)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        strconv.FormatInt(rand.Int63(), 10),
		},
		Id: id,
	})
	return uno_Sign(token)
}

func Uno_Verify(jwtStr string) error {
	_, err := jwt.Parse(jwtStr, func(t *jwt.Token) (interface{}, error) {
		return ecdsaPK.Public(), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func Uno_Unmarshal_AccessJWT(jwtStr string) (Uno_Sign_AccessJWT_Body, error) {
	token, err := jwt.ParseWithClaims(jwtStr, new(Uno_Sign_AccessJWT_Body), func(t *jwt.Token) (interface{}, error) {
		return ecdsaPK.Public(), nil
	})
	if err != nil {
		return Uno_Sign_AccessJWT_Body{}, err
	}
	if body, ok := token.Claims.(*Uno_Sign_AccessJWT_Body); !ok {
		return Uno_Sign_AccessJWT_Body{}, errors.New("异常错误")
	} else {
		return *body, nil
	}
}

func Uno_Unmarshal_RefreshJWT(jwtStr string) (Uno_Sign_RefreshJWT_Body, error) {
	token, err := jwt.ParseWithClaims(jwtStr, new(Uno_Sign_RefreshJWT_Body), func(t *jwt.Token) (interface{}, error) {
		return ecdsaPK.Public(), nil
	})
	if err != nil {
		return Uno_Sign_RefreshJWT_Body{}, err
	}
	if body, ok := token.Claims.(*Uno_Sign_RefreshJWT_Body); !ok {
		return Uno_Sign_RefreshJWT_Body{}, errors.New("异常错误")
	} else {
		return *body, nil
	}
}
