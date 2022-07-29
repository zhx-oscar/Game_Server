package dbutil

import (
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/Chat/rpcproc/logic/user/internal/dbdoc"
	"Cinder/DB"
	"context"
	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type _UsersGroupsUtil struct {
	groupID GroupID
}

func UsersGroupsUtil(groupID GroupID) *_UsersGroupsUtil {
	return &_UsersGroupsUtil{groupID: groupID}
}

func (u *_UsersGroupsUtil) c() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("chat.users.groups")
}

// AddGroupToUsers 为一批用户添加群记录，忽略重复。
// 用于群创建，也可用于单个用户加入群。
func (u *_UsersGroupsUtil) AddGroupToUsers(userIDs []UserID, sequenceID SequenceID) error {
	docs := make([]mongo.WriteModel, 0, len(userIDs))
	for _, userID := range userIDs {
		docs = append(docs, mongo.NewInsertOneModel().SetDocument(
			dbdoc.UsersGroupsDoc{
				UserID:     userID,
				GroupID:    u.groupID,
				SequenceID: sequenceID,
			}))
	}
	_, err := u.c().BulkWrite(context.Background(), docs, options.BulkWrite().SetOrdered(false))
	return db.SkipDupErr(err)
}

// DeleteGroupFromUsers 为一批用户删除群。
// 用于群删除，也用于单个用户退群。
func (u *_UsersGroupsUtil) DeleteGroupFromUsers(userIDs []UserID) error {
	var docs []mongo.WriteModel
	for _, userID := range userIDs {
		docs = append(docs, mongo.NewDeleteOneModel().SetFilter(
			bson.M{"userID": userID, "groupID": u.groupID},
		))
	}
	_, err := u.c().BulkWrite(context.Background(), docs, options.BulkWrite().SetOrdered(false))
	return err
}
