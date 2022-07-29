package Game

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Base/Mailbox"
	"Cinder/Base/Message"
	BaseUser "Cinder/Base/User"
	"errors"
	"fmt"
)

type _IBMsgProc interface {
	MsgProc(msg Message.IMessage)
}

var rpcProcInst interface{}

type _IInit interface {
	Init()
}

type _IDestroy interface {
	Destroy()
}

func _Init(areaID string, serverID string, userProto BaseUser.IUser, rpcProc interface{}) error {
	if areaID == "" {
		return errors.New("invalid AreaID")
	}
	if serverID == "" {
		return errors.New("invalid Server ID")
	}

	_ = Core.New()

	info := Core.NewDefaultInfo()
	info.ServiceType = Const.Game
	info.AreaID = areaID
	info.ServiceID = fmt.Sprintf("%s_%s_%s", info.ServiceType, areaID, serverID)
	info.RpcProc = rpcProc
	err := Core.Inst.Init(info)
	if err != nil {
		return err
	}

	UserMgr = BaseUser.NewUserMgr(userProto, Core.Inst, true)

	userMgrProc := BaseUser.NewUserMessageProc(UserMgr)
	Core.Inst.GetNetNode().AddMessageProc(userMgrProc)
	Core.Inst.GetNetNode().AddMessageProc(&_SrvMsgProc{})
	Core.Inst.GetNetNode().AddMessageProc(Mailbox.GetDefaultMgr())

	ii, ok := rpcProcInst.(_IInit)
	if ok {
		ii.Init()
	}

	return nil
}

func _Destroy() {
	UserMgr.Destroy()

	ii, ok := rpcProcInst.(_IDestroy)
	if ok {
		ii.Destroy()
	}

	Core.Inst.Destroy()
}
