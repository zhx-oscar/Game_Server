package rpcproc

import (
	"Cinder/Chat/rpcproc/logic/user"
	"encoding/json"

	log "github.com/cihub/seelog"
)

type _RPCProcGetInfos struct {
}

// RPC_GetFriendInfos 查询一系列玩家信息
func (r *_RPCProcGetInfos) RPC_GetFriendInfos(userIDsJson []byte) []byte {
	// log.Debugf("RPC_GetFriendInfos: %v", string(userIDsJson))

	var userIDs []UserID
	if err := json.Unmarshal(userIDsJson, &userIDs); err != nil {
		log.Errorf("json unmarshal error: %v", err)
		return nil
	}

	infos, err := user.GetFriendInfos(userIDs)
	if err != nil {
		log.Errorf("failed to get infos: %v", err)
		return nil
	}

	buf, errBuf := json.Marshal(infos)
	if errBuf != nil {
		log.Errorf("failed to json marshal friend infos: %v", errBuf)
		return nil
	}
	return buf
}
