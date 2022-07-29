package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Base/User"
	"github.com/spf13/viper"
)

var Inst Core.ICore
var userMgr *UserMgr

func serverInit(areaID string, serverID string) error {

	Inst = Core.New()
	info := Core.NewDefaultInfo()

	info.ServiceType = Const.Agent
	info.AreaID = areaID
	info.ServiceID = info.ServiceType + "_" + areaID + "_" + serverID
	info.NetServerMessageProc = newNetServerProc()
	if v := viper.GetInt("Agent.PortMax"); v != 0 {
		info.PortMax = v
	}
	if v := viper.GetInt("Agent.PortMin"); v != 0 {
		info.PortMin = v
	}
	if v := viper.GetString("Agent.ListenAddr"); v != "" {
		info.ListenAddr = v
	}
	if v := viper.GetString("Agent.OuterAddr"); v != "" {
		info.OuterAddr = v
	}

	if err := Inst.Init(info); err != nil {
		return err
	}

	userMgr = newUserMgr()

	Inst.GetNetNode().AddMessageProc(User.NewUserMessageProc(userMgr))
	Inst.GetNetNode().AddMessageProc(&_SrvMsgProc{})

	return nil
}

func serverDestroy() {
	userMgr.Destroy()
	Inst.Destroy()
}
