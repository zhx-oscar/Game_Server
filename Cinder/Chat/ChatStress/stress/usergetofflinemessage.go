package stress

func init() {
	register(getOfflineMessage)
}

func getOfflineMessage(iGo GoroutineIndex, i _RunIndex) {
	user := firstLoginUser(iGo, i)
	user.GetOfflineMessage()
}
