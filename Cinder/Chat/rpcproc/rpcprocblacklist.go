package rpcproc

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/usermgr"
	"encoding/json"

	log "github.com/cihub/seelog"
)

type _RPCProcBlacklist struct {
}

func (r *_RPCProcBlacklist) RPC_AddFriendToBlacklist(actorID, blacklistedID string) {
	// log.Debugf("RPC_AddFriendToBlacklist actorID=%v blacklistedID=%v", actorID, blacklistedID)
	user := usermgr.GetUserMgr().GetUser(UserID(actorID))
	if user == nil {
		log.Warnf("user '%s' is not online", actorID)
		return
	}
	if err := user.GetBlacklist().Add(UserID(blacklistedID)); err != nil {
		log.Errorf("failed to add '%s' to the blacklist of '%s': %v", blacklistedID, actorID, err)
		return
	}
}

func (r *_RPCProcBlacklist) RPC_RemoveFriendFromBlacklist(actorID, blacklistedID string) {
	// log.Debugf("RPC_RemoveFriendFromBlacklist actorID=%v blacklistedID=%v", actorID, blacklistedID)
	user := usermgr.GetUserMgr().GetUser(UserID(actorID))
	if user == nil {
		log.Warnf("user '%s' is not online", actorID)
		return
	}
	if err := user.GetBlacklist().Remove(UserID(blacklistedID)); err != nil {
		log.Errorf("failed to remove '%s' from the blacklist of '%s': %v", blacklistedID, actorID, err)
		return
	}
}

func (r *_RPCProcBlacklist) RPC_GetFriendBlacklist(actorID string) []byte {
	// log.Debugf("RPC_GetFriendBlacklist actorID=%v", actorID)
	user := usermgr.GetUserMgr().GetUser(UserID(actorID))
	if user == nil {
		log.Warnf("user '%s' is not online", actorID)
		return nil
	}
	var infos []*chatapi.FriendInfo = user.GetBlacklist().GetBlacklistInfos()
	bin, err := json.Marshal(infos)
	if err != nil {
		log.Errorf("json marshal error: %v", err)
		return nil
	}
	return bin
}
