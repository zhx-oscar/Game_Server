package oflinfos

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"

	assert "github.com/arl/assertgo"
)

type _IUserMgr interface {
	GetUserFriendInfo(userID types.UserID) *chatapi.FriendInfo
}

var userMgr _IUserMgr

func SetUserMgr(mgr _IUserMgr) {
	assert.True(mgr != nil)
	userMgr = mgr
}
