// bcmail 发送系统广播邮件
package bcmail

import (
	"Cinder/Mail/mailapi/types"
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/internal/bcshard"
	"Cinder/Mail/rpcproc/handler/internal/maildoc"
	"Cinder/Mail/rpcproc/handler/internal/notify"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/cihub/seelog"
)

// Broadcast 广播邮件
func Broadcast(mail types.Mail) error {
	mail.IsBroadcast = true

	// 复制多份保存，分散读压力
	originalID := primitive.NewObjectID()
	mailDoc := maildoc.MailToDoc(mail)
	docs := make([]interface{}, 0, bcshard.BroadcastShardCount)
	for shard := bcshard.ShardID(0); shard < bcshard.BroadcastShardCount; shard++ {
		docs = append(docs, &maildoc.BcMailDoc{
			Shard:      shard,
			OriginalID: originalID,
			Mail:       &mailDoc,
		})
	} // for

	if _, err := mgocol.BroadcastMails().InsertMany(context.Background(), docs); err != nil {
		return fmt.Errorf("db broadcast mails insert: %w", err)
	}
	// NotifyBroadcastMail() 中会填 mail.ID
	if err := notify.NotifyBroadcastMail(mail, originalID); err != nil {
		// 发邮件成功，但通知失败，按成功处理
		log.Warnf("failed to notify new broadcast mail: %v", err)
	}
	return nil
}
