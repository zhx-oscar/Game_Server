package bc

import (
	"Cinder/Chat/rpcproc/logic/rpc"
	ltypes "Cinder/Chat/rpcproc/logic/types"
	"encoding/json"

	log "github.com/cihub/seelog"
)

// GroupBroadcastMessage 群内广播消息
// 向在线成员发送消息，包括发送者。
func GroupBroadcastMessage(groupID ltypes.GroupID, groupOnlineMemberIDs map[ltypes.UserID]bool,
	fromRoleID ltypes.UserID, fromNick string, fromData []byte, msgContent []byte) {
	srvIDToTargets := getSrvIDToTargetsMap(groupOnlineMemberIDs)
	for srvID, targets := range srvIDToTargets {
		targetsJson, errJson := json.Marshal(targets)
		if errJson != nil {
			log.Errorf("failed to json marshal targets when broadcast group message, error=%s, srvID=%s", errJson, srvID)
			continue
		}
		ret := rpc.Rpc(srvID, "RPC_ChatRecvGroupMessageV2", string(groupID), targetsJson, string(fromRoleID), fromNick, fromData, msgContent)
		if ret.Err != nil {
			log.Errorf("rpc error: %s, srvID=%s", ret.Err, srvID)
			continue
		}
	}
}
