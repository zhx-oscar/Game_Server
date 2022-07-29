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

// 因为没有删除，所以禁止该测试
func xTestUserOflnMsgUtil(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = uomu().Insert(nil)
	assert.Nil(err)
	err = uomu().Insert([]*Msg{})
	assert.Nil(err)

	err = uomu().Insert([]*Msg{
		&Msg{
			From:       "from_test",
			FromNick:   "from_nixk",
			SendTime:   1234567,
			MsgContent: []byte("111111"),
		},
		&Msg{},
	})
	assert.Nil(err)
}
