package Game

import (
	"Cinder/Base/Core"
	"Cinder/Base/Message"
	"Cinder/Base/Util"
	"github.com/spf13/viper"

	log "github.com/cihub/seelog"
)

type _SrvMsgProc struct {
}

func (proc *_SrvMsgProc) MessageProc(srcAddr string, message Message.IMessage) {

	switch message.GetID() {
	case Message.ID_User_Login_Req:
		proc.onUserLoginReq(srcAddr, message.(*Message.UserLoginReq))
	case Message.ID_PropObject_Open_Ret:
		proc.onPropObjectOpenRet(srcAddr, message.(*Message.PropObjectOpenRet))
	case Message.ID_PropObject_Close_Ret:
		proc.onPropObjectCloseRet(srcAddr, message.(*Message.PropObjectCloseRet))
	case Message.ID_User_Destroy_Req:
		proc.onUserDestroyReq(srcAddr, message.(*Message.UserDestroyReq))
	}

}

func (proc *_SrvMsgProc) onUserLoginReq(srvAddr string, msg *Message.UserLoginReq) {

	go func() {
		log.Info("game server receive user login req " + msg.UserID)
		defer func() {
			if err := recover(); err != nil {
				log.Error("onUserLoginReq panic", err)
				if !viper.GetBool("Config.Recover") {
					panic(err)
				} else {
					log.Error(Util.GetPanicStackString())
				}
			}
		}()


		_, createNew, err := UserMgr.GetOrCreateUser(msg.UserID, nil)
		if err != nil {
			log.Error("get or create game user failed " + err.Error() + "  " + msg.UserID)
			return
		}

		if !createNew {
			log.Error("shouldn't been here ", msg.UserID)
		}
	}()
}

func (proc *_SrvMsgProc) onUserDestroyReq(srcAddr string, msg *Message.UserDestroyReq) {

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("onUserDestroyReq panic", err)
				if !viper.GetBool("Config.Recover") {
					panic(err)
				} else {
					log.Error(Util.GetPanicStackString())
				}
			}
		}()

		if err := UserMgr.DestroyUser(msg.UserID); err != nil {
			log.Error("onUserDestroyReq DestroyUser err ", err, " UserID ", msg.UserID)
		}
		if err := Core.Inst.Send(srcAddr, &Message.UserDestroyRet{UserID: msg.UserID}); err != nil {
			log.Error("onUserDestroyReq Send retMsg err ", err, " UserID ", msg.UserID)
		}
	}()
}

type _IUserPropObjectProc interface {
	onPropObjectOpenRet(m *Message.PropObjectOpenRet)
	onPropObjectCloseRet(m *Message.PropObjectCloseRet)
}

func (proc *_SrvMsgProc) onPropObjectOpenRet(srvAddr string, msg *Message.PropObjectOpenRet) {

	user, err := UserMgr.GetUser(msg.UserID)
	if err != nil {
		return
	}

	user.(_IUserPropObjectProc).onPropObjectOpenRet(msg)
}

func (proc *_SrvMsgProc) onPropObjectCloseRet(srvAddr string, msg *Message.PropObjectCloseRet) {

	user, err := UserMgr.GetUser(msg.UserID)
	if err != nil {
		return
	}

	user.(_IUserPropObjectProc).onPropObjectCloseRet(msg)
}
