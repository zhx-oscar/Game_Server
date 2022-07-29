package Game

import (
	BaseUser "Cinder/Base/User"
)

var UserMgr BaseUser.IUserMgr

func Init(areaID string, serverID string, userProto BaseUser.IUser, rpcProc interface{}) error {
	return _Init(areaID, serverID, userProto, rpcProc)
}
func Destroy() {
	_Destroy()
}
