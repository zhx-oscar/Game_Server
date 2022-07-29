package dbutil

import (
	"Cinder/Chat/chatapi"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Msg = chatapi.ChatMessage

func uomu() *_UserOflnMsgUtil {
	return UserOflnMsgUtil("test_user")
}

func TestUserOflnMsgUtil(t *testing.T) {
	var err error
	assert := assert.New(t)

	_, errLoad := uomu().Load()
	assert.Nil(errLoad)
	err = uomu().Remove()
	assert.Nil(err)
}
