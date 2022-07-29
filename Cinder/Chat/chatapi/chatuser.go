package chatapi

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Core"
	"Cinder/Chat/chatapi/internal"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/cihub/seelog"
)

// 必须由 Login(id string) 创建
type _ChatUser struct {
	id string
}

func newChatUser(id string) *_ChatUser {
	return &_ChatUser{
		id: id,
	}
}

// 必须在 Core.Inst 初始化之后使用。
func (user *_ChatUser) Login() error {
	ret := internal.Rpc("RPC_Login", user.id, Core.Inst.GetServiceID())
	return internal.GetStringErrorWithHint(ret, "RPC_Login")
}

// 必须在 Core.Inst 初始化之后使用。
func (user *_ChatUser) LoginWithNickData(nick string, activeData []byte) error {
	ret := internal.Rpc("RPC_LoginWithNickData", user.id, nick, activeData, Core.Inst.GetServiceID())
	return internal.GetStringErrorWithHint(ret, "login")
}

// Logout 退出，之后就不会有聊天服的回调了
func (user *_ChatUser) Logout() error {
	if ret := internal.Rpc("RPC_Logout", user.id); ret.Err != nil {
		return fmt.Errorf("logout failed: %w", ret.Err)
	}
	return nil
}

func (user *_ChatUser) SetNick(nick string) error {
	ret := internal.Rpc("RPC_SetNick", user.id, nick)
	return internal.GetStringErrorWithHint(ret, "RPC_SetNick")
}

func (user *_ChatUser) GetNick() (string, error) {
	ret := internal.Rpc("RPC_GetNick", user.id)
	if ret.Err != nil {
		return "", ret.Err
	}
	if errStr := ret.Ret[0].(string); errStr != "" {
		return "", fmt.Errorf("RPC_GetNick returns error: %s", errStr)
	}
	return ret.Ret[1].(string), nil
}

func (user *_ChatUser) SetActiveData(data []byte) error {
	return rpcSetData("RPC_SetActiveData", user.id, data)
}

func (user *_ChatUser) GetActiveData() ([]byte, error) {
	return rpcGetData("RPC_GetActiveData", user.id)
}

func (user *_ChatUser) SetPassiveData(data []byte) error {
	return rpcSetData("RPC_SetPassiveData", user.id, data)
}

func (user *_ChatUser) GetPassiveData() ([]byte, error) {
	return rpcGetData("RPC_GetPassiveData", user.id)
}

func rpcSetData(rpcMethod string, userID string, data []byte) error {
	ret := internal.Rpc(rpcMethod, userID, data)
	return internal.GetStringErrorWithHint(ret, rpcMethod)
}

func rpcGetData(rpcMethod string, userID string) ([]byte, error) {
	ret := internal.Rpc(rpcMethod, userID)
	if ret.Err != nil {
		return nil, ret.Err
	}
	if errStr := ret.Ret[0].(string); errStr != "" {
		return nil, fmt.Errorf("%s returns error: %s", rpcMethod, errStr)
	}
	return ret.Ret[1].([]byte), nil
}

func (user *_ChatUser) GetOfflineMessage() *OfflineMessage {
	ret := internal.Rpc("RPC_GetOfflineMessage", user.id)
	if ret.Err != nil {
		log.Error("GetOfflineMessage err: ", ret.Err)
		return nil
	}

	s := ret.Ret[0].([]byte)
	var result OfflineMessage
	if err := json.Unmarshal(s, &result); err != nil {
		log.Errorf("unmarshal error: %v", err)
	}
	return &result
}

func (user *_ChatUser) SendMessage(target string, msgContent []byte) error {
	// log.Debugf("SendMessage target='%s' msgContent='%v'", target, msgContent)
	ret := internal.Rpc("RPC_SendMessage", user.id, target, msgContent)
	return ret.Err
}

func (user *_ChatUser) CreateGroup(groupID string, members []string) error {
	return internal.CreateGroup(groupID, members)
}

func (user *_ChatUser) DeleteGroup(groupID string) error {
	return internal.DeleteGroup(groupID)
}

func (user *_ChatUser) JoinGroup(groupID string) error {
	ret := internal.Rpc("RPC_JoinGroup", user.id, groupID)
	return ret.Err
}

func (user *_ChatUser) LeaveGroup(groupID string) error {
	return internal.MemberLeaveGroup(user.id, groupID)
}

func (user *_ChatUser) KickFromGroup(groupID string, member string) error {
	// 同 LeaveGroup()
	return internal.MemberLeaveGroup(member, groupID)
}

func (user *_ChatUser) GetGroupMembers(groupID string) ([]string, error) {
	return internal.GetGroupMembers(groupID)
}

func (user *_ChatUser) SendGroupMessage(groupID string, msgContent []byte) error {
	// log.Debugf("SendGroupMessage groupID='%s' msgContent='%v'", groupID, msgContent)
	return internal.SendGroupMessage(user.id, groupID, msgContent)
}

func (user *_ChatUser) FollowFriendReq(friendID string) error {
	ret := internal.Rpc("RPC_FollowFriendReq", user.id, friendID)
	return ret.Err
}

func (user *_ChatUser) UnFollowFriendReq(friendID string) error {
	ret := internal.Rpc("RPC_UnFollowFriendReq", user.id, friendID)
	return ret.Err
}

func (user *_ChatUser) GetFollowingList() []FriendInfo {
	return internal.GetFollowingList(user.id)
}

func (user *_ChatUser) GetFollowerList() []FollowerInfo {
	return internal.GetFollowerList(user.id)
}

func (user *_ChatUser) AddFriendReq(friendID string, reqInfo []byte) error {
	ret := internal.Rpc("RPC_AddFriendReq", user.id, friendID, reqInfo)
	return getFriendRpcError(ret)
}

// 作废，改用 ReplyAddFriendReq()
func (user *_ChatUser) ApplyAddFriendReq(fromID string, ok bool) error {
	return user.ReplyAddFriendReq(fromID, ok)
}

func (user *_ChatUser) ReplyAddFriendReq(fromID string, ok bool) error {
	ret := internal.Rpc("RPC_ReplyAddFriendReq", user.id, fromID, ok)
	return getFriendRpcError(ret)
}

// getFriendRpcError 将好友RPC返回的某些错误转成预定义错误
func getFriendRpcError(ret CRpc.RpcRet) error {
	err := internal.GetStringError(ret)
	if err == nil {
		return nil
	}

	errStr := err.Error()
	if strings.Contains(errStr, SelfReachedMaxFriendCount) {
		return ErrSelfReachedMaxFriendCount
	}
	if strings.Contains(errStr, PeerReachedMaxFriendCount) {
		return ErrPeerReachedMaxFriendCount
	}
	if strings.Contains(errStr, PeerReachedMaxAddFriendReq) {
		return ErrPeerReachedMaxAddFriendReq
	}
	return err
}

func (user *_ChatUser) AddFriendToBlacklist(friendID string) error {
	ret := internal.Rpc("RPC_AddFriendToBlacklist", user.id, friendID)
	return ret.Err
}

func (user *_ChatUser) RemoveFriendFromBacklist(friendID string) error {
	return user.RemoveFriendFromBlacklist(friendID)
}

func (user *_ChatUser) RemoveFriendFromBlacklist(friendID string) error {
	ret := internal.Rpc("RPC_RemoveFriendFromBlacklist", user.id, friendID)
	return ret.Err
}

func (user *_ChatUser) DeleteFriend(friendID string) error {
	ret := internal.Rpc("RPC_DeleteFriend", user.id, friendID)
	return ret.Err
}

func (user *_ChatUser) GetFriendList() []FriendInfo {
	return internal.GetFriendList(user.id)
}

func (user *_ChatUser) GetFriendBlacklist() []FriendInfo {
	return internal.GetFriendBlacklist(user.id)
}
