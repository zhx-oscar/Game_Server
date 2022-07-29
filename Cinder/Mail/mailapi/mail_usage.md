# 邮件服使用说明

本文最新版本见：
http://gitlab.ztgame.com/Cinder/Server/Cinder/blob/master/Mail/mailapi/mail_usage.md

本文说明如何接入邮件服。邮件服的功能和部署见：[../README.md](../README.md)

请使用 `Cinder/Mail/mailapi` 包接入邮件服。仅支持 go 语言并应用 Cinder 框架开发的游戏。

## 示例
```
roleID := ...
ms := mailapi.GetMailService()
if err := ms.Login(roleID); err != err {
	log.Errorf("mail login error: %s", err)
	return
}

// 获取邮件列表
mails, errList := ms.ListMail(roleID)
if errList != nil {
	log.Errorf("list mail error: %s", errList)
	return
}
log.Debugf("list mail: %d", len(mails))

if len(mails) == 0 {
	return
}
mailStrID := mails[0].ID
if err := ms.MarkAsRead(roleID, mailStrID); err != nil {
	log.Errorf("mark mail as read error: %s", err)
	return
}
```

## 回调 RPC

要求 Login() 的调用服实现回调 RPC：

	func RPC_MailNotify(mailJsonData []byte) {}
		其中 mailJsonData 是Mail结构的Json打包块, 可使用 mailapi.UnmarshalMail() 解包。

用户上下线时用 `Login()`, `Logou()` 通知邮件服，用户在线就会收到新邮件通知，即 `RPC_MailNotify` 回调。
`Login()`/`Logout()` 仅影响新邮件通知，其他邮件操作不要求在线.
对于系统广播邮件，`RPC_MailNotify` 是对每个服广播新邮件通知，而不是对每个用户调用。

## 函数与接口

可以直接调用 mailapi 包的自由函数，也可以先用 `GetMailService()` 获取 `IMailService`, 通过该接口调用方法。

```
// GetMailService 获取邮件服务接口
func GetMailService() IMailService
```

自由函数和接口的功能是一样的。如：
```
mailapi.Login("...")
```
等同于
```
mailapi.GetMailService().Login("...")
```

接口功能如下：
```
	// 用户上下线，首次上线时间记录为出生时间
	Login(userID string) error
	Logout(userID string) error

	// 发送用户邮件
	Send(m *Mail) error
	// 广播系统邮件
	Broadcast(m *Mail) error

	// 列举邮件
	ListMail(userID string) ([]*Mail, error)
	// 删除邮件
	Delete(userID string, mailID string) error
	BatchDelete(userID string, mailIDs []string) error

	// 设置邮件状态：已读/未读/已收附件/未收附件
	MarkAsRead(userID string, mailID string) error
	MarkAsUnread(userID string, mailID string) error
	MarkAttachmentsAsReceived(userID string, mailID string) error
	MarkAttachmentsAsUnreceived(userID string, mailID string) error
```

## 数据类型

* `Mail struct` 邮件
* `Attachment struct` 附件
