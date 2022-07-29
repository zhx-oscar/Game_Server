package main

import (
	"Cinder/Mail/mailapi"
	"Daisy/Proto"
)

func (u *_User) mailOnline() {
	u.Debug("User mailOnline ")
	/*
		u.ms = mailapi.GetMailService()
		if u.ms == nil {
			u.Error("[mailOnline] u.mailService is nil")
			return
		}

		err := u.ms.Login(u.GetID())
		if err != nil {
			u.Error("[mailOnline] mail login error ", err)
			return
		}
		mails, errList := u.ms.ListMail(u.GetID())
		if errList != nil {
			u.Error("[mailOnline] list mail error: ", err)
			return
		}

		// 排序
		sortMails := u.

	*/
}

func (u *_User) mailSortFromMailApi(mails []*mailapi.Mail) []*Proto.Mail {
	/*
		now := time.Now().Unix()
		sortMails := make([]*Proto.Mail, 0)
		for i := 0; i < len(mails); i++ {
			mail :=
		}
	*/

	return nil
}
