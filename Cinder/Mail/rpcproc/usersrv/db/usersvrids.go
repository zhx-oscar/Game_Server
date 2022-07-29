package db

import (
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/Mail/mgocol"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type _Doc struct {
	SrvID string `bson:"srvID"`
}

type _UserSrvIDs struct {
}

var userSrvIDs = &_UserSrvIDs{}

func GetUserSrvIDs() *_UserSrvIDs {
	return userSrvIDs
}

func (l *_UserSrvIDs) c() *mongo.Collection {
	return mgocol.UserSrvIDs()
}

func (l *_UserSrvIDs) Insert(srvID string) error {
	_, err := l.c().InsertOne(context.Background(), _Doc{
		SrvID: srvID,
	})
	return db.SkipDupErr(err)
}

func (l *_UserSrvIDs) Load() (map[string]bool, error) {
	docs := []_Doc{}
	cursor, err := l.c().Find(context.Background(), bson.M{}, options.Find().SetLimit(65535))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, err
	}

	result := make(map[string]bool, len(docs))
	for _, doc := range docs {
		result[doc.SrvID] = true
	}
	return result, nil
}

// RemoveIDs remove user srv ids
func (l *_UserSrvIDs) RemoveIDs(ids []string) error {
	var docs []mongo.WriteModel
	for _, id := range ids {
		docs = append(docs, mongo.NewDeleteOneModel().SetFilter(&_Doc{SrvID: id}))
	}

	_, err := l.c().BulkWrite(context.Background(), docs, options.BulkWrite().SetOrdered(false))

	return err
}
