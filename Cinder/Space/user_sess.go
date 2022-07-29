package Space

import (
	"Cinder/Base/Const"
	"Cinder/Base/Message"
)

func (user *User) OnPeerChange(srvType string, srvID string) {
	if srvType == Const.Agent {
		user.refreshAgentSess(srvID)
	}
}

func (user *User) refreshAgentSess(agentID string) {

	oldAgent := user.agentID
	user.agentID = agentID

	user.GetSpace().(_ISpace).onUserAgentChanged(user.GetID(), oldAgent, agentID)
}

func (user *User) SendToAllClient(msg Message.IMessage) {
	user.GetSpace().SendToAllClient(msg)
}

func (user *User) SendToAllClientExceptMe(msg Message.IMessage) {
	user.GetSpace().SendToAllClientExceptOne(msg, user.GetID())
}
