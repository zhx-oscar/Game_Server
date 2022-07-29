package dbutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDB() *groupUtil {
	return GroupUtil("test_group")
}

func TestInsertDeleteMember(t *testing.T) {
	var err error
	err = getDB().InsertMembers([]UserID{"test_user"})
	assert.Nil(t, err)
	err = getDB().DeleteMember("test_user")
	assert.Nil(t, err)
}

func TestDeleteGroup(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = getDB().DeleteGroup()
	assert.Nil(err)
	err = getDB().InsertMembers([]UserID{"test_user1", "test_user2"})
	assert.Nil(err)
	err = getDB().DeleteGroup()
	assert.Nil(err)
}

func TestLoadMemberIDs(t *testing.T) {
	var err error
	assert := assert.New(t)
	db := getDB()
	err = db.InsertMembers([]UserID{"u1", "u2", "u3"})
	assert.Nil(err)

	ids, errIDs := db.LoadMemberIDs()
	assert.Nil(errIDs)
	assert.Equal(3, len(ids))

	err = db.DeleteGroup()
	assert.Nil(err)
}
