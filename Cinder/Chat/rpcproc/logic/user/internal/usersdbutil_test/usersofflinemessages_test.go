package usersdbutil_test

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/user/internal/dbutil"
	offutil "Cinder/Chat/rpcproc/logic/user/internal/oflmsg/dbutil"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type _ChatMsg = chatapi.ChatMessage

/*
 dbutil: Load(), Remove()
 offutil: Insert()
*/

func TestOfflineMessages(t *testing.T) {
	assert := assert.New(t)
	var err error
	assert.Nil(err)

	msgs := []*_ChatMsg{
		&_ChatMsg{
			From:       "fromID",
			FromNick:   "fromNick",
			FromData:   []byte("fromData"),
			SendTime:   1234,
			MsgContent: []byte("content"),
		},
		&_ChatMsg{
			From: "from2",
		},
		&_ChatMsg{
			From: "from3",
		},
	}
	err = offutil.UserOflnMsgUtil(kU1).Insert(msgs)
	assert.Nil(err)
	db := dbutil.UserOflnMsgUtil(kU1)
	msgsLoaded, errLoad := db.Load()
	assert.Nil(errLoad)

	for _, m := range msgsLoaded {
		fmt.Printf("msg: %v\n", m)
	}

	err = db.Remove()
	assert.Nil(err)
}
