package Mailbox

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Core"
	"Cinder/Base/Message"
	"encoding/json"
	"errors"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

type _UserMailBox struct {
	ID     string
	SrvID  string
	UserID string

	retList sync.Map
}

func (mb *_UserMailBox) Rpc(methodName string, args ...interface{}) chan *CRpc.RpcRet {
	ret, retID := CRpc.NewRpcRet()
	msg := &Message.MailboxReq{
		TargetID:   mb.UserID,
		MethodName: methodName,
		Args:       Message.PackArgs(args...),
	}

	err := Core.Inst.Send(mb.SrvID, msg)
	if err != nil {
		ret.Err = err
		ret.Done <- ret
		close(ret.Done)

		return ret.Done
	}

	mb.retList.Store(retID, ret)

	time.AfterFunc(3*time.Second, func() {
		if v, ok := mb.retList.Load(retID); ok {
			mb.retList.Delete(retID)
			info := v.(*CRpc.RpcRet)
			info.Err = errors.New("time out")
			select {
			case info.Done <- info:
			default:
			}
		}
	})

	return nil
}

func (mb *_UserMailBox) onMailboxRet(retID string, err error, ret []interface{}) {
	v, ok := mb.retList.Load(retID)
	mb.retList.Delete(retID)
	if !ok {
		log.Error("onMailboxRet fail, get nil ", retID)
		return
	}

	info := v.(*CRpc.RpcRet)
	if err != nil {
		info.Err = err
	}
	info.Ret = ret
	select {
	case info.Done <- info:
	default:
	}
	close(info.Done)
}

func (mb *_UserMailBox) getType() uint8 {
	return TypeUser
}

func (mb *_UserMailBox) getMailboxID() string {
	return mb.ID
}

func (mb *_UserMailBox) marshal() ([]byte, error) {
	return json.Marshal(mb)
}

func (mb *_UserMailBox) unmarshal(data []byte) error {
	return json.Unmarshal(data, mb)
}
