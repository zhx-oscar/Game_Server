package internal

import (
	"Cinder/Chat/rpcproc/logic/bc"
	"Cinder/Chat/rpcproc/logic/usermgr"
)

// GroupBroadcastMessage 群内广播消息
// 向在线成员发送消息，包括发送者，并更新成员已收序号。
func GroupBroadcastMessage(groupID GroupID, groupOnlineMemberIDs map[UserID]bool,
	fromRoleID UserID, fromNick string, fromData []byte, msgContent []byte, sequenceID SequenceID) {
	bc.GroupBroadcastMessage(groupID, groupOnlineMemberIDs, fromRoleID, fromNick, fromData, msgContent)
	updateGroupMsgSeq(groupID, groupOnlineMemberIDs, sequenceID)
}

// updateGroupMsgSeq 更新已收序号
func updateGroupMsgSeq(groupID GroupID, groupOnlineMemberIDs map[UserID]bool, sequenceID SequenceID) {
	userMgr := usermgr.GetUserMgr()
	for userID, _ := range groupOnlineMemberIDs {
		user := userMgr.GetUser(userID)
		if user == nil {
			continue // 有可能并发，UserMgr中已下线，但Group中还未下线
		}

		// 更新已收序号, 保证单调增加，因为并发有可能错序
		user.GetUserGroupMgr().UpdateGroupMsgSeq(groupID, sequenceID)
	}
}
