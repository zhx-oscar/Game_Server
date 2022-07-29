package stress

func init() {
	register(sendMessage)
}

func sendMessage(iGo GoroutineIndex, i _RunIndex) {
	user := firstLoginUser(iGo, i)
	user.SendMessage("target", []byte("test msg"))
}
