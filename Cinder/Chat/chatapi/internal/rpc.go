package internal

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"fmt"
)

// Rpc 调用Chat服的RPC
func Rpc(methodName string, args ...interface{}) CRpc.RpcRet {
	// log.Debugf("rpc: %s", methodName)
	ret := <-Core.Inst.RpcByType(Const.Chat, methodName, args...)
	if ret != nil {
		return *ret
	}
	return CRpc.RpcRet{
		Err: fmt.Errorf("Core.Inst.RpcByType() channel pops nil"),
	}
}
