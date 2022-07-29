package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Chat/rpcproc"
	"errors"
	"fmt"
)

var Inst Core.ICore

func serverInit(areaID string, serverID string) error {
	if areaID == "" {
		return errors.New("invalid AreaID")
	}
	if serverID == "" {
		return errors.New("invalid Server ID")
	}

	Inst = Core.New()
	info := Core.NewDefaultInfo()

	info.ServiceType = Const.Chat
	info.AreaID = areaID
	info.ServiceID = fmt.Sprintf("%s_%s_%s", info.ServiceType, areaID, serverID)
	info.RpcProc = rpcproc.NewRPCProc()
	return Inst.Init(info)
}

func serverDestroy() {
	Inst.Destroy()
}
