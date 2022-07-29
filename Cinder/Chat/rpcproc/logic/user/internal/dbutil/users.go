package dbutil

import (
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/DB"
	"context"
	"errors"
	"time"

	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// chat.users 中的 User 文档
type UserDoc struct {
	UserID      UserID `bson:"userID"`
	Nick        string `bson:"nick"`
	ActiveData  []byte `bson:"activeData"`
	PassiveData []byte `bson:"passiveData"`

	FriendIDs []UserID `bson:"friendIDs"`
	FollowIDs []UserID `bson:"followIDs"`
	Blacklist []UserID `bson:"blacklist"`

	OfflineTime    time.Time `bson:"offlineTime"`
	FollowerNumber int       `bson:"followerNumber"`
}

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

func (u *_UsersUtil) Load() (*UserDoc, error) {
	doc := &UserDoc{}
	if rv := u.c().FindOne(context.Background(), bson.M{"userID": u.userID}); rv.Err() != nil {
		if !errors.Is(rv.Err(), mongo.ErrNoDocuments) {
			return nil, rv.Err()
		}
	} else {
		if err := rv.Decode(doc); err != nil {
			return nil, err
		}
		return doc, nil
	}
	// DB 插入，防止后续更新找不到
	if err := u.insertNew(); err != nil {
		return nil, err
	}
	return &UserDoc{UserID: u.userID}, nil
}

func (u *_UsersUtil) insertNew() error {
	_, err := u.c().InsertOne(context.Background(), &UserDoc{
		UserID: u.userID,
		// Must not nil orelse:
		// Cannot apply $addToSet to non-array field. Field named 'blacklist' has non-array type null
		FriendIDs:   []UserID{},
		FollowIDs:   []UserID{},
		Blacklist:   []UserID{},
		OfflineTime: time.Now(),
	})
	return db.SkipDupErr(err)
}

// SaveOfflineTime DB写 OfflineTime
func (u *_UsersUtil) SaveOfflineTime() error {
	selector := bson.M{"userID": u.userID}
	update := bson.M{"$set": bson.M{"offlineTime": time.Now()}}
	_, err := u.c().UpdateOne(context.Background(), selector, update)
	return err
}
