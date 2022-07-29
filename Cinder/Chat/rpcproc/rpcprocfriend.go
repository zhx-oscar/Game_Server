package rpcproc

import (
	"Cinder/Chat/rpcproc/logic/usermgr"
	"encoding/json"
	"fmt"

	log "github.com/cihub/seelog"
)

type _RPCProcFriend struct {
}

// RPC_AddFriendReq actorID 请求加 friendID 为好友
func (r *_RPCProcFriend) RPC_AddFriendReq(actorID, friendID string, reqInfo []byte) (errStr string) {
	// log.Debugf("RPC_AddFriendReq actorID=%v friendID=%v", actorID, friendID)
	user := usermgr.GetUserMgr().GetUser(UserID(actorID))
	if user == nil {
		return fmt.Sprintf("requester '%s' is not online", actorID)
	}
	if err := user.GetFriendMgr().SendRequest(UserID(friendID), reqInfo); err != nil {
		return fmt.Sprintf("failed to add friend request: %s", err)
	}
	return ""
}

// RPC_ReplyAddFriendReq responderID 响应 fromID 的加好友请求
func (r *_RPCProcFriend) RPC_ReplyAddFriendReq(responderID, fromID string, ok bool) (errStr string) {
	// log.Debugf("RPC_ReplyAddFriendReq responderID=%v fromID=%v ok=%v", responderID, fromID, ok)
	user := usermgr.GetUserMgr().GetUser(UserID(responderID))
	if user == nil {
		return fmt.Sprintf("friend responder '%s' is not online", responderID)
	}
	if err := user.GetFriendMgr().SendResponse(UserID(fromID), ok); err != nil {
		return fmt.Sprintf("failed to send add friend response: %s", err)
	}
	return ""
}

// RPC_DeleteFriend actorID 删除好友 friendID
func (r *_RPCProcFriend) RPC_DeleteFriend(actorID, friendID string) {
	// log.Debugf("RPC_DeleteFriend actorID=%v friendID=%v", actorID, friendID)
	user := usermgr.GetUserMgr().GetUser(UserID(actorID))
	if user == nil {
		log.Errorf("'%s' is not online", actorID)
		return
	}
	if err := user.GetFriendMgr().DeleteFriendActive(UserID(friendID)); err != nil {
		log.Errorf("failed to delete friend: %v", err)
	}
}

// RPC_GetFriendList 取 userID 的好友列表。
// userID 必须在线。
func (r *_RPCProcFriend) RPC_GetFriendList(userID string) []byte {
	// log.Debugf("RPC_GetFriendList userID=%v", userID)
	user := usermgr.GetUserMgr().GetUser(UserID(userID))
	if user == nil {
		log.Errorf("'%s' is not online", userID)
		return nil
	}
	infos, err := user.GetFriendMgr().GetFriendList()
	if err != nil {
		log.Errorf("failed to get friend list: %v", err)
		return nil
	}

	bin, errJson := json.Marshal(infos)
	if errJson != nil {
		log.Errorf("json marshal error: %v", err)
		return nil
	}

	return bin
}
