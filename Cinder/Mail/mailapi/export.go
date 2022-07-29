/* mailapi 包是邮件服的SDK工具包。
其他服调用 mailapi 操作邮件服。
使用说明见：[mail_usage.md](mail_usage.md)
*/

package mailapi

import (
	"Cinder/Mail/mailapi/internal"
	"Cinder/Mail/mailapi/types"
	"time"

	assert "github.com/arl/assertgo"
)

var iService IMailService

type Mail = types.Mail

func init() {
	iService = internal.NewMailService()
}

// SetMailService 设置邮件服务接口
func SetMailService(ms IMailService) {
	assert.True(ms != nil)
	iService = ms
}

// GetMailService 获取邮件服务接口
func GetMailService() IMailService {
	assert.True(iService != nil)
	return iService
}

func Login(userID string) error {
	return iService.Login(userID)
}

func Logout(userID string) error {
	return iService.Logout(userID)
}

func Send(m *Mail) error {
	return iService.Send(m)
}

func Broadcast(m *Mail) error {
	return iService.Broadcast(m)
}

func ListMail(userID string) ([]*Mail, error) {
	return iService.ListMail(userID)
}

func Delete(userID string, mailID string) error {
	return iService.Delete(userID, mailID)
}

func BatchDelete(userID string, mailIDs []string) error {
	return iService.BatchDelete(userID, mailIDs)
}

func MarkAsRead(userID string, mailID string) error {
	return iService.MarkAsRead(userID, mailID)
}

func MarkAsUnread(userID string, mailID string) error {
	return iService.MarkAsUnread(userID, mailID)
}

func MarkAttachmentsAsReceived(userID string, mailID string) error {
	return iService.MarkAttachmentsAsReceived(userID, mailID)
}

func MarkAttachmentsAsUnreceived(userID string, mailID string) error {
	return iService.MarkAttachmentsAsUnreceived(userID, mailID)
}

func SetExtData(userID string, mailID string, extData []byte) error {
	return iService.SetExtData(userID, mailID, extData)
}

// IMailService 通用邮件服务器接口
type IMailService interface {
	// 用户上下线，首次上线时间记录为出生时间
	Login(userID string) error
	Logout(userID string) error

	// 发送用户邮件和广播系统邮件
	Send(m *Mail) error
	// 弃用，广播系统邮件，请改用 Send()
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
	// 更改过期时间. 只能更改个人邮件，更改全服邮件则返回错误。
	SetExpireTime(userID string, mailID string, expireTime time.Time) error
	// 设置额外数据
	SetExtData(userID string, mailID string, extData []byte) error
}
