package dbutil

import (
	"Cinder/DB"
	"context"
	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// _UsersUtil 操作 chat.users 集合
type _UsersUtil struct {
	userID UserID
}

func UsersUtil(userID UserID) *_UsersUtil {
	return &_UsersUtil{
		userID: userID,
	}
}

func (u *_UsersUtil) c() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("chat.users")
}

func (u *_UsersUtil) IncreaseFollowerNumber() error {
	return u.incFollowerNumber(1)
}

func (u *_UsersUtil) DecreaseFollowerNumber() error {
	return u.incFollowerNumber(-1)
}

func (u *_UsersUtil) incFollowerNumber(delta int) error {
	selector := bson.M{"userID": u.userID}
	update := bson.M{"$inc": bson.M{"followerNumber": delta}}
	_, err := u.c().UpdateOne(context.Background(), selector, update)
	return err
}
