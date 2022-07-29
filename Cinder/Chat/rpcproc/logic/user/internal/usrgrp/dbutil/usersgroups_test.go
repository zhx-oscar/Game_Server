package dbutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersGroupUtil(t *testing.T) {
	assert := assert.New(t)
	var err error

	db := UsersGroupsUtil("test_group")
	err = db.AddGroupToUsers([]UserID{"u1", "u2", "u3", "u4"}, 123)
	assert.Nil(err)

	err = db.DeleteGroupFromUsers([]UserID{"u1", "u2"})
	assert.Nil(err)
	err = db.DeleteGroupFromUsers([]UserID{"u1", "u2", "u3", "u4"})
	assert.Nil(err)
}
