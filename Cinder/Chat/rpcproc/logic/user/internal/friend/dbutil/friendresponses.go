package dbutil

import (
	"Cinder/DB"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocFriendResponse struct {
	UserID      UserID `bson:"userID"`
	ResponderID UserID `bson:"responderID"`
	OK          bool   `bson:"ok"`
}

type _FriendResponsesUtil struct {
	userID UserID
}

func FriendResponsesUtil(userID UserID) *_FriendResponsesUtil {
	return &_FriendResponsesUtil{
		userID: userID,
	}
}

func (f *_FriendResponsesUtil) c() *mongo.Collection {
	return DB.MongoDB.Collection("chat.users.friend_responses")
}

func (f *_FriendResponsesUtil) Add(responderID UserID, ok bool) error {
	selector := bson.M{"userID": f.userID, "responderID": responderID}
	update := bson.M{"$set": DocFriendResponse{
		UserID:      f.userID,
		ResponderID: responderID,
		OK:          ok,
	}}
	if _, err := f.c().UpdateOne(context.Background(), selector, update, options.Update().SetUpsert(true)); err != nil {
		return fmt.Errorf("upsert error: %w", err)
	}
	return nil
}

func (f *_FriendResponsesUtil) Pop() ([]*DocFriendResponse, error) {
	docs := []*DocFriendResponse{}
	selector := bson.M{"userID": f.userID}
	cursor, err := f.c().Find(context.Background(), selector, options.Find().SetLimit(1000))
	if err != nil {
		return nil, fmt.Errorf("find err: %w", err)
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, fmt.Errorf("cursor err: %w", err)
	}

	if len(docs) == 0 {
		return docs, nil
	}

	// 只取1000，多余的也会被删
	if _, err := f.c().DeleteMany(context.Background(), selector); err != nil {
		return nil, fmt.Errorf("delete err: %w", err)
	}

	return docs, nil
}
