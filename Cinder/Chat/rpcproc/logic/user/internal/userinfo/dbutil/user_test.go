package dbutil

import (
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/DB"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpdate(t *testing.T) {
	assert := require.New(t)
	const kUser = "test_user_userinfo"
	_, err := DB.MongoDB.Collection("chat.users").InsertOne(context.Background(), bson.M{"userID": kUser})
	err = db.SkipDupErr(err)
	assert.NoError(err)

	db := UsersUtil(kUser)
	err = db.UpdateNick("nick")
	assert.NoError(err)
	err = db.UpdateActiveData([]byte("!!!data"))
	assert.NoError(err)
	err = db.UpdatePassiveData([]byte("!!!passiave data!!!"))
	assert.NoError(err)

	DB.MongoDB.Collection("chat.users").DeleteOne(context.Background(), bson.M{"userID": kUser})
}
