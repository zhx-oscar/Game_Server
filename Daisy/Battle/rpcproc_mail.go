package main

import (
	"Cinder/Mail/mailapi"
	mailtypes "Cinder/Mail/mailapi/types"
	"Cinder/Space"
	log "github.com/cihub/seelog"
)

// RPC_MailNotify 邮件服新邮件通知回调
// 其中 mailJsonData 是Mail结构的Json打包块, 可使用 mailapi.UnmarshalMail() 解包。
func (proc *_RPCProc) RPC_MailNotify(mailJsonData []byte) {
	log.Debugf("RPC_MailNotify: data=%d bytes", len(mailJsonData))
	mail, err := mailapi.UnmarshalMail(mailJsonData)
	if err != nil {
		log.Errorf("RPC_MailNotify unmarshal mail error: %s", err)
		return
	}

	if !mail.IsBroadcast {
		user, err2 := Space.Inst.GetUser(mail.To)
		if err2 != nil {
			_ = log.Error("RPC_AddFriendReq can't find targetUser: ", mail.To, err)
			return
		}
		u, ok := user.(*_User)
		if !ok || u == nil || u.role == nil {
			return
		}
		u.role.SyncInsertMailFromSrv(&mail)
	} else {
		Space.Inst.TraversalSpace(func(space Space.ISpace) {
			notifyMailInSpace(mail, space)
		})
	}
}

// notifyMailInSpace 在space中查找接收者，通知新邮件。
// 如果非广播邮件并找到接收者，则设置done为true.
func notifyMailInSpace(mail mailtypes.Mail, space Space.ISpace) {
	log.Info("notifyMailInSpace begin ")

	space.TraversalActor(func(actor Space.IActor) {
		role, ok := actor.(*_Role)
		if !ok {
			return
		}
		role.SyncInsertMailFromSrv(&mail)
	})
}
