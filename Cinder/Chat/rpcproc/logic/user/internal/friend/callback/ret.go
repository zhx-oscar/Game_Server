package callback

import (
	"Cinder/Chat/rpcproc/logic/rpc"
	"Cinder/Chat/rpcproc/logic/user/internal/friend/dbutil"

	log "github.com/cihub/seelog"
)

// RpcOneAddFriendRet 回调一个 RPC_AddFriendRet
func RpcOneAddFriendRet(srvID string, rsp *dbutil.DocFriendResponse) {
	ret := rpc.Rpc(srvID, "RPC_AddFriendRet",
		string(rsp.UserID), string(rsp.ResponderID), rsp.OK)
	if ret.Err != nil {
		log.Errorf("RPC_AddFriendRet failed: %s", ret.Err)
	}
}

// RpcManyAddFriendRets 回调多个 RPC_AddFriendRet
func RpcManyAddFriendRets(srvID string, responses []*dbutil.DocFriendResponse) {
	for _, rsp := range responses {
		RpcOneAddFriendRet(srvID, rsp)
	}
}
