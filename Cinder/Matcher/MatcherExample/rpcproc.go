package main

import (
	"Cinder/Matcher/matcherlib"
)

type RPCProc struct {
	matcherlib.RPCProc
}

// 可以添加自定义 RPC

func (r *RPCProc) RPC_MyRPCTest(arg string) uint32 {
	return 0
}

func (r *RPCProc) RPC_Ping() {
}
