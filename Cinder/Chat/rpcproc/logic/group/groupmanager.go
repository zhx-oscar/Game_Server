package group

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/group/internal/dbutil"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user"
	"sync"

	assert "github.com/arl/assertgo"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

type GroupManager struct {
	groupMap *sync.Map
}

var _GroupMgr = &GroupManager{
	groupMap: &sync.Map{},
}

func GetGroupMgr() *GroupManager {
	return _GroupMgr
}

// SaveAllToDBOnExit 保存所有数据到DB
func (mgr *GroupManager) SaveAllToDBOnExit() {
	log.Info("group manager save all to DB on exit...")
	mgr.groupMap.Range(func(_, group interface{}) bool {
		group.(*Group).Save()
		return true
	})
	log.Info("group manager done")
}

// getGroup 仅在内存中查找群，不会DB加载。
// 一般应该用 GetOrLoadGroup().
func (mgr *GroupManager) getGroup(groupID GroupID) *Group {
	res, ok := mgr.groupMap.Load(groupID)
	if ok {
		return res.(*Group)
	}
	return nil
}

// GetOrLoadGroup 获取群对象，如果还未加载则先加载。
// 因为DB中不区分群成员为空和群不存在，所以加载群可能会加载到一个空群，即认为群总是存在的，只是成员为空。
// 返回 error 为空时，*Group 返回必定非空。
func (mgr *GroupManager) GetOrLoadGroup(groupID GroupID) (types.IGroup, error) {
	res := mgr.getGroup(groupID)
	if res != nil {
		return res, nil
	}

	group, err := LoadGroup(groupID)
	if err != nil {
		return nil, errors.Wrap(err, "LoadGroup")
	}
	assert.True(group != nil)
	mgr.groupMap.Store(groupID, group)
	return group, nil
}

// InsertMembersToGroup 创建群组并添加成员，如果群已存在，则是成员合并。
// 群将先加载, 然后添加成员。
func (mgr *GroupManager) InsertMembersToGroup(groupID GroupID, memberList []UserID) error {
	// log.Debug("Create Group: ", groupID, memberList)
	group, err := mgr.GetOrLoadGroup(groupID)
	if err != nil {
		return errors.Wrap(err, "GetORLoadGroup")
	}
	assert.True(group != nil)
	group.AddMembers(memberList)
	return nil
}

// DeleteGroup 删除群组
func (mgr *GroupManager) DeleteGroup(groupID GroupID) error {
	// log.Debug("Delete Group: ", groupID)
	group, err := mgr.GetOrLoadGroup(groupID)
	if err != nil {
		return errors.Wrap(err, "GetOrLoadGroup")
	}
	assert.True(group != nil)

	mgr.groupMap.Delete(groupID)
	groupUtil := dbutil.GroupUtil(groupID)
	if err := groupUtil.DeleteGroup(); err != nil {
		return errors.Wrap(err, "db delete group")
	}

	// 未删除群消息，下次创建同名群消息序号会顺沿下去，只是新加成员看不到旧消息, 新成员开始序号为当前序号.

	// 从所有成员的 group 列表中删除
	members := group.GetGroupMemberIDs()
	if err := user.DeleteGroupFromUsers(groupID, members); err != nil {
		// TODO: 出错时数据不一致, 仍按删除成功处理
		log.Errorf("failed to delete group '%v' from users '%v': %v", groupID, members, err)
	}
	return nil
}

// LogoutGroups 设置群中为离线状态
func (mgr *GroupManager) LogoutGroups(userID UserID, groupIDs []GroupID) {
	// log.Debugf("GroupManager LogoutGroups: userID=%s groupIDs=%v", userID, groupIDs)
	for _, groupID := range groupIDs {
		group := mgr.getGroup(groupID) // 不要从DB加载
		if group == nil {
			log.Warnf("group '%s' not found when logout '%s'", groupID, userID) // 群应该在User加载时已加载
			continue
		}
		group.LogoutMember(userID)
	}
}

// LoginGroups 设置群中的在线状态，会触发群加载
func (mgr *GroupManager) LoginGroups(userID UserID, groupIDs []GroupID) {
	// log.Debugf("GroupManager LoginGroups: userID=%s groupIDs=%v", userID, groupIDs)
	for _, groupID := range groupIDs {
		group, err := mgr.GetOrLoadGroup(groupID) // 可能从DB加载
		if err != nil {
			log.Warnf("ignore GetOrLoadGroup error: %v", err)
		}
		if group != nil {
			group.LoginMember(userID)
		}
	}
}

// GetGroupOfflineMessageMap 获取群离线消息。
func (mgr *GroupManager) GetGroupOfflineMessageMap(groupToSeq map[GroupID]SequenceID) map[string][]*chatapi.ChatMessage {
	// log.Debugf("GroupManager GetGroupOfflineMessageMap: groupToSeq=%v", groupToSeq)
	result := make(map[string][]*chatapi.ChatMessage)
	for groupID, seqID := range groupToSeq {
		group := mgr.getGroup(groupID) // 不要从DB加载
		if group == nil {
			log.Warnf("group '%s' not found when getting offline message", groupID) // 群应该在User加载时已加载
			continue
		}
		result[string(groupID)] = group.GetOfflineMessagesAfter(seqID)
	}
	return result
}

// GetGroupMemberCounts 返回一指群的成员个数
func (mgr *GroupManager) GetGroupMemberCounts(groupIDs []GroupID) map[GroupID]int {
	result := make(map[GroupID]int)
	for _, groupID := range groupIDs {
		group, err := mgr.GetOrLoadGroup(groupID)
		if err != nil {
			continue
		}
		result[groupID] = group.GetGroupMemberCount()
	}
	return result
}
