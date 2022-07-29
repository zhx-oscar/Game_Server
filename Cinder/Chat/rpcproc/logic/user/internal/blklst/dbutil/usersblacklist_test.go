package dbutil

import (
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/DB"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestBlacklist(t *testing.T) {
	var err error
	assert := assert.New(t)
	const kUser = "test_user_blacklist"
	_, err = DB.MongoDB.Collection("chat.users").InsertOne(context.Background(), bson.M{"userID": kUser})
	err = db.SkipDupErr(err)
	assert.Nil(err)

	db := UsersBlacklistUtil(kUser)
	const kBlocked = "blocked_user0"
	err = db.AddToBlacklist(kBlocked)
	assert.Nil(err)
	err = db.RemoveFromBlacklist(kBlocked)
	assert.Nil(err)

	_, err = DB.MongoDB.Collection("chat.users").DeleteOne(context.Background(), bson.M{"userID": kUser})
	assert.Nil(err)
}
