package db

import (
	"Cinder/Mail/mgocol"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/stretchr/testify/require"
)

func TestUpsert(t *testing.T) {
	assert := require.New(t)
	var err error

	err = mgocol.Users().Remove(bson.M{"userID": "test_user"})
	if !errors.Is(err, mgo.ErrNotFound) {
		assert.NoError(err)
	}

	err = GetUsers("test_user").Upsert("srvID_abce")
	assert.NoError(err)
}
