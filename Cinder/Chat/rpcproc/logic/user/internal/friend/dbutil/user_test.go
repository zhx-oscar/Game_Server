package dbutil

import (
	"Cinder/DB"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetFriendCount(t *testing.T) {
	assert := require.New(t)
	var err error
	var cnt int
	const kUser = "user_friend_tester"
	db := UsersUtil(kUser)

	col := DB.MongoDB.Collection("chat.users")
	selector := bson.M{"userID": kUser}
	_, err = col.DeleteMany(context.Background(), selector)
	assert.NoError(err)

	cnt, err = db.GetFriendCount()
	assert.NoError(err)
	assert.Equal(0, cnt)

	// 必须先有记录，才能 AddFriend()
	_, err = col.InsertOne(context.Background(), bson.M{"userID": kUser})
	assert.NoError(err)
	defer func() {
		_, err = col.DeleteMany(context.Background(), selector)
		assert.NoError(err)
	}()

	cnt, err = db.GetFriendCount()
	assert.NoError(err)
	assert.Equal(0, cnt)

	err = db.AddFriend("f1")
	assert.NoError(err)
	cnt, err = db.GetFriendCount()
	assert.NoError(err)
	assert.Equal(1, cnt)

	err = db.AddFriend("f2")
	assert.NoError(err)
	err = db.AddFriend("f3")
	assert.NoError(err)
	cnt, err = db.GetFriendCount()
	assert.NoError(err)
	assert.Equal(3, cnt)

	err = db.DeleteFriend("f1")
	assert.NoError(err)
	err = db.DeleteFriend("f2")
	assert.NoError(err)
	err = db.DeleteFriend("f3")
	assert.NoError(err)

	cnt, err = db.GetFriendCount()
	assert.NoError(err)
	assert.Equal(0, cnt)
}
