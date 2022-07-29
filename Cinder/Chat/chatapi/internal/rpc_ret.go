package internal

import (
	"Cinder/Base/CRpc"
	"errors"
	"fmt"
)

func GetStringErrorWithHint(ret CRpc.RpcRet, hint string) error {
	return formatError(hint, GetStringError(ret))
}

func formatError(hint string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", hint, err)
}

func GetStringError(ret CRpc.RpcRet) error {
	if first, err := GetRpcFirstRet(ret); err != nil {
		return err
	} else if errStr, ok := first.(string); !ok {
		return errors.New("rpc returns non-string error")
	} else if errStr != "" {
		return errors.New(errStr)
	}
	return nil
}

// GetRpcFirstRet 返回RPC第一个返回值
func GetRpcFirstRet(ret CRpc.RpcRet) (interface{}, error) {
	if ret.Err != nil {
		return nil, ret.Err
	}
	if len(ret.Ret) == 0 {
		return nil, fmt.Errorf("returns nil")
	}
	return ret.Ret[0], nil
}
