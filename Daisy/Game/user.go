package main

import (
	"Cinder/Chat/chatapi"
	"Cinder/Game"
	"Daisy/Prop"
	"time"
)

type _User struct {
	Game.User
	prop     *Prop.RoleProp
	chatUser chatapi.IUser // 聊天服入口

	curTeamSrvID string // 所处队伍所在的服务器ID
}

func (u *_User) Init() {
	u.prop = u.GetProp().(*Prop.RoleProp)
	t := time.Now().Unix()
	u.prop.SyncSetOnline(true, u.prop.Data.Base.LastLogoutTime, t)

	u.Debug("User Inited")
}

func (u *_User) Destroy() {
	u.Debug("User Destroy")
	t := time.Now().Unix()
	u.prop.SyncSetOnline(false, t, u.prop.Data.Base.LastLoginTime)
	u.chatOffline()
}

func (u *_User) Start() {
	u.Debug("User start")
	u.initWorldChannel()
	u.chatOnline()
	u.friendOnline()
	u.Debug("user start", u.prop.Data.Base.UID)
}

func (u *_User) GetPropInfo() (string, string) {
	return Prop.RolePropType, u.GetID()
}
