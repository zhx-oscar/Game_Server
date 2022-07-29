package user

import (
	"Cinder/Chat/rpcproc/logic/user/internal/usrgrp"
)

// AddGroupToUsers 给玩家添加聊天群, 初始化已读消息序号。
// 玩家可能不在线。如果在线，则需更新内存对象。
// 本函数不改变Group数据。
func AddGroupToUsers(groupID GroupID, userIDs []UserID, seqID SequenceID) error {
	return usrgrp.AddGroupToUsers(groupID, userIDs, seqID)
}

// DeleteGroupFromUsers 删除玩家上的聊天群。
// 玩家可能不在线。如果在线，则需更新内存对象。
// 本函数不改变Group数据。
func DeleteGroupFromUsers(groupID GroupID, userIDs []UserID) error {
	return usrgrp.DeleteGroupFromUsers(groupID, userIDs)
}
