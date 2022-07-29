package DBAgent

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"errors"
	"fmt"
)

var rpcProcInst interface{}

type _IInit interface {
	Init()
}

type _IDestroy interface {
	Destroy()
}

func _Init(areaID string, serverID string, rpcProc interface{}) error {
	if areaID == "" {
		return errors.New("invalid AreaID")
	}
	if serverID == "" {
		return errors.New("invalid ServerID")
	}

	Inst = Core.New()

	info := Core.NewDefaultInfo()
	info.ServiceType = Const.DB
	info.AreaID = areaID
	info.ServiceID = fmt.Sprintf("%s_%s_%s", info.ServiceType, areaID, serverID)
	info.RpcProc = rpcProc
	err := Inst.Init(info)
	if err != nil {
		return err
	}

	propMgr = newPropMgr()

	Inst.GetNetNode().AddMessageProc(&_SrvMsgProc{})

	ii, ok := rpcProcInst.(_IInit)
	if ok {
		ii.Init()
	}

	return nil
}

func _Destroy() {
	ii, ok := rpcProcInst.(_IDestroy)
	if ok {
		ii.Destroy()
	}

	propMgr.Destroy()
	Inst.Destroy()
}
