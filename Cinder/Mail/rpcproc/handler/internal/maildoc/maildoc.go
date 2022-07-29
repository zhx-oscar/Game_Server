// maildoc 包定义 mongodb 中邮件文档。
package maildoc

import (
	"Cinder/Mail/rpcproc/handler/internal/bcshard"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserMailDoc 用户邮件文档
type UserMailDoc struct {
	OID primitive.ObjectID `bson:"_id"`

	Mail *Mail `bson:"mail"`
}

// BcMailDoc 是广播邮件文档。
// 一分邮件分片保存多个，分散读压力
type BcMailDoc struct {
	Shard      bcshard.ShardID    `bson:"shard"`      // 分片号
	OriginalID primitive.ObjectID `bson:"originalID"` // 原始ID

	Mail *Mail `bson:"mail"`
}

// Mail 是mongodb中普通邮件和系统邮件文档的主体。
// 对于普通邮件，需添加 ObjectId。
// 对于系统广播邮件，需添加分片号，原始 ObjectId.
type Mail struct {
	From     string `bson:"from"`
	FromNick string `bson:"fromNick"`
	To       string `bson:"to"`
	ToNick   string `bson:"toNick"`
	Title    string `bson:"title"`
	Body     string `bson:"body"`

	State *MailState `bson:"state"`

	Attachments []*Attachment `bson:"attachments"`

	SendTime   time.Time `bson:"sendTime"`
	ExpireTime time.Time `bson:"expireTime"`
}

type MailState struct {
	IsRead                bool   `bson:"isRead"`
	IsAttachmentsReceived bool   `bson:"isAttachmentsReceived"`
	ExtData               []byte `bson:"extData"`
}

type Attachment struct {
	ItemID uint32 `bson:"itemID"`
	Count  uint32 `bson:"count"`
	Data   []byte `bson:"data"`
}
