package bc

import (
	"Cinder/Chat/rpcproc/logic/types"
)

var userMgr types.IUserManager

func SetUserMgr(mgr types.IUserManager) {
	userMgr = mgr
}
