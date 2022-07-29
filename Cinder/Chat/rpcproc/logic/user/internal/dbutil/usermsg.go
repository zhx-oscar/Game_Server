package dbutil

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/user/internal/dbdoc"
	"Cinder/DB"
	"context"
	assert "github.com/arl/assertgo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User offline message util.
type _UserOflnMsgUtil struct {
	userID UserID
}

func UserOflnMsgUtil(userID UserID) *_UserOflnMsgUtil {
	return &_UserOflnMsgUtil{
		userID: userID,
	}
}

func (u *_UserOflnMsgUtil) c() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("chat.users.offline_messages")
}

// Load 加载离线消息，最多1000条，返回按 sendTime 从小到大排列
func (u *_UserOflnMsgUtil) Load() ([]*chatapi.ChatMessage, error) {
	var docs []dbdoc.UsersOfflineMessagesDoc
	cursor, err := u.c().Find(context.Background(), bson.M{"userID": u.userID}, options.Find().SetLimit(1000).SetSort(bson.M{"sendTime": -1}))
	if err != nil {
		return nil, errors.Wrap(err, "find err")
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, errors.Wrap(err, "cursor err")
	}
	// 需要反序
	l := len(docs)
	res := make([]*chatapi.ChatMessage, l)
	for i := 0; i < l; i++ {
		msg := docs[l-i-1]
		res[i] = &chatapi.ChatMessage{
			From:       msg.FromID,
			FromNick:   msg.FromNick,
			FromData:   msg.FromData,
			SendTime:   msg.SendTime,
			MsgContent: msg.MsgContent,
		}
	}
	return res, nil
}

func (u *_UserOflnMsgUtil) Remove() error {
	_, err := u.c().DeleteMany(context.Background(), bson.M{"userID": u.userID})
	return err
}
