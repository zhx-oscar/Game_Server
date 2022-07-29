package usrgrp

import (
	"Cinder/Chat/rpcproc/logic/types"

	assert "github.com/arl/assertgo"
)

type _IUserMgr interface {
	GetUserGroupMgr(userID UserID) types.IUserGroupMgr
}

var userMgr _IUserMgr

func SetUserMgr(mgr _IUserMgr) {
	assert.True(mgr != nil)
	userMgr = mgr
}
