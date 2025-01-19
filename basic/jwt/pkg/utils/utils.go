package utils

import (
	"math/rand"
	"strings"
)

func PortRangeCheck(p int64) bool { return !(p <= 0 || p > 65535) }

const Dict string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int, dict string) string {
	ret := new(strings.Builder)
	for i := 0; i < length; i++ {
		ret.WriteByte(dict[rand.Intn(len(dict))])
	}
	return ret.String()
}
