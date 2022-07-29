package dbutil

import (
	"Cinder/Chat/rpcproc/logic/user/internal/dbdoc"
	"Cinder/DB"
	"context"
	"fmt"
	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type _UsersGroupsSeqUtil struct {
	userID UserID
}

func UsersGroupsSeqUtil(userID UserID) *_UsersGroupsSeqUtil {
	return &_UsersGroupsSeqUtil{userID: userID}
}

func (u *_UsersGroupsSeqUtil) c() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("chat.users.groups")
}

func (u *_UsersGroupsSeqUtil) LoadGroupToSeq() (map[GroupID]SequenceID, error) {
	docs := []dbdoc.UsersGroupsDoc{}
	cursor, err := u.c().Find(context.Background(), bson.M{"userID": u.userID})
	if err != nil {
		return nil, fmt.Errorf("find err: %w", err)
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, fmt.Errorf("cursor err: %w", err)
	}

	result := make(map[GroupID]SequenceID)
	for _, doc := range docs {
		result[doc.GroupID] = doc.SequenceID
	}
	return result, nil
}

func (u *_UsersGroupsSeqUtil) UpdateGroupToSeq(g2s map[GroupID]SequenceID) error {
	docs := make([]mongo.WriteModel, 0, len(g2s)*2)
	for groupID, seqID := range g2s {
		docs = append(docs, mongo.NewUpdateOneModel().SetFilter(dbdoc.UsersGroupsDocKey{
			UserID:  u.userID,
			GroupID: groupID,
		}).SetUpdate(bson.M{
			"$set": bson.M{"sequenceID": seqID},
		}))
	}
	_, err := u.c().BulkWrite(context.Background(), docs, options.BulkWrite().SetOrdered(false))
	return err
}
