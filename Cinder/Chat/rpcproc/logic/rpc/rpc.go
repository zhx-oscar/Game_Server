package rpc

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Core"
	"fmt"
)

// Rpc 回调其他服的RPC
func Rpc(srvID string, methodName string, args ...interface{}) CRpc.RpcRet {
	// log.Debugf("rpc: %s %s", srvID, methodName)
	ret := <-Core.Inst.RpcByID(srvID, methodName, args...)
	if ret != nil {
		return *ret
	}
	return CRpc.RpcRet{
		Err: fmt.Errorf("Core.Inst.RpcByID() channel pops nil"),
	}
}
