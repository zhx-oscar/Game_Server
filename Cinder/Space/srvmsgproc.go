package Space

import (
	"Cinder/Base/Message"
	log "github.com/cihub/seelog"
)

type _SrvMsgProc struct {
}

func (proc *_SrvMsgProc) MessageProc(srvAddr string, message Message.IMessage) {
	switch message.GetID() {
	case Message.ID_User_Destroy_Req:
		go proc.onUserDestroyReq(srvAddr, message.(*Message.UserDestroyReq))
	}
}

func (proc *_SrvMsgProc) onUserDestroyReq(srvAddr string, msg *Message.UserDestroyReq) {
	if err := Inst.LeaveSpace(msg.UserID); err != nil {
		log.Error("onUserDestroyReq LeaveSpace err ", err, " UserID ", msg.UserID)
	}

	if err := Inst.Send(srvAddr, &Message.UserDestroyRet{UserID: msg.UserID}); err != nil {
		log.Error("onUserDestroyReq Send UserDestroyRet err ", err, " UserID ", msg.UserID)
	}
}
