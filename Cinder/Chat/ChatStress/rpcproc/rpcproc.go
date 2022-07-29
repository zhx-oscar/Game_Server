package rpcproc

type RpcProc struct {
}

func NewRpcProc() *RpcProc {
	return &RpcProc{}
}

func (r *RpcProc) RPC_ChatRecvP2PMessage(targetID string, targetData []byte, fromID string, fromNick string, fromData []byte, msgContent []byte) {
}
func (r *RpcProc) RPC_ChatRecvGroupMessageV2(groupID string, targetsJson []byte, fromID string, fromNick string, fromData []byte, msgContent []byte) {
}
func (r *RpcProc) RPC_AddFriendReq(targetID string, fromID string, reqInfo []byte) {}
func (r *RpcProc) RPC_AddFriendRet(targetID string, fromID string, ok bool)        {}
func (r *RpcProc) RPC_AddDelFollower(targetID string, targetData []byte, followerID string, isAdd bool) {
}
func (r *RpcProc) RPC_ChatUpdateUser(updateUserMsgJson []byte) {
}
func (r *RpcProc) RPC_FriendDeleted(targetID string, friendID string) {
}
