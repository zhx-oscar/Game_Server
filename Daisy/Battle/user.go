package main

import (
	cConst "Cinder/Base/Const"
	"Cinder/Space"
)

type _User struct {
	Space.User

	role *_Role
}

func (user *_User) Init() {
	user.setRole()

	if user.role != nil {
		user.role.OnOnline()
	}
	user.GetSpace().(*_Team).MemberOnline(user.GetID())

	user.Info("Init")
}

func (user *_User) Destroy() {
	if user.role != nil {
		user.role.OnOffline()
		user.role = nil
	}
	user.GetSpace().(*_Team).MemberOffline(user.GetID())

	user.Info("Destroy")
}

func (user *_User) GetPropInfo() (string, string) {
	return "", ""
}

func (user *_User) setRole() {
	team := user.GetSpace().(*_Team)
	ia := team.GetActorByUserID(user.GetID())
	if ia == nil {
		user.Errorf("online get role err")
		return
	}

	user.role = ia.(*_Role)
	return
}

func (user *_User) GetRole() *_Role {
	return user.role
}

// Start user创建完成 （调用role.online时 user并没有创建完成）
func (user *_User) Start() {
	user.addToTeamChannel()
}

// addToTeamChannel 玩家初始化完，加入队伍后，加入到队伍聊天频道
func (user *_User) addToTeamChannel() {
	team := user.GetSpace().(*_Team)
	tid := team.getTeamChatChannelKey()
	user.Rpc(cConst.Game, "RPC_TeamChatChannelAddMember", tid)
}
