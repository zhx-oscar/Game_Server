package Space

import (
	"Cinder/Base/Mailbox"
	BaseUser "Cinder/Base/User"
	"errors"
)

var rpcProcInst interface{}

func _Init(areaID string, serverID string, spacePT ISpace, userPT IUser, rpcProc interface{}) error {
	if areaID == "" {
		return errors.New("invalid AreaID")
	}
	if serverID == "" {
		return errors.New("invalid ServerID")
	}

	Inst = newServer()
	err := Inst.(*_Server).InitSrv(areaID, serverID, spacePT, userPT, rpcProc)
	if err != nil {
		return err
	}

	userMgrProc := BaseUser.NewUserMessageProc(Inst.(*_Server))
	Inst.GetNetNode().AddMessageProc(userMgrProc)
	Inst.GetNetNode().AddMessageProc(&_SrvMsgProc{})
	Inst.GetNetNode().AddMessageProc(Mailbox.GetDefaultMgr())

	if rpcProc != nil {
		ii, ok := rpcProc.(_IInit)
		if ok {
			ii.Init()
		}

		rpcProcInst = rpcProc
	}

	return nil
}

func _Destroy() {

	ii, ok := rpcProcInst.(_IDestroy)
	if ok {
		ii.Destroy()
	}

	if Inst != nil {
		Inst.(*_Server).DestroySrv()
	}
}
