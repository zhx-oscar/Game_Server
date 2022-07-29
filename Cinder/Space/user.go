package Space

import (
	"Cinder/Base/Const"
	"Cinder/Base/Message"
	"Cinder/Base/Net"
	BaseUser "Cinder/Base/User"
	"sync"
	"time"
)

type _IUser interface {
	IUser
	GetAgentID() string
	IsClientNetOK() bool
}

type User struct {
	BaseUser.User

	agentID   string
	agentSess Net.ISess

	space ISpace

	lastHeartbeatTimeMtx sync.Mutex
	lastHeartbeatTime    time.Time
}

type _IBMsgProc interface {
	MsgProc(msg Message.IMessage)
}

func (user *User) InitBase() {
	user.space = user.GetUserData().(ISpace)

	agentID, err := user.GetPeerServerID(Const.Agent)
	if err != nil {
		user.Warn("InitBase space user create but no agent part existed")
	} else {
		user.refreshAgentSess(agentID)
	}

	user.SetParentCaller(user.space.(_ISpace).GetSafeCall())
}

func (user *User) LateInitBase() {
	user.clientNotifyUserEnter()
	user.space.(_ISpace).onAddUser(user.GetRealPtr().(IUser))
}

func (user *User) DestroyBase() {

	user.space.(_ISpace).onRemoveUser(user.GetRealPtr().(IUser))
	user.clientNotifyUserLeave()
}

func (user *User) GetSpace() ISpace {
	return user.space
}

func (user *User) GetAgentID() string {
	return user.agentID
}

func (user *User) LoopBase() {

}
