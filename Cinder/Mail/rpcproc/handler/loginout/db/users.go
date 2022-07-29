package db

import (
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/userid"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type _UserDoc struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID userid.UserID      `bson:"userID"`
	SrvID  string             `bson:"srvID"`
}

type _Users struct {
	userID userid.UserID
}

func GetUsers(userID userid.UserID) *_Users {
	return &_Users{
		userID: userID,
	}
}

func (l *_Users) Upsert(srvID string) error {
	selector := bson.M{"userID": l.userID}
	update := bson.M{"$set": bson.M{"srvID": srvID}}
	_, err := mgocol.Users().UpdateOne(context.Background(), selector, update, options.Update().SetUpsert(true))
	return err
}
