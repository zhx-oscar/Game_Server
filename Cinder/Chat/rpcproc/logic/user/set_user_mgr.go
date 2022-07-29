package user

import (
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user/internal/follow"
	"Cinder/Chat/rpcproc/logic/user/internal/friend"
	"Cinder/Chat/rpcproc/logic/user/internal/oflinfos"
	"Cinder/Chat/rpcproc/logic/user/internal/usrgrp"
)

var userMgr types.IUserManager

func SetUserMgr(mgr types.IUserManager) {
	userMgr = mgr
	friend.SetUserMgr(mgr)
	usrgrp.SetUserMgr(mgr)
	follow.SetUserMgr(mgr)
	oflinfos.SetUserMgr(mgr)
}
