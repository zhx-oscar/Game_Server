package loginout

import (
	"Cinder/Mail/rpcproc/handler/loginout/db"
	"Cinder/Mail/rpcproc/usersrv"
	"Cinder/Mail/rpcproc/userid"
	"fmt"
)

func Login(userID userid.UserID, peerSrvID string) error {
	// 写DB并异步广播新增 peerSrvID，希望在函数返回时所有节点已更新到
	if err := usersrv.InsertAndBroadcast(peerSrvID); err != nil {
		return fmt.Errorf("insert and broadcast: %w", err)
	}
	// DB插入 userID, peerSrvID
	if err := db.GetUsers(userID).Upsert(peerSrvID); err != nil {
		return fmt.Errorf("upsert users: %w", err)
	}

	return nil
}

func Logout(userID userid.UserID) error {
	if err := db.GetUsers(userID).Upsert(""); err != nil {
		return fmt.Errorf("upset users to reset srvID: %w", err)
	}
	return nil

}
