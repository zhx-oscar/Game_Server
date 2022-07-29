package rpcproc

import (
	"Cinder/Chat/rpcproc/logic/bc"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/usermgr"
	"encoding/json"
	"fmt"

	log "github.com/cihub/seelog"
)

type _RPCProcUser struct {
}

type UserID = types.UserID

// RPC_LoginWithNickData 玩家登录, 附带昵称和数据
// peerSrvID 记录玩家实例所在服，用于向该服推送聊天消息。
func (r *_RPCProcUser) RPC_LoginWithNickData(roleID string, nick string, activeData []byte, peerSrvID string) (errStr string) {
	// log.Debugf("RPC_LoginWithNickData: roleID=%s nick=%s activeData=%v srvID=%s", roleID, nick, activeData, peerSrvID)
	if u, err := usermgr.GetUserMgr().LoginUser(UserID(roleID), peerSrvID); err != nil {
		return fmt.Sprintf("LoginUser: %s", err)
	} else if _, err := u.SetNick(nick); err != nil {
		return fmt.Sprintf("SetNick: %s", err)
	} else if _, err := u.SetActiveData(activeData); err != nil {
		return fmt.Sprintf("SetActiveData: %s", err)
	}
	bc.BroadcastUserLogin(UserID(roleID))
	// TODO: 应该在登录时就自动发送离线聊天，不然 Login 和 GetOfflineMessage 之间可能会有实时消息
	return ""
}

// RPC_Login 玩家登录
// peerSrvID 记录玩家实例所在服，用于向该服推送聊天消息。
func (r *_RPCProcUser) RPC_Login(roleID string, peerSrvID string) (errStr string) {
	// log.Debugf("RPC_Login: roleID=%s srvID=%s", roleID, peerSrvID)
	_, err := usermgr.GetUserMgr().LoginUser(UserID(roleID), peerSrvID)
	if err != nil {
		return fmt.Sprintf("LoginUser: %v", err)
	}
	bc.BroadcastUserLogin(UserID(roleID))
	// TODO: 应该在登录时就自动发送离线聊天，不然 Login 和 GetOfflineMessage 之间可能会有实时消息
	return ""
}

// RPC_Logout 玩家登出
func (r *_RPCProcUser) RPC_Logout(roleID string) {
	// log.Debugf("RPC_Logout: %s", roleID)
	usermgr.GetUserMgr().RemoveUser(UserID(roleID))
}

// RPC_SetNick 设置昵称
func (r *_RPCProcUser) RPC_SetNick(roleID string, nick string) (errStr string) {
	if u := getUser(roleID); u == nil {
		return "can not find user: " + roleID
	} else if changed, err := u.SetNick(nick); err != nil {
		return err.Error()
	} else if changed {
		bc.BroadcastUserNick(UserID(roleID), nick)
	}
	return ""
}

// RPC_GetNick 获取昵称
func (r *_RPCProcUser) RPC_GetNick(roleID string) (errStr string, nick string) {
	u := getUser(roleID)
	if u == nil {
		return "can not find user: " + roleID, ""
	}
	return "", u.GetNick()
}

// RPC_SetActiveData 设置主动数据
func (r *_RPCProcUser) RPC_SetActiveData(roleID string, activeData []byte) (errStr string) {
	if u := getUser(roleID); u == nil {
		return "can not find user: " + roleID
	} else if changed, err := u.SetActiveData(activeData); err != nil {
		return err.Error()
	} else if changed {
		bc.BroadcastUserData(UserID(roleID), activeData)
	}
	return ""
}

// RPC_GetActiveData 获取主动数据
func (r *_RPCProcUser) RPC_GetActiveData(roleID string) (errStr string, activeData []byte) {
	u := getUser(roleID)
	if u == nil {
		return "can not find user: " + roleID, nil
	}
	return "", u.GetActiveData()
}

// RPC_SetPassiveData 设置被动数据
func (r *_RPCProcUser) RPC_SetPassiveData(roleID string, passiveData []byte) (errStr string) {
	if u := getUser(roleID); u == nil {
		return "can not find user: " + roleID
	} else if err := u.SetPassiveData(passiveData); err != nil {
		return err.Error()
	}
	return ""
}

// RPC_GetPassiveData 获取被动数据
func (r *_RPCProcUser) RPC_GetPassiveData(roleID string) (errStr string, passiveData []byte) {
	u := getUser(roleID)
	if u == nil {
		return "can not find user: " + roleID, nil
	}
	return "", u.GetPassiveData()
}

// RPC_SendMessage 发送一对一聊天消息。
func (r *_RPCProcUser) RPC_SendMessage(fromRoleID string, targetRoleID string, msgContent []byte) {
	// log.Debugf("RPC_SendMessage '%s'->'%s' msgContent='%v'", fromRoleID, targetRoleID, msgContent)
	user := getUser(fromRoleID)
	if user == nil {
		log.Warnf("message sender '%s' is not online", fromRoleID)
		return
	}
	user.SendChatMessage(UserID(targetRoleID), msgContent)
}

// RPC_GetOfflineMessage 获取历史消息。
func (r *_RPCProcUser) RPC_GetOfflineMessage(sRoleID string) []byte {
	// log.Debugf("RPC_GetOfflineMessage roleID = '%s'", roleID)
	user := getUser(sRoleID)
	if user == nil {
		log.Warnf("user '%s' is not online when RPC_GetOfflineMessage", sRoleID)
		return nil
	}
	msg, errLoad := user.LoadOfflineMessage()
	if errLoad != nil {
		log.Errorf("load offline message error: %v", errLoad)
		return nil
	}
	buf, errMarshal := json.Marshal(msg)
	if errMarshal != nil {
		log.Errorf("marshal error: %v", errMarshal)
		return nil
	}
	return buf
}

func getUser(roleID string) types.IUser {
	return usermgr.GetUserMgr().GetUser(UserID(roleID))
}
