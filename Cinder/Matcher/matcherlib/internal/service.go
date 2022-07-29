package internal

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Matcher/matcherlib/internal/rpcproc/room"
	"Cinder/Matcher/matcherlib/ltypes"
	"Cinder/Matcher/svcid"

	assert "github.com/arl/assertgo"
)

var rpcProcInst interface{}

type _IInit interface {
	Init()
}

type _IDestroy interface {
	Destroy()
}

func Init(areaID string, serverID string, rpcProc interface{}, roomEvtHdlr ltypes.IRoomEventHandler) error {
	assert.True(roomEvtHdlr != nil)
	room.SetRoomEventHandler(roomEvtHdlr)

	_ = Core.New()
	info := Core.NewDefaultInfo()
	info.ServiceType = Const.Matcher
	info.AreaID = areaID
	info.ServiceID = svcid.FormatSvcID(info.ServiceType, areaID, serverID)
	info.RpcProc = rpcProc
	err := Core.Inst.Init(info)
	if err != nil {
		return err
	}

	rpcProcInst = rpcProc
	ii, ok := rpcProcInst.(_IInit)
	if ok {
		ii.Init()
	}

	return nil
}

func Destroy() {
	ii, ok := rpcProcInst.(_IDestroy)
	if ok {
		ii.Destroy()
	}

	Core.Inst.Destroy()
}
