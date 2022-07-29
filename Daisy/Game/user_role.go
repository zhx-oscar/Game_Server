package main

import (
	"Daisy/Const"
	"Daisy/DB"
	"Daisy/ErrorCode"
)

func (u *_User) RPC_SetName(name string, head uint32) int32 {
	if u.prop.Data.Base.Name != "" {
		u.Error("重复取名")
		return ErrorCode.Failure
	}

	if Const.UTF8Width(name) > 10 {
		u.Error("名字长度过长")
		return ErrorCode.NameInvalid
	}

	util := DB.NameUtil(name)
	ok, err := util.LockName()
	if err != nil {
		u.Error(err)
		return ErrorCode.DBOpErr
	}
	if !ok {
		u.Error("名字重复")
		return ErrorCode.NameDuplicate
	}

	err = util.SetRoleID(u.GetID())
	if err != nil {
		util.Del()
		u.Error(err)
		return ErrorCode.DBOpErr
	}

	u.prop.SyncSetName(name, head)
	u.Info("RPC_SetName success")

	u.updateChatUser() // 不判断返回值。更新聊天服数据失败了，并不妨碍取名的流程
	return ErrorCode.Success
}

// updateChatUser 更新聊天服名字头像信息
func (u *_User) updateChatUser() int32 {
	err := u.chatUser.SetNick(u.prop.Data.Base.Name)
	if err != nil {
		u.Error("[updateChatUser] setNick failed ", err)
		return ErrorCode.Failure
	}
	return u.updateChatUserActivateData()
}
