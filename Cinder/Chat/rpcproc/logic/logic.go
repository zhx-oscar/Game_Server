// logic 包处理聊天服逻辑
package logic

import (
	"Cinder/Chat/rpcproc/logic/bc"
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/Chat/rpcproc/logic/group"
	"Cinder/Chat/rpcproc/logic/user"
	"Cinder/Chat/rpcproc/logic/usermgr"
	"fmt"
)

func init() {
	user.SetUserMgr(usermgr.GetUserMgr())
	// 包变量 group._GroupMgr 应该在 init() 之前初始化
	user.SetGroupMgr(group.GetGroupMgr())

	bc.SetUserMgr(usermgr.GetUserMgr())
	bc.SetGroupMgr(group.GetGroupMgr())
}

// EnsureMongoDBIndexes 创建 mongodb 索引
func EnsureMongoDBIndexes() error {
	if err := db.EnsureIndexes(); err != nil {
		return fmt.Errorf("ensure db indexes: %w", err)
	}
	return nil
}

// SaveAllToDBOnExit 用于进程退出时保存信息到DB
func SaveAllToDBOnExit() {
	group.GetGroupMgr().SaveAllToDBOnExit()
	usermgr.GetUserMgr().SaveAllToDBOnExit()
}
