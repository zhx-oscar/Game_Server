package main

import (
	"Daisy/Proto"
)

// RPC_OneKeyMailColReq 一键领取附件
func (user *_User) RPC_OneKeyMailColReq() (int32, *Proto.MailAwardItems) {
	user.Debug("RPC_OneKeyMailColReq roleId ", user.role.GetID(), " ")
	return user.role.oneKeyGetAttachmentsOfMail()
}

// RPC_MailColReq 领取单封附件
func (user *_User) RPC_MailColReq(mailID string) (int32, *Proto.MailAwardItems) {
	user.Debug("RPC_MailColReq roleId ", user.role.GetID(), " ", mailID)
	return user.role.getAttachmentsOfMail(mailID)
}

// RPC_MailReadReq 读取邮件
func (user *_User) RPC_MailReadReq(mailID string) int32 {
	user.Debug("RPC_MailReadReq roleId ", user.role.GetID(), " ", mailID)
	return user.role.markMailAsRead(mailID)
}

// RPC_MailDelHasReadReq 删除已读邮件
func (user *_User) RPC_MailDelHasReadReq() int32 {
	user.Debug("RPC_MailDelReq roleId ", user.role.GetID(), " ")
	return user.role.deleteMailsHasRead()
}
