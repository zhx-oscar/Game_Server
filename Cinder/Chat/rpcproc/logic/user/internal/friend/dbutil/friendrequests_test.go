package dbutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	assert := require.New(t)
	var err error
	const kUser = "user_friend_request"
	db := FriendRequestsUtil(kUser)
	err = db.Add("fromID", []byte("reqInfo"))
	assert.Nil(err)
	reqs, errPop := db.Query()
	assert.Nil(errPop)
	fmt.Printf("req: %+v\n", reqs)
	reqs, errPop = db.Query()
	assert.Nil(errPop)
	fmt.Printf("req: %+v\n", reqs)
}

func TestHas(t *testing.T) {
	assert := require.New(t)
	var err error
	var has bool
	const kUser = "user_friend_has_tester"
	db := FriendRequestsUtil(kUser)
	has, err = db.Has("no_such_user")
	assert.NoError(err)
	assert.False(has)

	err = db.Add("fromID", []byte("reqInfo"))
	assert.NoError(err)
	has, err = db.Has("fromID")
	assert.NoError(err)
	assert.True(has)
}

func TestGetCount(t *testing.T) {
	assert := require.New(t)
	var cnt int64
	var err error
	const kUser = "user_test_getcount"
	const kFrom1 = "user_test_getcount_from1"
	const kFrom2 = "user_test_getcount_from2"

	db := FriendRequestsUtil(kUser)
	err = db.DeleteRequest(kFrom1)
	assert.NoError(err)
	err = db.DeleteRequest(kFrom2)
	assert.NoError(err)

	cnt, err = db.GetCount()
	assert.NoError(err)
	assert.Empty(cnt)
	err = db.Add(kFrom1, []byte("abcde"))
	assert.NoError(err)
	cnt, err = db.GetCount()
	assert.NoError(err)
	assert.Equal(int64(1), cnt)
	err = db.Add(kFrom2, []byte("abcde"))
	assert.NoError(err)
	cnt, err = db.GetCount()
	assert.NoError(err)
	assert.Equal(int64(2), cnt)

	err = db.DeleteRequest(kFrom1)
	assert.NoError(err)
	err = db.DeleteRequest(kFrom2)
	assert.NoError(err)
}
