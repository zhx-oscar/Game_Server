package internal

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Matcher/rpcmsg"
	"Cinder/Matcher/svcid"
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/cihub/seelog"
)

type _RpcCaller struct {
	matcherSvcID string // 跨区匹配服服务号，为空表示本区
}

func newRpcCaller(matcherAreaID string, matcherServerID string) *_RpcCaller {
	if matcherAreaID == "" || matcherServerID == "" {
		return &_RpcCaller{}
	}

	return &_RpcCaller{
		matcherSvcID: svcid.FormatSvcID(Const.Matcher, matcherAreaID, matcherServerID),
	}
}

// rpcWithRet 调用Matcher服的RPC, 返回 RpcRet
func (r *_RpcCaller) rpcWithRet(methodName string, args ...interface{}) CRpc.RpcRet {
	log.Debugf("rpc: %s", methodName)

	// 允许直接输入 rpcmsg 中的 Request 类型，自动打包。(但是没有自动解包)
	args2, errArg := marshalReqArgs(args)
	if errArg != nil {
		return CRpc.RpcRet{
			Err: errArg,
		}
	}

	// log.Debugf("args2: %v", args2)
	var ret *CRpc.RpcRet
	if "" == r.matcherSvcID {
		ret = <-Core.Inst.RpcByType(Const.Matcher, methodName, args2...) // 本区
	} else {
		ret = <-Core.Inst.RpcByID(r.matcherSvcID, methodName, args2...) // 跨区
	}
	// log.Debugf("ret: %#v", ret)
	if ret != nil {
		return *ret
	}
	return CRpc.RpcRet{
		Err: fmt.Errorf("Core.Inst.RpcByType() channel pops nil"),
	}
}

// rpc 调用Matcher服的RPC, 返回 error
func (r *_RpcCaller) rpc(methodName string, args ...interface{}) error {
	ret := r.rpcWithRet(methodName, args...)
	if ret.Err != nil {
		return ret.Err
	}
	// 一般 RPC 仅返回一个错误串
	sErr := ret.Ret[0].(string)
	if sErr != "" {
		return errors.New(sErr)
	}
	return nil
}

// rpc 调用Matcher服的RPC, 返回 RPCResponse
func (r *_RpcCaller) rpc2(methodName string, args ...interface{}) (rpcmsg.RPCResponse, error) {
	rsp := rpcmsg.RPCResponse{}
	ret := r.rpcWithRet(methodName, args...)
	if ret.Err != nil {
		return rsp, ret.Err
	}
	// 返回 rpcmsg.RPCResponse json 和 errStr
	sErr := ret.Ret[1].(string)
	if sErr != "" {
		return rsp, errors.New(sErr)
	}
	json := ret.Ret[0].([]byte)
	if err := rsp.Unmarshal(json); err != nil {
		return rsp, err
	}
	return rsp, nil
}

func marshalReqArgs(args []interface{}) ([]interface{}, error) {
	result := []interface{}{}
	for _, arg := range args {
		if !rpcmsg.IsRequest(arg) {
			result = append(result, arg)
			continue
		}

		buf, errJson := json.Marshal(arg)
		if errJson != nil {
			return nil, errJson
		}
		result = append(result, buf)
	}
	return result, nil
}
