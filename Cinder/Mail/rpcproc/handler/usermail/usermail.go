package usermail

import (
	"Cinder/Mail/mailapi/types"
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/internal/maildoc"
	"Cinder/Mail/rpcproc/handler/internal/notify"
	"context"
	"fmt"

	log "github.com/cihub/seelog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Send 发送普通邮件。
// 存DB, 然后查询用户，如果在线就通知
func Send(mail types.Mail) error {
	mailDoc := maildoc.MailToDoc(mail)
	oid := primitive.NewObjectID()
	userMailDoc := maildoc.UserMailDoc{
		OID:  oid, // 指定 _id
		Mail: &mailDoc,
	}

	if _, err := mgocol.UsersMails().InsertOne(context.Background(), userMailDoc); err != nil {
		return fmt.Errorf("db users mails insert: %w", err)
	}

	// NotifyMail()时才填 mail.ID
	if err := notify.NotifyMail(mail, oid); err != nil {
		// 发邮件成功，但通知失败，按成功处理
		log.Warnf("failed to notify new mail: %v", err)
	}
	return nil
}
