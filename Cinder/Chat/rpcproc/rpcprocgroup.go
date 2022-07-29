package rpcproc

import (
	"Cinder/Chat/rpcproc/logic/group"
	"Cinder/Chat/rpcproc/logic/types"
	"encoding/json"

	assert "github.com/arl/assertgo"
	log "github.com/cihub/seelog"
)

type _RPCProcGroup struct {
}

type GroupID = types.GroupID

// RPC_CreateGroup 处理创建聊天群的RPC.
// 如果群已存在，则合并成员。
// binMembers 为 []string json 打包，是成员列表。
func (r *_RPCProcGroup) RPC_CreateGroup(sGroupID string, binMembers []byte) (errStr string) {
	groupID := GroupID(sGroupID)
	members := []UserID{}
	if err := json.Unmarshal(binMembers, &members); err != nil {
		log.Debug("Unmarshal error: ", err)
		return "unmarshal error: " + err.Error()
	}

	// log.Debugf("RPC_CreateGroup groupID=%s members=%v", groupID, members)
	if err := group.GetGroupMgr().InsertMembersToGroup(groupID, members); err != nil {
		log.Debug("InsertMembersToGroup error")
		return "InsertMembersToGroup error: " + err.Error()
	}
	return ""
}

func (r *_RPCProcGroup) RPC_DeleteGroup(groupID string) {
	// log.Debugf("RPC_DeleteGroup groupID=%s", groupID)
	err := group.GetGroupMgr().DeleteGroup(GroupID(groupID))
	if err != nil {
		log.Errorf("failed to delete group '%s': %v", groupID, err)
	}
}

func (r *_RPCProcGroup) RPC_JoinGroup(sRoleID string, sGroupID string) {
	// log.Debugf("RPC_JoinGroup sRoleID=%s sGroupID=%s", sRoleID, sGroupID)
	roleID := UserID(sRoleID)
	groupID := GroupID(sGroupID)
	mgr := group.GetGroupMgr()
	assert.True(mgr != nil)
	group, err := mgr.GetOrLoadGroup(groupID)
	if err != nil {
		log.Errorf("get or load group %s failed: %v", groupID, err)
		return
	}
	assert.True(group != nil)
	if err := group.AddMember(roleID); err != nil {
		log.Errorf("join group error: %v", err)
	}
}

func (r *_RPCProcGroup) RPC_LeaveGroup(sRoleID string, sGroupID string) {
	// log.Debugf("RPC_LeaveGroup sRoleID=%s sGroupID=%s", sRoleID, sGroupID)
	roleID := UserID(sRoleID)
	groupID := GroupID(sGroupID)
	group, err := group.GetGroupMgr().GetOrLoadGroup(groupID)
	if err != nil {
		log.Errorf("get or load group '%s' error: %v", groupID, err)
		return
	}
	assert.True(group != nil)
	group.KickMember(roleID)
}

// RPC_GetGroupMembers 获取聊天群成员ID列表。
// 返回 []byte 是成员列表 []string json 打包。
func (r *_RPCProcGroup) RPC_GetGroupMembers(groupID string) []byte {
	// log.Debugf("RPC_GetGroupMembers groupID=%s", groupID)
	group, err := group.GetGroupMgr().GetOrLoadGroup(GroupID(groupID))
	if err != nil {
		log.Errorf("can not get or load group '%s': %v", groupID, err)
		return nil
	}
	assert.True(group != nil)

	ids := group.GetGroupMemberIDs()
	bin, err2 := json.Marshal(ids)
	if err2 != nil {
		log.Error("json marshal error: ", err2)
		return nil
	}
	return bin
}

// RPC_SendGroupMessage 发送群聊消息。
func (r *_RPCProcGroup) RPC_SendGroupMessage(sFromRoleID string, sGroupID string, msgContent []byte) {
	// log.Debugf("RPC_SendGroupMessage sFromRoleID=%s sGroupID=%s msgContent=%v", sFromRoleID, sGroupID, msgContent)
	fromRoleID := UserID(sFromRoleID)
	groupID := GroupID(sGroupID)
	// log.Debugf("RPC_SendGroupMessage '%s'->group(ID='%s') msgContent='%v'", fromRoleID, groupID, msgContent)
	group, err := group.GetGroupMgr().GetOrLoadGroup(groupID)
	if err != nil {
		log.Errorf("get or load group '%s' error: %v", groupID, err)
		return
	}
	assert.True(group != nil)
	group.SendGroupMessage(fromRoleID, msgContent)
}

// RPC_GetGroupMemberCounts 获取一批群的成员数
func (r *_RPCProcGroup) RPC_GetGroupMemberCounts(binGroupIDs []byte) []byte {
	// log.Debugf("RPC_GetGroupMemberCounts binGroupIDs=%v", binGroupIDs)
	groupIDs := []GroupID{}
	if err := json.Unmarshal(binGroupIDs, &groupIDs); err != nil {
		log.Errorf("Unmarshal error: %v", err)
		return nil
	}

	var cnts map[GroupID]int = group.GetGroupMgr().GetGroupMemberCounts(groupIDs)
	bin, err2 := json.Marshal(cnts)
	if err2 != nil {
		log.Errorf("json marshal error: %v", err2)
		return nil
	}
	return bin
}

// RPC_GetGroupHistoryMessages 获取群的最近几条历史消息
func (r *_RPCProcGroup) RPC_GetGroupHistoryMessages(groupID string, count uint16) []byte {
	// log.Debugf("RPC_GetGroupHistoryMessages groupID=%v, count=%v", groupID, count)
	group, err := group.GetGroupMgr().GetOrLoadGroup(GroupID(groupID))
	if err != nil {
		log.Errorf("get or load group error: %v", err)
		return nil
	}
	msgs := group.GetHistoryMessages(count)
	bin, err2 := json.Marshal(msgs)
	if err2 != nil {
		log.Errorf("json marshal error: %v", err2)
		return nil
	}
	return bin
}
