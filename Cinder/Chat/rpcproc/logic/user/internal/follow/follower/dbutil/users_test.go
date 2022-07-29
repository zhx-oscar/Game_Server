package dbutil

import (
	"Cinder/DB"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestIncDec(t *testing.T) {
	assert := require.New(t)
	const kUser = "test_user"
	db := UsersUtil(kUser)
	var err error

	_, _ = DB.MongoDB.Collection("chat.users").InsertOne(context.Background(), bson.M{"userID": kUser})

	err = db.IncreaseFollowerNumber()
	assert.NoError(err)
	err = db.DecreaseFollowerNumber()
	assert.NoError(err)

	_, err = DB.MongoDB.Collection("chat.users").DeleteOne(context.Background(), bson.M{"userID": kUser})
	assert.NoError(err)
}
