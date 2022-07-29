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

// _UsersUtil 对应 chat.users 集合的操作
type _UsersUtil struct {
	userID UserID // 主人ID，粉丝的关注对象
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

func (u *_UsersUtil) AddFollowID(id UserID) error {
	selector := bson.M{"userID": u.userID}
	update := bson.M{"$addToSet": bson.M{"followIDs": id}}
	_, err := u.c().UpdateOne(context.Background(), selector, update)
	return err
}

func (u *_UsersUtil) DeleteFollowID(id UserID) error {
	selector := bson.M{"userID": u.userID}
	update := bson.M{"$pull": bson.M{"followIDs": id}}
	_, err := u.c().UpdateOne(context.Background(), selector, update)
	return err
}
