package usersdbutil_test

import (
	"Cinder/Chat/rpcproc/logic/types"
	bl "Cinder/Chat/rpcproc/logic/user/internal/blklst/dbutil"
	user "Cinder/Chat/rpcproc/logic/user/internal/dbutil"
	flw "Cinder/Chat/rpcproc/logic/user/internal/follow/following/dbutil"
	frnd "Cinder/Chat/rpcproc/logic/user/internal/friend/dbutil"
	"Cinder/DB"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

const kU1 = "test_user1"
const kU2 = "test_user2"
const kU3 = "test_user3"
const kG1 = "test_group1"
const kG2 = "test_group2"

type UserID = types.UserID

// RemoveUser 删除User记录
func removeUser(userID UserID) error {
	_, err := DB.MongoDB.Collection("chat.users").DeleteOne(context.Background(), bson.M{"userID": userID})
	return err
}

/* Test:
* user: Load()
* blacklist: AddToBlacklist(), RemoveFromBlacklist()
* following: AddFollowID(), DeleteFollowID()
* friend: AddFriend(), DeleteFriend()
 */

func TestLoad(t *testing.T) {
	assert := require.New(t)
	var err error
	assert.Nil(err)

	err = removeUser(kU1)
	assert.Nil(err)

	db := user.UsersUtil(kU1)
	user, errLoad := db.Load()
	assert.Nil(errLoad)
	fmt.Printf("User: %+v\n", user)

	err = removeUser(kU1)
	assert.Nil(err)
}

func TestBlacklist(t *testing.T) {
	assert := require.New(t)
	var err error
	assert.Nil(err)

	_, err = user.UsersUtil(kU1).Load()
	assert.Nil(err)

	db := bl.UsersBlacklistUtil(kU1)
	err = db.AddToBlacklist(kU2)
	assert.Nil(err)
	err = db.AddToBlacklist(kU3)
	assert.Nil(err)
	err = db.RemoveFromBlacklist(kU2)
	assert.Nil(err)
	err = db.RemoveFromBlacklist(kU3)
	assert.Nil(err)
	err = db.RemoveFromBlacklist("no_such_blacklist_member")
	assert.Nil(err)

	err = removeUser(kU1)
	assert.Nil(err)
}

func TestFollowing(t *testing.T) {
	assert := require.New(t)
	var err error
	assert.Nil(err)

	_, err = user.UsersUtil(kU1).Load()
	assert.Nil(err)

	db := flw.UsersUtil(kU1)
	err = db.AddFollowID(kU2)
	assert.Nil(err)
	err = db.AddFollowID(kU3)
	assert.Nil(err)

	err = db.DeleteFollowID(kU2)
	assert.Nil(err)
	err = db.DeleteFollowID(kU3)
	assert.Nil(err)
	err = db.DeleteFollowID("no_such_user")
	assert.Nil(err)

	err = removeUser(kU1)
	assert.Nil(err)
}

func TestFriend(t *testing.T) {
	assert := require.New(t)
	var err error
	assert.Nil(err)

	_, err = user.UsersUtil(kU1).Load()
	assert.Nil(err)

	db := frnd.UsersUtil(kU1)
	err = db.AddFriend(kU2)
	assert.Nil(err)
	err = db.AddFriend(kU3)
	assert.Nil(err)

	err = db.DeleteFriend(kU2)
	assert.Nil(err)
	err = db.DeleteFriend(kU3)
	assert.Nil(err)
	err = db.DeleteFriend("no_such_user")
	assert.Nil(err)

	err = removeUser(kU1)
	assert.Nil(err)
}
