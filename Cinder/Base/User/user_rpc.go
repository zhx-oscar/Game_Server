package User

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Const"
	"Cinder/Base/Message"
	"Cinder/Base/Util"
	"context"
	"errors"
	"time"
)

func (u *User) Rpc(srvType string, methodName string, args ...interface{}) chan *CRpc.RpcRet {

	ret, retID := u.genRpcRet()

	msg := &Message.UserRpcReq{
		UserID:     u.GetID(),
		MethodName: methodName,
		Args:       Message.PackArgs(args...),
		RetID:      retID,
	}

	if srvType == Const.Agent {
		err := u.SendToClient(msg)
		if err != nil {
			ret.Err = err
			ret.Done <- ret
			close(ret.Done)
			return ret.Done
		}
		return nil
	}

	err := u.SendToPeerServer(srvType, msg)
	if err != nil {
		ret.Err = err
		ret.Done <- ret
		close(ret.Done)
		return ret.Done
	}

	u.retList.Store(retID, ret)

	go func() {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer ctxCancel()

		select {
		case <-ctx.Done():
			ii, _ := u.retList.LoadOrStore(retID, nil)
			u.retList.Delete(retID)
			if ii != nil {
				r := ii.(*CRpc.RpcRet)
				r.Err = errors.New("time out " + retID)
				r.Done <- r
			}
		}
	}()

	return ret.Done
}

func (u *User) genRpcRet() (*CRpc.RpcRet, string) {
	ret := &CRpc.RpcRet{
		Done: make(chan *CRpc.RpcRet, 1),
	}
	return ret, Util.GetGUID()
}

func (u *User) OnRpcRet(retID string, err string, ret []interface{}) {

	ii, _ := u.retList.LoadOrStore(retID, nil)
	u.retList.Delete(retID)

	if ii == nil {
		u.Error("onRpcRet fail, get nil info", retID)
		return
	}

	info := ii.(*CRpc.RpcRet)
	if err != "" {
		info.Err = errors.New(err)
	}
	info.Ret = ret
	info.Done <- info
	close(info.Done)
}
