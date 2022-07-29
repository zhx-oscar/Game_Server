package dbutil

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user/internal/dbdoc"
	"Cinder/DB"
	"context"

	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserID = types.UserID

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

func (u *_UserOflnMsgUtil) Insert(msgs []*chatapi.ChatMessage) error {
	if len(msgs) == 0 {
		return nil // Got 0 operations.
	}
	docs := make([]interface{}, 0, len(msgs))
	for _, msg := range msgs {
		docs = append(docs, dbdoc.UsersOfflineMessagesDoc{
			UserID:     u.userID,
			FromID:     msg.From,
			FromNick:   msg.FromNick,
			FromData:   msg.FromData,
			SendTime:   msg.SendTime,
			MsgContent: msg.MsgContent,
		})
	}
	_, err := u.c().InsertMany(context.Background(), docs) // 没有唯一键
	return err
}
