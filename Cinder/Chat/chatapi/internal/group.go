package internal

import (
	"Cinder/Chat/chatapi/types"
	"encoding/json"

	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

func CreateGroup(groupID string, members []string) error {
	// RPC 不支持 []string, 需要自己打包
	binMembers, err := json.Marshal(members)
	if err != nil {
		return errors.Wrap(err, "marshal error")
	}

	ret := Rpc("RPC_CreateGroup", groupID, binMembers)
	return GetStringError(ret)
}

func DeleteGroup(groupID string) error {
	ret := Rpc("RPC_DeleteGroup", groupID)
	return ret.Err
}

func MemberJoinGroup(memberID string, groupID string) error {
	ret := Rpc("RPC_JoinGroup", memberID, groupID)
	return ret.Err
}

func MemberLeaveGroup(memberID string, groupID string) error {
	ret := Rpc("RPC_LeaveGroup", memberID, groupID)
	return ret.Err
}

func GetGroupMembers(groupID string) ([]string, error) {
	ret := Rpc("RPC_GetGroupMembers", groupID)
	if ret.Err != nil {
		log.Debug("RPC_GetGroupMembers err: ", ret.Err)
		return nil, ret.Err
	}

	binMembers := ret.Ret[0].([]byte)
	members := []string{}
	err := json.Unmarshal(binMembers, &members)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal error")
	}
	return members, nil
}

func SendGroupMessage(userID string, groupID string, msgContent []byte) error {
	// log.Debugf("SendGroupMessage uesrID='%s' groupID='%s' msgContent='%v'", userID, groupID, msgContent)
	ret := Rpc("RPC_SendGroupMessage", userID, groupID, msgContent)
	return ret.Err
}

// 获取一批群的成员数，用于判断多个世界群的大小(繁忙程度)
// 返回map, 以群ID为键。返回 nil 表示出错。
func GetGroupMemberCounts(groupIDs []string) map[string]int {
	// RPC 不支持 []string, 需要自己打包
	binGroupIDs, err := json.Marshal(groupIDs)
	if err != nil {
		log.Errorf("json marshal error: %v", err)
		return nil
	}

	ret := Rpc("RPC_GetGroupMemberCounts", binGroupIDs)
	if ret.Err != nil {
		log.Errorf("RPC_GetGroupMemberCounts error: %v", ret.Err)
		return nil
	}

	bin := ret.Ret[0].([]byte)
	if len(bin) == 0 {
		log.Errorf("RPC_GetGroupMemberCounts returns nil")
		return nil
	}
	result := map[string]int{}
	if err := json.Unmarshal(bin, &result); err != nil {
		log.Errorf("Unmarshal error: %v", err)
		return nil
	}
	return result
}

// GetGroupHistoryMessages 获取群最近 count 条历史消息。
// 出错则返回nil
func GetGroupHistoryMessages(groupID string, count uint16) []*types.ChatMessage {
	ret := Rpc("RPC_GetGroupHistoryMessages", groupID, count)
	if ret.Err != nil {
		log.Errorf("RPC_GetGroupMemberCounts error: %v", ret.Err)
		return nil
	}

	bin := ret.Ret[0].([]byte)
	if len(bin) == 0 {
		log.Errorf("RPC_GetGroupMemberCounts returns nil")
		return nil
	}
	result := []*types.ChatMessage{}
	if err := json.Unmarshal(bin, &result); err != nil {
		log.Errorf("Unmarshal error: %v", err)
		return nil
	}
	return result
}
