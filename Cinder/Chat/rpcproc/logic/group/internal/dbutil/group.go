package dbutil

import (
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/DB"
	"context"
	assert "github.com/arl/assertgo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type groupUtil struct {
	groupID GroupID
}

type groupMember struct {
	GroupID  GroupID `bson:"groupID"`
	MemberID UserID  `bson:"memberID"`
}

func GroupUtil(groupID GroupID) *groupUtil {
	return &groupUtil{groupID: groupID}
}

func (u *groupUtil) c() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("chat.groups.members")
}

func (u *groupUtil) InsertMembers(members []UserID) error {
	docs := make([]mongo.WriteModel, 0, len(members))
	for _, memberID := range members {
		docs = append(docs, mongo.NewInsertOneModel().SetDocument(groupMember{
			GroupID:  u.groupID,
			MemberID: memberID,
		}))
	}

	_, err := u.c().BulkWrite(context.Background(), docs, options.BulkWrite().SetOrdered(false))

	return db.SkipDupErr(err)
}

func (u *groupUtil) LoadMemberIDs() ([]UserID, error) {
	var docs []groupMember
	cursor, err := u.c().Find(context.Background(), bson.M{"groupID": u.groupID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, errors.Wrap(err, "find")
	}
	result := make([]UserID, 0, len(docs))
	for _, doc := range docs {
		result = append(result, doc.MemberID)
	}
	return result, nil
}

// DeleteGroup 删除群。
// 即删除所有成员。
func (u *groupUtil) DeleteGroup() error {
	_, err := u.c().DeleteMany(context.Background(), bson.M{"groupID": u.groupID})
	return err
}

func (u *groupUtil) DeleteMember(memberID UserID) error {
	_, err := u.c().DeleteOne(context.Background(), groupMember{
		GroupID:  u.groupID,
		MemberID: memberID,
	})
	return err
}
