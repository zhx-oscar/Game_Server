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

func (u *_UsersUtil) UpdateNick(nick string) error {
	update := bson.M{"$set": bson.M{"nick": nick}}
	return u.update(update)
}

func (u *_UsersUtil) UpdateActiveData(activeData []byte) error {
	update := bson.M{"$set": bson.M{"activeData": activeData}}
	return u.update(update)
}

func (u *_UsersUtil) UpdatePassiveData(passiveData []byte) error {
	update := bson.M{"$set": bson.M{"passiveData": passiveData}}
	return u.update(update)
}

func (u *_UsersUtil) update(update bson.M) error {
	selector := bson.M{"userID": u.userID}
	_, err := u.c().UpdateOne(context.Background(), selector, update)
	return err
}
