package utils

func PortRangeCheck(p int64) bool { return !(p <= 0 || p > 65535) }
