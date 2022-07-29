package rpcproc

import (
	"Cinder/Chat/rpcproc/logic/usermgr"
	"encoding/json"

	log "github.com/cihub/seelog"
)

type _RPCProcFollow struct {
}

func (r *_RPCProcFollow) RPC_FollowFriendReq(followerID, followedID string) {
	// log.Debugf("RPC_FollowFriendReq followerID=%v followedID=%v", followerID, followedID)
	user := usermgr.GetUserMgr().GetUser(UserID(followerID))
	if user == nil {
		log.Warnf("follower '%s' is not online", followerID)
		return
	}
	if err := user.GetFollowMgr().Follow(UserID(followedID)); err != nil {
		log.Errorf("'%s' failed to follow '%s': %v", followerID, followedID, err)
		return
	}
}

func (r *_RPCProcFollow) RPC_UnFollowFriendReq(unfollowerID, unfollowedID string) {
	// log.Debugf("RPC_UnFollowFriendReq unfollowerID=%v unfollowedID=%v", unfollowerID, unfollowedID)
	user := usermgr.GetUserMgr().GetUser(UserID(unfollowerID))
	if user == nil {
		log.Warnf("user '%s' is not online", unfollowerID)
		return
	}
	if err := user.GetFollowMgr().Unfollow(UserID(unfollowedID)); err != nil {
		log.Errorf("'%s' failed to unfollow '%s': %v", unfollowerID, unfollowedID, err)
		return
	}
}

func (r *_RPCProcFollow) RPC_GetFollowingList(followerID string) []byte {
	// log.Debugf("RPC_GetFollowingList followerID=%v", followerID)
	user := usermgr.GetUserMgr().GetUser(UserID(followerID))
	if user == nil {
		log.Warnf("user '%s' is not online", followerID)
		return nil
	}
	infos, err := user.GetFollowMgr().GetFollowingList()
	if err != nil {
		log.Errorf("'%s' failed to get following list: %v", followerID, err)
		return nil
	}
	bin, err2 := json.Marshal(infos)
	if err2 != nil {
		log.Errorf("marshal error: %v", err)
		return nil
	}
	return bin
}

func (r *_RPCProcFollow) RPC_GetFollowerList(followedID string) []byte {
	// log.Debugf("RPC_GetFollowerList followedID=%v", followedID)
	user := usermgr.GetUserMgr().GetUser(UserID(followedID))
	if user == nil {
		log.Warnf("user '%s' is not online", followedID)
		return nil
	}
	infos, err := user.GetFollowMgr().GetFollowerList()
	if err != nil {
		log.Errorf("'%s' failed to get follower list: %v", followedID, err)
		return nil
	}
	bin, err2 := json.Marshal(infos)
	if err2 != nil {
		log.Errorf("marshal error: %v", err)
		return nil
	}
	return bin
}
