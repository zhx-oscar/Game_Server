package Prop

import (
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	"Daisy/Proto"
)

func (u *RoleProp) SyncMailClear(mailType int32) {
	u.Sync("MailBoxClear", Message.PackArgs(mailType), false, Prop.Target_Client)
}
func (u *RoleProp) MailBoxClear() {
	u.Data.MailBox.Mails = nil
}

func (u *RoleProp) SyncInsertMail(mail *Proto.Mail) {
	u.Sync("MailBoxInsertMail", Message.PackArgs(mail), false, Prop.Target_Client)
}
func (u *RoleProp) MailBoxInsertMail(mail *Proto.Mail) {
	u.Data.MailBox.Mails = append(u.Data.MailBox.Mails, mail)
}

func (u *RoleProp) SyncMailBoxRemoveMail(mailStrID string) {
	u.Sync("MailBoxRemoveMail", Message.PackArgs(mailStrID), false, Prop.Target_Client)
}
func (u *RoleProp) MailBoxRemoveMail(mailStrID string) {
	mails := u.Data.MailBox.Mails
	for i := 0; i < len(mails); i++ {
		if mails[i].MailID == mailStrID {
			u.Data.MailBox.Mails = append(mails[:i], mails[i+1:]...)
			break
		}
	}
}

func (u *RoleProp) SyncUpdateMailRead(mailStrID string, read bool) {
	u.Sync("MailBoxMailRead", Message.PackArgs(mailStrID, read), false, Prop.Target_Client)
}
func (u *RoleProp) MailBoxMailRead(mailStrID string, read bool) {
	mails := u.Data.MailBox.Mails
	for i := 0; i < len(mails); i++ {
		if mails[i].MailID == mailStrID {
			mails[i].IsRead = read
			break
		}
	}
}

func (u *RoleProp) SyncUpdateMailReceived(mailStrID string, received bool) {
	u.Sync("MailBoxMailReceived", Message.PackArgs(mailStrID, received), false, Prop.Target_Client)
}
func (u *RoleProp) MailBoxMailReceived(mailStrID string, received bool) {
	mails := u.Data.MailBox.Mails
	for i := 0; i < len(mails); i++ {
		if mails[i].MailID == mailStrID {
			mails[i].IsReceived = received
			break
		}
	}
}

// SyncUpdateMailExpireTimeAndReadTime 更新邮件过期时间和已读时间
func (u *RoleProp) SyncUpdateMailExpireTimeAndReadTime(mailStrID string, expireTime int64, readTime int64) {
	u.Sync("MailBoxUpdateExpireTimeAndReadTime", Message.PackArgs(mailStrID, expireTime, readTime), false, Prop.Target_Client)
}
func (u *RoleProp) MailBoxUpdateExpireTimeAndReadTime(mailStrID string, expireTime int64, readTime int64) {
	mails := u.Data.MailBox.Mails
	for i := 0; i < len(mails); i++ {
		if mails[i].MailID == mailStrID {
			mails[i].ExpireTime = expireTime
			mails[i].ReadTime = readTime
			break
		}
	}
}
