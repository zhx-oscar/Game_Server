package dbutil

import (
	"Cinder/DB"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLoad(t *testing.T) {
	const kUser = "test_user3"
	db := UsersUtil(kUser)
	user, err := db.Load()
	assert.Nil(t, err)
	fmt.Printf("user: %+v\n", user)
	_, errRm := DB.MongoDB.Collection("chat.users").DeleteOne(context.Background(), bson.M{"userID": kUser})
	assert.Nil(t, errRm)
}

func TestSaveOfflineTime(t *testing.T) {
	assert := require.New(t)
	const kUser = "test_user"
	db := UsersUtil(kUser)
	var err error
	_, err = db.Load()
	assert.NoError(err)
	err = db.SaveOfflineTime()
	assert.NoError(err)
	doc, errDoc := db.Load()
	assert.NoError(errDoc)
	assert.InDelta(time.Now().Unix(), doc.OfflineTime.Unix(), 2)

	_, errRm := DB.MongoDB.Collection("chat.users").DeleteOne(context.Background(), bson.M{"userID": kUser})
	assert.NoError(errRm)
}
