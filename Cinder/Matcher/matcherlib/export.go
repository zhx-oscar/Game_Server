package matcherlib

import (
	"Cinder/Base/Core"
	"Cinder/Matcher/matcherlib/internal"
	"Cinder/Matcher/matcherlib/internal/rpcproc"
	"Cinder/Matcher/matcherlib/ltypes"

	assert "github.com/arl/assertgo"
)

var Inst Core.ICore

type (
	RPCProc struct {
		rpcproc.RPCProc
	}

	IMatcherLibRPCProcer interface {
		MatcherLibRPCProc_60daab5588df4747bf7c6aa23abe4552()
	}
) // type

func (r *RPCProc) MatcherLibRPCProc_60daab5588df4747bf7c6aa23abe4552() {}

// Init 初始化。
// rpcProc 必须包含 RPCProc.
func Init(areaID string, serverID string, rpcProc IMatcherLibRPCProcer, roomEvtHdlr ltypes.IRoomEventHandler) error {
	assert.True(roomEvtHdlr != nil)
	return internal.Init(areaID, serverID, rpcProc, roomEvtHdlr)
}

func Destroy() {
	internal.Destroy()
}
