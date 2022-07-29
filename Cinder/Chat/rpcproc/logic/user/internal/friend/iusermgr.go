package friend

import (
	"Cinder/Chat/rpcproc/logic/types"

	assert "github.com/arl/assertgo"
)

type _IUserMgr interface {
	GetUserFriendMgr(userID UserID) types.IFriendMgr
}

// userMgr 可查询用户并取 FriendMgr。
// 避免循环引用 user 包。
var userMgr _IUserMgr

func SetUserMgr(mgr _IUserMgr) {
	assert.True(mgr != nil)
	userMgr = mgr
}
