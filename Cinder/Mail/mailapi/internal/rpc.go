package internal

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"fmt"
)

// rpcWithRet 调用Mail服的RPC, 返回 RpcRet
func rpcWithRet(methodName string, args ...interface{}) CRpc.RpcRet {
	// log.Debugf("rpc: %s", methodName)
	ret := <-Core.Inst.RpcByType(Const.Mail, methodName, args...)
	// log.Debugf("ret: %#v", ret)
	if ret != nil {
		return *ret
	}
	return CRpc.RpcRet{
		Err: fmt.Errorf("Core.Inst.RpcByType() channel pops nil"),
	}
}

// rpc 调用Mail服的RPC, 返回 error
func rpc(methodName string, args ...interface{}) error {
	ret := rpcWithRet(methodName, args...)
	if ret.Err != nil {
		return ret.Err
	}
	// 一般 RPC 仅返回一个错误串
	sErr := ret.Ret[0].(string)
	if sErr != "" {
		return fmt.Errorf("rpc %s returns error: %s", methodName, sErr)
	}
	return nil
}
