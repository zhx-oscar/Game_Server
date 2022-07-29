package dbutil

import (
	"Cinder/Chat/chatapi"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Msg = chatapi.ChatMessage

func TestGroupsMessagessUtil(t *testing.T) {
	var err error
	assert := assert.New(t)
	db := GroupsMessagessUtil("test_group")

	msgs := make(map[SequenceID]*Msg)
	for i := 0; i < 10000; i++ {
		msgs[SequenceID(i)] = &Msg{}
	}
	err = db.Insert(msgs, 1, 10)
	assert.Nil(err)

	msgsLoaded, err := db.Load()
	assert.Nil(err)
	for seq, msg := range msgsLoaded {
		fmt.Printf("msg[%v]: %#v\n", seq, msg)
		break
	}
}

func TestInsertSeq(t *testing.T) {
	var err error
	assert := assert.New(t)
	db := GroupsMessagessUtil("test_group")
	const kMaxSeq = ^SequenceID(0)

	err = db.Insert(nil, 0, 0)
	assert.Nil(err)
	err = db.Insert(nil, 1, 1)
	assert.Nil(err)
	// err = db.Insert(nil, kMaxSeq, kMaxSeq)
	// assert.Nil(err)
	err = db.Insert(nil, 0, 1)
	assert.Nil(err)
	err = db.Insert(nil, 10, 1)
	assert.Nil(err)
	err = db.Insert(nil, 0, 9999999)
	assert.Nil(err)
	// err = db.Insert(nil, 0, kMaxSeq)
	// assert.Nil(err)
}

func TestNilMap(t *testing.T) {
	var m map[int]bool
	assert := assert.New(t)
	assert.Nil(m)
	b, ok := m[1111]
	assert.False(ok)
	assert.False(b)
	assert.False(m[123])
}
