package follow

import (
	"Cinder/Chat/rpcproc/logic/types"

	assert "github.com/arl/assertgo"
)

type _IUserMgr interface {
	GetUserFollowMgr(userID UserID) types.IFollowMgr
}

var userMgr _IUserMgr

func SetUserMgr(mgr _IUserMgr) {
	assert.True(mgr != nil)
	userMgr = mgr
}
