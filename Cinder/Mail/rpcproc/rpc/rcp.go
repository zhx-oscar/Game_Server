package rpc

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Core"
	"fmt"

	log "github.com/cihub/seelog"
)

// RpcByIDs 向一批ID调用RPC, 忽略返回值
func RpcByIDs(srvIDs []string, methodName string, args ...interface{}) {
	for _, id := range srvIDs {
		if ret := RpcWithRet(id, methodName, args...); ret.Err != nil {
			log.Errorf("RpcByID(%s, %s) error: %v", id, methodName, ret.Err)
		}
	}
}

// RpcWithRet 调用RPC, 返回 RpcRet
func RpcWithRet(srvID, methodName string, args ...interface{}) CRpc.RpcRet {
	// log.Debugf("RpcWithRet: srvID=%s method=%s", srvID, methodName)
	ret := <-Core.Inst.RpcByID(srvID, methodName, args...)
	// log.Debugf("ret: %#v", ret)
	if ret != nil {
		return *ret
	}
	return CRpc.RpcRet{
		Err: fmt.Errorf("Core.Inst.RpcByID() channel pops nil"),
	}
}
