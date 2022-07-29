package Game

import (
	"Cinder/Base/Core"
	"Cinder/Base/Message"
	"Cinder/Cache"
)

func (user *User) MsgProc(msg Message.IMessage) {
	switch msg.GetID() {
	case Message.ID_PropObject_Open_Req:
		m := msg.(*Message.PropObjectOpenReq)
		user.onPropObjectOpenReq(m)
	case Message.ID_PropObject_Close_Req:
		m := msg.(*Message.PropObjectCloseReq)
		user.onPropObjectCloseReq(m)
	}
}

func (user *User) onPropObjectOpenReq(m *Message.PropObjectOpenReq) {
	srvID, err := Cache.GetPropObjectSrvID(m.Typ, m.ID)
	if err != nil {
		user.Error("onPropObjectOpenReq GetPropObjectSrvID err ", err)
		return
	}

	m.UserID = user.GetID()
	if err = Core.Inst.Send(srvID, m); err != nil {
		user.Error("onPropObjectOpenReq Send Msg err ", err)
	}
}

func (user *User) onPropObjectOpenRet(m *Message.PropObjectOpenRet) {
	user.propObjectMap.Store(m.ID, m.SrvID)
	user.SendToClient(m)
}

func (user *User) onPropObjectCloseReq(m *Message.PropObjectCloseReq) {
	srvID, ok := user.propObjectMap.Load(m.ID)
	if !ok {
		user.Error("onPropObjectCloseReq prop object not exist")
		return
	}

	m.UserID = user.GetID()
	if err := Core.Inst.Send(srvID.(string), m); err != nil {
		user.Error("onPropObjectCloseReq Send err", err)
	}
}

func (user *User) onPropObjectCloseRet(m *Message.PropObjectCloseRet) {
	user.propObjectMap.Delete(m.ID)
	user.SendToClient(m)
}
