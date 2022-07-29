package usrgrp

import (
	"Cinder/Chat/chatapi"
	"fmt"
)

// UserGroupMgr 管理单个用户的所有群
type UserGroupMgr struct {
	// 因为 GroupToSeq 是协程安全的，其他成员不变，所以不需要加锁

	userID UserID

	// 初始化时的各群已读序号，用于收集离线消息，初始化后只读。更新数据在下面的groupToSeq。
	initGroupToSeq map[GroupID]SequenceID

	groupToSeq *GroupToSeq // 记录所有聊天群，及群内已读序号
}

func NewUserGroupMgr(userID UserID) *UserGroupMgr {
	return &UserGroupMgr{
		userID: userID,

		initGroupToSeq: make(map[GroupID]SequenceID),
		groupToSeq:     NewGroupToSeq(userID),
	}
}

// LoadGroups 从DB加载用户的所有群.
// 此时User还没加入UserMgr.
func (u *UserGroupMgr) LoadGroups() error {
	// log.Debug("UserGroupMgr LoadGroups()")
	if err := u.groupToSeq.Load(); err != nil {
		return fmt.Errorf("load group to sequence ID map error: %w", err)
	}
	u.initGroupToSeq = u.groupToSeq.CopyGroupToSeq()

	// 设置所有群中在线状态, 会触发群加载
	groupIDs := u.groupToSeq.GetGroupIDs()
	groupMgr.LoginGroups(u.userID, groupIDs)
	return nil
}

// GetGroupOfflineMessageMap 收集群聊离线消息.
// 群已预先加载，所以此时不会有DB加载。
func (u *UserGroupMgr) GetGroupOfflineMessageMap() map[string][]*chatapi.ChatMessage {
	return groupMgr.GetGroupOfflineMessageMap(u.initGroupToSeq)
}

// Save 保存数据到DB.
func (u *UserGroupMgr) Save() {
	// 只有群已读序号需写DB, 其他都是变化时就写DB的。
	u.groupToSeq.Save()
}

// AddGroup 添加群。
// 仅在内存中添加.
// DB加群删群不是User的功能，而是Group的功能，因为User可能不在线。
func (u *UserGroupMgr) AddGroup(groupID GroupID, seqID SequenceID) {
	u.groupToSeq.AddGroupIfNot(groupID, seqID)
}

// DeleteGroup 删除群。
// 仅在内存中删除。
// DB加群删群不是User的功能，而是Group的功能，因为User可能不在线。
func (u *UserGroupMgr) DeleteGroup(groupID GroupID) {
	u.groupToSeq.DeleteGroup(groupID)
}

// LogoutAllGroups 通知所有群下线.
func (u *UserGroupMgr) LogoutAllGroups() {
	groupIDs := u.groupToSeq.GetGroupIDs()

	// 通知所有群下线
	groupMgr.LogoutGroups(u.userID, groupIDs)
}

// UpdateGroupMsgSeq 更新聊天群中已收消息数，保证单调增加。
// 因为群异步广播消息，所以有可能跳过某些序号。
func (u *UserGroupMgr) UpdateGroupMsgSeq(groupID GroupID, sequenceID SequenceID) {
	u.groupToSeq.Update(groupID, sequenceID)
}

func (u *UserGroupMgr) GetGroupIDs() []GroupID {
	return u.groupToSeq.GetGroupIDs()
}
