package db

import (
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/internal/bcshard"
	"Cinder/Mail/rpcproc/handler/internal/maildoc"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// mail.broadcast_mails 文档中的state 和 time 字段
type StateAndTime struct {
	State maildoc.MailState `bson:"state"`
	Mail  struct {
		SendTime   time.Time `bson:"sendTime"`
		ExpireTime time.Time `bson:"expireTime"`
	} `bson:"mail"`
}

// BroadcastMails 对应 mail.broadcast_mails 集合操作
type BroadcastMails struct {
	originalID primitive.ObjectID
}

func NewBroadcastMails(originalID primitive.ObjectID) *BroadcastMails {
	return &BroadcastMails{
		originalID: originalID,
	}
}

func (b *BroadcastMails) LoadStateAndTime() (StateAndTime, error) {
	query := bson.M{"shard": bcshard.GetRandBcShardID(), "originalID": b.originalID}
	selector := bson.M{"state": 1, "mail.sendTime": 1, "mail.expireTime": 1}
	var doc StateAndTime
	if rv := mgocol.BroadcastMails().FindOne(context.Background(), query, options.FindOne().SetProjection(selector)); rv.Err() != nil {
		return doc, rv.Err()
	} else {
		if err := rv.Decode(&doc); err != nil {
			return doc, err
		}
	}

	return doc, nil
}
