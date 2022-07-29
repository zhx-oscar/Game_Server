package rpc

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Core"
	"Cinder/Matcher/matchapi/mtypes"
	"fmt"

	log "github.com/cihub/seelog"
)

// RpcByID 向SrvID调用RPC, 忽略返回值
func RpcByID(srvID mtypes.SrvID, methodName string, args ...interface{}) {
	if ret := RpcWithRet(srvID, methodName, args...); ret.Err != nil {
		log.Errorf("RpcByID(%s, %s) error: %v", srvID, methodName, ret.Err)
	}
}

// RpcWithRet 调用RPC, 返回 RpcRet
func RpcWithRet(srvID mtypes.SrvID, methodName string, args ...interface{}) CRpc.RpcRet {
	// log.Debugf("RpcWithRet: srvID=%s method=%s", srvID, methodName)
	ret := <-Core.Inst.RpcByID(string(srvID), methodName, args...)
	// log.Debugf("ret: %#v", ret)
	if ret != nil {
		return *ret
	}
	return CRpc.RpcRet{
		Err: fmt.Errorf("Core.Inst.RpcByID() channel pops nil"),
	}
}
