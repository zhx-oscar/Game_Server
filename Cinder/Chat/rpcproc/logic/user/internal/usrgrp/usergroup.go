package usrgrp

import (
	"Cinder/Chat/rpcproc/logic/user/internal/usrgrp/dbutil"

	assert "github.com/arl/assertgo"
	"github.com/pkg/errors"
)

// AddGroupToUsers 给玩家添加聊天群, 初始化已读消息序号。
// 玩家可能不在线。如果在线，则需更新内存对象。
// 本函数不改变Group数据。
func AddGroupToUsers(groupID GroupID, userIDs []UserID, seqID SequenceID) error {
	db := dbutil.UsersGroupsUtil(groupID)
	if err := db.AddGroupToUsers(userIDs, seqID); err != nil {
		return errors.Wrap(err, "db add group to users")
	}

	// 如果 User 在线，则还需更新 User
	assert.True(userMgr != nil)
	for _, userID := range userIDs {
		if userGroupMgr := userMgr.GetUserGroupMgr(userID); userGroupMgr != nil {
			userGroupMgr.AddGroup(groupID, seqID)
		}
	}
	return nil
}

// DeleteGroupFromUsers 删除玩家上的聊天群。
// 玩家可能不在线。如果在线，则需更新内存对象。
// 本函数不改变Group数据。
func DeleteGroupFromUsers(groupID GroupID, userIDs []UserID) error {
	db := dbutil.UsersGroupsUtil(groupID)
	if err := db.DeleteGroupFromUsers(userIDs); err != nil {
		return errors.Wrap(err, "db delete group from users")
	}

	for _, userID := range userIDs {
		if userGroupMgr := userMgr.GetUserGroupMgr(userID); userGroupMgr != nil {
			userGroupMgr.DeleteGroup(groupID)
		}
	}
	return nil
}
