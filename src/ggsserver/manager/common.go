package manager

func GetUserKey(userId string) string {
	return "user:" + userId + ":info"
}
