package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Mail/rpcproc"
	"fmt"

	log "github.com/cihub/seelog"
)

var Inst Core.ICore

func serverInit(areaID string, serverID string) error {
	Inst = Core.New()
	info := Core.NewDefaultInfo()

	info.ServiceType = Const.Mail
	info.AreaID = areaID
	info.ServiceID = fmt.Sprintf("%s_%s_%s", info.ServiceType, areaID, serverID)
	info.RpcProc = rpcproc.NewRPCProc()
	log.Debugf("init core, serviceID=%s", info.ServiceID)
	return Inst.Init(info)
}

func serverDestroy() {
	Inst.Destroy()
}
