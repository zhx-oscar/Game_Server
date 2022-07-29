package main

import (
	"Cinder/Base/Const"
)

//SendCliNotify 发送非常驻通知，多参数用|分隔
func (user *_User) SendCliNotify(id uint32, args string) {
	user.Debug("SendCliNotify ", id, args)
	user.Rpc(Const.Agent, "RPC_SendCliNotify", id, args)
}

//ShowCliNotify 显示常驻通知，多参数用|分隔
func (user *_User) ShowCliNotify(id uint32, args string) {
	user.Debug("ShowCliNotify ", id, args)
	user.Rpc(Const.Agent, "RPC_ShowCliNotify", id, args)
}

//HideCliNotify 隐藏常驻通知，多参数用|分隔
func (user *_User) HideCliNotify(id uint32) {
	user.Debug("HideCliNotify ", id)
	user.Rpc(Const.Agent, "RPC_HideCliNotify", id)
}
