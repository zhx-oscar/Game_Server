package internal

import (
	"Cinder/Base/CRpc"
	"Cinder/Chat/chatapi/types"
	"encoding/json"

	log "github.com/cihub/seelog"
)

type FriendInfo = types.FriendInfo

func GetFollowingList(userID string) []FriendInfo {
	ret := Rpc("RPC_GetFollowingList", userID)
	return getFriendInfosFromRpcRet(ret)
}

func GetFollowerList(userID string) []types.FollowerInfo {
	ret := Rpc("RPC_GetFollowerList", userID)
	if ret.Err != nil {
		log.Errorf("rpc get follower list error: %v", ret.Err)
		return nil
	}
	bin := ret.Ret[0].([]byte)
	var result []types.FollowerInfo
	if err := json.Unmarshal(bin, &result); err != nil {
		log.Errorf("json unmarshal follower info list error: %v", err)
		return nil
	}
	return result
}

func GetFriendList(userID string) []FriendInfo {
	ret := Rpc("RPC_GetFriendList", userID)
	return getFriendInfosFromRpcRet(ret)
}

func GetFriendBlacklist(userID string) []FriendInfo {
	ret := Rpc("RPC_GetFriendBlacklist", userID)
	return getFriendInfosFromRpcRet(ret)
}

func getFriendInfosFromRpcRet(ret CRpc.RpcRet) []FriendInfo {
	if ret.Err != nil {
		log.Errorf("get friend infos rpc error: %v", ret.Err)
		return nil
	}
	bin := ret.Ret[0].([]byte)
	var result []FriendInfo
	if err := json.Unmarshal(bin, &result); err != nil {
		log.Errorf("json unmarshal friend info error: %v", err)
		return nil
	}
	return result
}

func GetFriendInfos(userIDs []string) []FriendInfo {
	// RPC 不支持 []string, 需要自己打包
	binIDs, err := json.Marshal(userIDs)
	if err != nil {
		log.Errorf("get friend infos marshal error: %v", err)
		return nil
	}

	ret := Rpc("RPC_GetFriendInfos", binIDs)
	return getFriendInfosFromRpcRet(ret)
}
