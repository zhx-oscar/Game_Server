package stress

import (
	"Cinder/Chat/chatapi"
)

func init() {
	register(loginLogout)
	register(loginSame)
	register(loginLogoutSame)
}

func loginLogout(iGo GoroutineIndex, i _RunIndex) {
	roleID := getRoleID(iGo, i)
	_, errLogin := chatapi.Login(roleID, "nick", []byte("data"))
	panicIfError(errLogin)
	errLogout := chatapi.Logout(roleID)
	panicIfError(errLogout)
}

func loginSame(iGo GoroutineIndex, _ _RunIndex) {
	_, errLogin := chatapi.Login("roleID", "nick", []byte("data"))
	panicIfError(errLogin)
}

func loginLogoutSame(iGo GoroutineIndex, _ _RunIndex) {
	roleID := "roleID"
	_, errLogin := chatapi.Login(roleID, "nick", []byte("data"))
	panicIfError(errLogin)
	errLogout := chatapi.Logout(roleID)
	panicIfError(errLogout)
}
