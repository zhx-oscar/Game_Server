package Space

import (
	"Cinder/Base/Message"
	//log "github.com/cihub/seelog"
	"time"
)

func (user *User) MsgProc(msg Message.IMessage) {

	switch msg.GetID() {
	case Message.ID_Client_User_Rpc:
		m := msg.(*Message.ClientUserRpc)
		user.onClientUserRpc(m)
	case Message.ID_Heart_Beat:
		m := msg.(*Message.HeartBeat)
		user.onClientHeartbeat(m)
	}

}

func (user *User) onClientUserRpc(msg *Message.ClientUserRpc) {
	switch msg.Target {
	case 0:
		user.SendToAllClient(msg)
	case 1:
		user.SendToAllClientExceptMe(msg)
	case 2:
		u := user.GetSpace().GetOwnerUser()
		if u != nil {
			_ = u.SendToClient(msg)
		}
	}
}

func (user *User) onClientHeartbeat(msg *Message.HeartBeat) {
	user.lastHeartbeatTimeMtx.Lock()
	user.lastHeartbeatTime = time.Now()
	user.lastHeartbeatTimeMtx.Unlock()

	msg.ServerTime = time.Now().UnixNano()
	_ = user.SendToClient(msg)
}

func (user *User) IsClientNetOK() bool {
	user.lastHeartbeatTimeMtx.Lock()
	lastHb := user.lastHeartbeatTime
	user.lastHeartbeatTimeMtx.Unlock()

	if lastHb.IsZero() {
		return true
	}

	if time.Now().Sub(lastHb) < 11*time.Second {
		return true
	}

	return false
}
