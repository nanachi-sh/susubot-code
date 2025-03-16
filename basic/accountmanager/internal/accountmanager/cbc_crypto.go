package accountmanager

import (
	"math"
	"strconv"
	"time"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

func decryptPassword(pwd string) string {
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
