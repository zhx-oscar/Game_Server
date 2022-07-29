package usrgrp

import (
	"Cinder/Chat/chatapi"

	assert "github.com/arl/assertgo"
)

// usrgrp 包不能导入 group 包，所以通过 IGroupMgr 接口执行 group 功能。
// 需要在初始化时 SetGroupMgr().
type IGroupMgr interface {
	LoginGroups(UserID, []GroupID)
	LogoutGroups(UserID, []GroupID)
	GetGroupOfflineMessageMap(map[GroupID]SequenceID) map[string][]*chatapi.ChatMessage
}

var groupMgr IGroupMgr

func SetGroupMgr(mgr IGroupMgr) {
	assert.True(mgr != nil)
	groupMgr = mgr
}
