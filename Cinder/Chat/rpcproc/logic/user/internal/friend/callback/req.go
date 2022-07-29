package callback

import (
	"Cinder/Chat/rpcproc/logic/rpc"
	"Cinder/Chat/rpcproc/logic/user/internal/friend/dbutil"

	log "github.com/cihub/seelog"
)

// RpcOneAddFriendReq 回调一个 RPC_AddFriendReq
func RpcOneAddFriendReq(srvID string, req *dbutil.DocFriendRequest) {
	ret := rpc.Rpc(srvID, "RPC_AddFriendReq", string(req.UserID), string(req.FromID), req.ReqInfo)
	if ret.Err != nil {
		log.Errorf("RPC_AddFriendReq failed: %s", ret.Err)
	}
}

// RpcManyAddFriendReqs 回调多个 RPC_AddFriendReq
func RpcManyAddFriendReqs(srvID string, requests []*dbutil.DocFriendRequest) {
	for _, req := range requests {
		RpcOneAddFriendReq(srvID, req)
	}
}
