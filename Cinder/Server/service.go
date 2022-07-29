package Server

import (
	"Cinder/Base/Core"
	"fmt"
)

var rpcProcInst interface{}

type _IInit interface {
	Init()
}

type _IDestroy interface {
	Destroy()
}

func _Init(srvType string, areaID string, serverID string, rpcProc interface{}) error {

	_ = Core.New()

	info := Core.NewDefaultInfo()
	info.ServiceType = srvType
	info.AreaID = areaID
	info.ServiceID = fmt.Sprintf("%s_%s_%s", srvType, areaID, serverID)
	info.RpcProc = rpcProc

	rpcProcInst = rpcProc

	ret := Core.Inst.Init(info)

	if ret == nil {
		ii, ok := rpcProcInst.(_IInit)
		if ok {
			ii.Init()
		}
	}

	return ret
}

func _Destroy() {

	ii, ok := rpcProcInst.(_IDestroy)
	if ok {
		ii.Destroy()
	}

	Core.Inst.Destroy()
}
