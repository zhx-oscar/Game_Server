package rpcproc

type _RPCProc struct {
	_RPCProcUser
	_RPCProcGroup
	_RPCProcFollow
	_RPCProcFriend
	_RPCProcBlacklist
	_RPCProcGetInfos
}

func NewRPCProc() *_RPCProc {
	return &_RPCProc{}
}
