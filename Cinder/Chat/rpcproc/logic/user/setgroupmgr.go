package user

import (
	"Cinder/Chat/rpcproc/logic/user/internal/usrgrp"

	assert "github.com/arl/assertgo"
)

// 需要在初始化时 SetGroupMgr().
func SetGroupMgr(mgr usrgrp.IGroupMgr) {
	assert.True(mgr != nil)
	usrgrp.SetGroupMgr(mgr)
}
