package dbutil

import (
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/DB"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserID = types.UserID

type DocFriendRequest struct {
	UserID  UserID `bson:"userID"`
	FromID  UserID `bson:"fromID"`
	ReqInfo []byte `bson:"reqInfo"`
}

type _FriendRequestsUtil struct {
	userID UserID
}

func FriendRequestsUtil(userID UserID) *_FriendRequestsUtil {
	return &_FriendRequestsUtil{
		userID: userID,
	}
}

func (f *_FriendRequestsUtil) c() *mongo.Collection {
	return DB.MongoDB.Collection("chat.users.friend_requests")
}

func (f *_FriendRequestsUtil) Add(fromID UserID, reqInfo []byte) error {
	selector := bson.M{"userID": f.userID, "fromID": fromID}
	update := bson.M{"$set": DocFriendRequest{
		UserID:  f.userID,
		FromID:  fromID,
		ReqInfo: reqInfo,
	}}
	if _, err := f.c().UpdateOne(context.Background(), selector, update, options.Update().SetUpsert(true)); err != nil {
		return fmt.Errorf("upsert error: %w", err)
	}
	return nil
}

func (f *_FriendRequestsUtil) Query() ([]*DocFriendRequest, error) {
	docs := []*DocFriendRequest{}
	cursor, err := f.c().Find(context.Background(), bson.M{"userID": f.userID}, options.Find().SetLimit(1000))
	if err != nil {
		return nil, fmt.Errorf("find err: %w", err)
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, fmt.Errorf("cursor err: %w", err)
	}
	return docs, nil
}

func (f *_FriendRequestsUtil) Has(fromID UserID) (bool, error) {
	err := f.c().FindOne(context.Background(), f.selector(fromID)).Err()
	if err == nil {
		return true, nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	return false, err
}

func (f *_FriendRequestsUtil) DeleteRequest(fromID UserID) error {
	if _, err := f.c().DeleteOne(context.Background(), f.selector(fromID)); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (f *_FriendRequestsUtil) selector(fromID UserID) bson.M {
	return bson.M{"userID": f.userID, "fromID": fromID}
}

func (f *_FriendRequestsUtil) GetCount() (int64, error) {
	selector := bson.M{"userID": f.userID}
	return f.c().CountDocuments(context.Background(), selector)
}
