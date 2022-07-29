package dbutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddFollowID(t *testing.T) {
	assert := assert.New(t)
	db := UsersUtil("test_user")
	err := db.AddFollowID("bigman2")
	assert.Nil(err)
}
