package test

import (
	"Cinder/Matcher/matchapi/mtypes"

	"github.com/stretchr/testify/require"
)

type _RpcProc struct {
	*require.Assertions
}

func NewRpcProc(a *require.Assertions) *_RpcProc {
	return &_RpcProc{
		Assertions: a,
	}
}

func (r *_RpcProc) RPC_MatchNotify(msgJson []byte) {
	// fmt.Printf("RPC_MatchNotify(msgJson=%s)\n", string(msgJson))
	msg, errJson := mtypes.UnmarshalNotifyOneSrvMsg(msgJson)
	r.NoError(errJson)
	for _, notifyMsg := range msg.Msgs {
		handleNotifyMsg(notifyMsg)
	}
}

func handleNotifyMsg(msg mtypes.NotifyMsgToOneSrv) {
	// ...
}
