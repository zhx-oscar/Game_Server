package dbutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	assert := assert.New(t)
	var err error
	const kUser = "user_friend_response"
	db := FriendResponsesUtil(kUser)
	err = db.Add("fromID", true)
	assert.Nil(err)
	reqs, errPop := db.Pop()
	assert.Nil(errPop)
	fmt.Printf("resps: %+v\n", reqs)
	reqs, errPop = db.Pop()
	assert.Nil(errPop)
	fmt.Printf("resps: %+v\n", reqs)
}
