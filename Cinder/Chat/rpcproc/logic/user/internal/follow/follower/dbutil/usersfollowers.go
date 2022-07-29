package dbutil

import (
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/DB"
	"context"
	"fmt"
	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserID = types.UserID

type _UserFollower struct {
	UserID   UserID `bson:"userID"`
	Follower UserID `bson:"follower"`
}

type _Follower struct {
	OID      primitive.ObjectID `bson:"_id"`
	Follower UserID             `bson:"follower"`
}

// _UsersFollowersUtil 对应 chat.users.followers 集合的操作
type _UsersFollowersUtil struct {
	userID UserID // 主人ID，粉丝的关注对象
}

func UsersFollowersUtil(userID UserID) *_UsersFollowersUtil {
	return &_UsersFollowersUtil{
		userID: userID,
	}
}
func (u *_UsersFollowersUtil) c() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("chat.users.followers")
}

func (u *_UsersFollowersUtil) Add(followerID UserID) error {
	_, err := u.c().InsertOne(context.Background(), _UserFollower{
		UserID:   u.userID,
		Follower: followerID,
	})
	return db.SkipDupErr(err)
}

func (u *_UsersFollowersUtil) Remove(followerID UserID) error {
	_, err := u.c().DeleteOne(context.Background(), _UserFollower{
		UserID:   u.userID,
		Follower: followerID,
	})
	return err
}

func (u *_UsersFollowersUtil) Load() ([]_Follower, error) {
	docs := []_Follower{}
	cursor, err := u.c().Find(context.Background(), bson.M{"userID": u.userID}, options.Find().SetLimit(1024*8).SetProjection(bson.M{"follower": 1, "_id": 1}))
	if err != nil {
		return nil, fmt.Errorf("find error: %w", err)
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return docs, nil
}
