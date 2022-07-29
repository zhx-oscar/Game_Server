package dbutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFollower(t *testing.T) {
	assert := assert.New(t)
	var err error

	db := UsersFollowersUtil("test_user")
	err = db.Add("1")
	assert.Nil(err)
	err = db.Add("2")
	assert.Nil(err)

	followers, err := db.Load()
	assert.Nil(err)
	fmt.Printf("followers: %#v\n", followers)

	for _, flwr := range followers {
		err = db.Remove(flwr.Follower)
		assert.Nil(err)
	}
}
