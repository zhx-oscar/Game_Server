package rpcproc

import (
	"Cinder/Matcher/matchapi/mtypes"
)

type _RpcProc struct {
}

func NewRpcProc() *_RpcProc {
	return &_RpcProc{}
}

func (r *_RpcProc) RPC_MatchNotify(msgJson []byte) {
	// fmt.Printf("RPC_MatchNotify(msgJson=%s)\n", string(msgJson))
	_, errJson := mtypes.UnmarshalNotifyOneSrvMsg(msgJson)
	if errJson != nil {
		panic(errJson)
	}
}
