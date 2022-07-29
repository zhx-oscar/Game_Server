package internal

import (
	"Cinder/Chat/rpcproc/logic/group/internal/dbutil"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user"
	"Cinder/Chat/rpcproc/logic/usermgr"

	assert "github.com/arl/assertgo"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

type UserID = types.UserID

// 聊天群成员。
// 不必加锁，因为Group中已有锁。
type GroupMembers struct {
	groupID GroupID

	// 在线成员ID, 用于发消息时遍历跳过离线ID。
	onlineMembers map[UserID]bool
	// 所有成员, 用于返回成员列表，及KickAllMembers()
	members map[UserID]bool
}

// NewGroupMembers 创建对象
func NewGroupMembers(groupID GroupID) *GroupMembers {
	result := &GroupMembers{
		groupID: groupID,
	}
	result.reset()
	return result
}

// reset 重置成员
func (g *GroupMembers) reset() {
	g.onlineMembers = make(map[UserID]bool)
	g.members = make(map[UserID]bool)
}

// Login 设置上线状态。
func (g *GroupMembers) Login(userID UserID) {
	if g.HasMember(userID) {
		g.onlineMembers[userID] = true
	}
}

// Logout 设置下线状态。
func (g *GroupMembers) Logout(userID UserID) {
	delete(g.onlineMembers, userID)
}

// HasMember 判断是否成员
func (g *GroupMembers) HasMember(userID UserID) bool {
	_, ok := g.members[userID]
	return ok
}

// AddMembers 加成员入群组, 立即存DB.
// 允许添加离线成员。
// 因为成员可能未上线，所以没有 User 对象，需直接在 DB 中为 User 添加群.
// seqID 为当前群最大消息序号。
func (g *GroupMembers) AddMembers(members []UserID, seqID SequenceID) error {
	// log.Debugf("add members to group '%s': %v", g.groupID, members)
	toAddIDs := g.getNonMembers(members) // 忽略已有成员
	if len(toAddIDs) == 0 {
		// log.Debugf("all are already members")
		return nil
	}

	groupID := g.groupID
	groupUtil := dbutil.GroupUtil(groupID)
	err := groupUtil.InsertMembers(toAddIDs)
	if err != nil {
		return errors.Wrap(err, "db insert members")
	}

	// users 添加该聊天群, user不一定在线. 一次性写DB, 并添加当前序号
	if err := user.AddGroupToUsers(groupID, toAddIDs, seqID); err != nil {
		// 群数据中已有成员，但玩家数据中添加群失败
		// TODO: 是否可忽略不一致性？是否会在玩家保存群序号时自动修复？
		log.Errorf("failed to add group '%v' to users '%v': %s", groupID, toAddIDs, err)
	}

	userMgr := usermgr.GetUserMgr()
	assert.True(userMgr != nil)
	for _, userID := range toAddIDs {
		// 更新成员列表
		g.members[userID] = true
		// 需要判断在线状态
		if userMgr.IsUserOnline(userID) {
			g.onlineMembers[userID] = true
		}
	}
	return nil
}

// getNonMembers 获取非成员。
// 用于添加成员时忽略已有成员。
func (g *GroupMembers) getNonMembers(userIDs []UserID) []UserID {
	result := make([]UserID, 0, len(userIDs))
	for _, userID := range userIDs {
		if !g.HasMember(userID) {
			result = append(result, userID)
		}
	}
	return result
}

// DeleteMember 让成员离开群组, 立即存DB
func (g *GroupMembers) DeleteMember(userID UserID) {
	// log.Debugf("'%s' leave group '%s'", userID, g.groupID)
	if !g.HasMember(userID) {
		log.Debugf("DeleteMember: no user '%s' in group '%s'", userID, g.groupID)
		return
	}
	delete(g.onlineMembers, userID) // 更新在线列表
	delete(g.members, userID)

	groupUtil := dbutil.GroupUtil(g.groupID)
	groupUtil.DeleteMember(userID)

	user.DeleteGroupFromUsers(g.groupID, []UserID{userID})
}

// GetGroupMemberIDs 获取成员ID列表。
// 以复制方式获取。
func (g *GroupMembers) GetMemberIDs() []UserID {
	result := make([]UserID, 0, len(g.members))
	for userID, _ := range g.members {
		result = append(result, userID)
	}
	return result
}

func (g *GroupMembers) GetMemberCount() int {
	return len(g.members)
}

// Load 从DB加载群成员列表, 并更新在线列表
func (g *GroupMembers) Load() error {
	members, err := dbutil.GroupUtil(g.groupID).LoadMemberIDs()
	if err != nil {
		return errors.Wrap(err, "db load group members")
	}

	g.reset()
	for _, userID := range members {
		g.setMember(userID) // 设置为成员，并更新在线状态
	}
	return nil
}

// CopyOnlineMemberIDs 复制在线成员ID
func (g *GroupMembers) CopyOnlineMemberIDs() map[UserID]bool {
	result := make(map[UserID]bool, len(g.onlineMembers))
	for id, _ := range g.onlineMembers {
		result[id] = true
	}
	return result
}

// set 设置为成员并更新在线状态。
func (g *GroupMembers) setMember(userID UserID) {
	g.members[userID] = true

	// 检查并更新在线成员
	if usermgr.GetUserMgr().IsUserOnline(userID) {
		g.onlineMembers[userID] = true
	}
}
