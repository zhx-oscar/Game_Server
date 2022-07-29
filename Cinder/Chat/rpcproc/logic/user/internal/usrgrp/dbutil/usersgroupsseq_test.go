package dbutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const kUser = "test_user"
const kGroup = "test_group"

func TestUpdateGroupToSeq(t *testing.T) {
	assert := assert.New(t)
	db := UsersGroupsSeqUtil(kUser)
	var err error

	err = UsersGroupsUtil("g1").AddGroupToUsers([]UserID{kUser}, 0)
	assert.Nil(err)
	err = UsersGroupsUtil("g2").AddGroupToUsers([]UserID{kUser}, 0)
	assert.Nil(err)

	err = db.UpdateGroupToSeq(map[GroupID]SequenceID{
		"g1":            1111,
		"g2":            2222,
		"no_such_group": 1234,
	})
	assert.Nil(err)

	g2s, errG2s := db.LoadGroupToSeq()
	assert.Nil(errG2s)
	fmt.Printf("g2s: %v", g2s)

	err = UsersGroupsUtil("g1").DeleteGroupFromUsers([]UserID{kUser})
	assert.Nil(err)
	err = UsersGroupsUtil("g2").DeleteGroupFromUsers([]UserID{kUser})
	assert.Nil(err)
}
