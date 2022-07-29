package dbutil

import (
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/DB"
	"context"
	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserID = types.UserID

// _UsersBlacklistUtil 操作 chat.users 集合中的blacklist数组
type _UsersBlacklistUtil struct {
	userID UserID
}

func UsersBlacklistUtil(userID UserID) *_UsersBlacklistUtil {
	return &_UsersBlacklistUtil{
		userID: userID,
	}
}

func (u *_UsersBlacklistUtil) c() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("chat.users")
}

func (u *_UsersBlacklistUtil) AddToBlacklist(userID UserID) error {
	selector := bson.M{"userID": u.userID}
	update := bson.M{"$addToSet": bson.M{"blacklist": userID}}
	// TODO: 限制大小, 虽然内存中已有限制
	_, err := u.c().UpdateOne(context.Background(), selector, update)
	return err
}

func (u *_UsersBlacklistUtil) RemoveFromBlacklist(userID UserID) error {
	selector := bson.M{"userID": u.userID}
	update := bson.M{"$pull": bson.M{"blacklist": userID}}
	_, err := u.c().UpdateOne(context.Background(), selector, update)
	return err
}
