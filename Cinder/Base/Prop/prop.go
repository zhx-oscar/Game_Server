package Prop

import (
	"Cinder/Base/Const"
	"Cinder/Base/Message"
	_ "Cinder/Base/ServerConfig"
	"Cinder/Base/Util"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

type _IProp interface {
	InitProp(owner IPropOwner, prop interface{})
	DestroyProp()
}

type Prop struct {
	caller Util.ISafeCall
	owner  IPropOwner
}

func (p *Prop) InitProp(owner IPropOwner, propRealPtr interface{}) {
	p.owner = owner
	p.caller = Util.NewSafeCall(propRealPtr, viper.GetBool("Config.Recover"))
}

func (p *Prop) DestroyProp() {
	p.caller.SafeCallDestroy()
}

func (p *Prop) Sync(methodName string, args []byte, syncToDB bool, target ...int) {
	// 检查target
	for _, t := range target {
		switch t {
		case Target_Game:
			if p.owner.GetSrvNode().GetType() == Const.Game {
				panic("prop can't sync self")
			}

		case Target_Space:
			if p.owner.GetSrvNode().GetType() == Const.Space {
				panic("prop can't sync self")
			}

		default:

		}
	}

	if syncToDB && p.owner.GetDBSrvID() != "" {

		msg := &Message.PropNotify{
			ID:          p.owner.GetPropID(),
			Type:        p.owner.GetPropType(),
			MethodName:  methodName,
			Args:        args,
			ServiceType: p.owner.GetSrvNode().GetType(),
		}

		if err := p.owner.GetSrvNode().Send(p.owner.GetDBSrvID(), msg); err != nil {
			log.Error("Sync Send PropNotify err ", err, " Msg ", msg)
			return
		}
	}

	if p.owner.GetSync() != nil {
		p.owner.GetSync().SyncProp(methodName, args, target...)
	}
}

func (p *Prop) GetCaller() Util.ISafeCall {
	return p.caller
}
