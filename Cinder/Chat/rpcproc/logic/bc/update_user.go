package bc

import (
	"Cinder/Chat/chatapi/types"
	"Cinder/Chat/rpcproc/logic/rpc"
	ltypes "Cinder/Chat/rpcproc/logic/types"
	"encoding/json"

	log "github.com/cihub/seelog"
)

func BroadcastUserNick(userID ltypes.UserID, newNick string) {
	ids := getBcUserIDs(userID)
	broadcastUpdateUserNick(userID, newNick, ids)
}

func BroadcastUserData(userID ltypes.UserID, activeData []byte) {
	ids := getBcUserIDs(userID)
	broadcastUpdateUserData(userID, activeData, ids)
}

func BroadcastUserLogin(userID ltypes.UserID) {
	u := userMgr.GetUser(userID)
	if u == nil {
		log.Errorf("can not find user when BroadcastUserLogin") // 应该先添加 userMgr
		return
	}
	nick, data := u.GetNickAndData()
	msg := types.UpdateUserMsg{
		UserID: string(userID),
		Type:   types.UUT_LOGIN_WITH_NICK_AND_DATA,
		Nick:   nick,
		Data:   data,
	}
	ids := getBcUserIDs(userID)
	goBroadcastUpdateUserMsg(msg, ids)
}

func BroadcastUserLogout(userID ltypes.UserID) {
	if userMgr.GetUser(userID) == nil {
		log.Errorf("can not find user when BroadcastUserLogout") // 不要删 userMgr 中的 User
		return
	}

	msg := types.UpdateUserMsg{
		UserID: string(userID),
		Type:   types.UUT_LOGOUT,
	}
	ids := getBcUserIDs(userID)
	goBroadcastUpdateUserMsg(msg, ids)
}

func broadcastUpdateUserNick(userID ltypes.UserID, nick string, ids userIDs) {
	msg := types.UpdateUserMsg{
		UserID: string(userID),
		Type:   types.UUT_NICK,
		Nick:   nick,
	}
	goBroadcastUpdateUserMsg(msg, ids)
}

func broadcastUpdateUserData(userID ltypes.UserID, data []byte, ids userIDs) {
	msg := types.UpdateUserMsg{
		UserID: string(userID),
		Type:   types.UUT_DATA,
		Data:   data,
	}
	goBroadcastUpdateUserMsg(msg, ids)
}

func goBroadcastUpdateUserMsg(msg types.UpdateUserMsg, ids userIDs) {
	if len(ids) != 0 {
		go broadcastUpdateUserMsg(msg, ids)
	}
}

func broadcastUpdateUserMsg(msg types.UpdateUserMsg, ids userIDs) {
	srvIDToTargets := getSrvIDToTargetsMap(ids)
	for srvID, targets := range srvIDToTargets {
		msg.Targets = targets
		sendUpdateUserMsg(msg, srvID)
	}
}

func sendUpdateUserMsg(msg types.UpdateUserMsg, srvID string) {
	bin, binErr := json.Marshal(msg)
	if binErr != nil {
		log.Errorf("marshal UpdateUserMsg error: %v", binErr)
		return
	}

	ret := rpc.Rpc(srvID, "RPC_ChatUpdateUser", bin)
	if ret.Err != nil {
		log.Errorf("RPC_ChatUpdateUser error: %v", ret.Err)
	}
}
