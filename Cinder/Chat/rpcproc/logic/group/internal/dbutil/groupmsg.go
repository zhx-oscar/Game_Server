package dbutil

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/DB"
	"context"
	"fmt"
	assert "github.com/arl/assertgo"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Group offline message util.
type _GroupsMessagesUtil struct {
	groupID GroupID
}

type _GroupMsg struct {
	GroupID    GroupID    `bson:"groupID"`
	FromID     string     `bson:"fromID"`
	FromNick   string     `bson:"fromNick"`
	FromData   []byte     `bson:"fromData"`
	SendTime   int64      `bson:"sendTime"`
	MsgContent []byte     `bson:"msgContent"`
	SequenceID SequenceID `bson:"sequenceID"`
}

func GroupsMessagessUtil(groupID GroupID) *_GroupsMessagesUtil {
	return &_GroupsMessagesUtil{
		groupID: groupID,
	}
}

func (u *_GroupsMessagesUtil) c() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("chat.groups.messages")
}

func (g *_GroupsMessagesUtil) Load() (map[SequenceID]*chatapi.ChatMessage, error) {
	var docs []_GroupMsg
	cursor, err := g.c().Find(context.Background(), bson.M{"groupID": g.groupID}, options.Find().SetLimit(1000).SetSort(bson.M{"sequenceID": -1}))
	if err != nil {
		return nil, errors.Wrap(err, "find err")
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, errors.Wrap(err, "cursor err")
	}

	ret := make(map[SequenceID]*chatapi.ChatMessage)
	for _, rec := range docs {
		ret[rec.SequenceID] = &chatapi.ChatMessage{
			From:       rec.FromID,
			FromNick:   rec.FromNick,
			FromData:   rec.FromData,
			SendTime:   rec.SendTime,
			MsgContent: rec.MsgContent,
		}
	}
	return ret, nil
}

// Save 保存离线消息, 保存序号为[startSeq, endSeq]
// 允许 msgs 为nil, 或空。
// 最多仅保存最后 10000 个.
// seq 接近最大值时返回错误
func (g *_GroupsMessagesUtil) Insert(msgs map[SequenceID]*chatapi.ChatMessage, startSeq, endSeq SequenceID) error {
	if endSeq < startSeq {
		log.Warnf("start(%d) is larger than end(%d)", startSeq, endSeq)
		return nil
	}

	// 如不限制，Insert(nil, ^0, ^0), Insert(nil, 0, ^0) 会死循环。
	if endSeq > ^SequenceID(0)-100 {
		return fmt.Errorf("sequence ID is too large: %v", endSeq)
	}

	const kMaxCount = 10000
	cnt := endSeq - startSeq + 1
	if cnt > kMaxCount {
		log.Warnf("insert count(%d..%d = %d) is larger than %d", startSeq, endSeq, cnt, kMaxCount)
		cnt = kMaxCount
	}

	docs := make([]mongo.WriteModel, 0, cnt)
	for seq := endSeq - cnt + 1; seq <= endSeq; seq++ {
		msg, _ := msgs[seq]
		if msg == nil {
			continue
		}

		docs = append(docs, mongo.NewInsertOneModel().SetDocument(_GroupMsg{
			GroupID:    g.groupID,
			FromID:     msg.From,
			FromNick:   msg.FromNick,
			FromData:   msg.FromData,
			SendTime:   msg.SendTime,
			MsgContent: msg.MsgContent,
			SequenceID: seq,
		}))
	}
	if len(docs) == 0 {
		return nil
	}
	_, err := g.c().BulkWrite(context.Background(), docs, options.BulkWrite().SetOrdered(false))
	return db.SkipDupErr(err)
}
