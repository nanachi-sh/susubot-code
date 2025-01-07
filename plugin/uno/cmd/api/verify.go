package api

//检查是否为特权用户(调用无需验证)
func CheckPrivilegeUser(cs []string) bool {
	return false
}
