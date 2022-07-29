package rpcproc

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/mockcore"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_RPC_AddFriendReq(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()
	a := mock.Anything
	mockCore.On("RpcByID", a, a, a, a, a).Return(mockCore.Ch())

	loginUser1()
	defer proc.RPC_Logout(kUser1)
	proc.RPC_AddFriendReq(kUser1, "user2", []byte("111"))

	loginUser2()
	defer proc.RPC_Logout(kUser2)
	proc.RPC_AddFriendReq(kUser1, "user2", []byte("111"))
}

func Test_RPC_ReplyAddFriendReq(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()
	a := mock.Anything
	mockCore.On("RpcByID", a, a, a, a, a).Return(mockCore.Ch())

	loginUser1()
	defer proc.RPC_Logout(kUser1)
	proc.RPC_ReplyAddFriendReq(kUser1, kUser2, true)
	proc.RPC_ReplyAddFriendReq(kUser1, "no_user_fromID", false)

	loginUser2()
	defer proc.RPC_Logout(kUser2)
	proc.RPC_ReplyAddFriendReq(kUser1, kUser2, true)
}

func Test_RPC_DeleteFriend(t *testing.T) {
	loginUser1()
	defer proc.RPC_Logout(kUser1)
	proc.RPC_DeleteFriend(kUser1, "friendID")
}

func Test_RPC_GetFriendList(t *testing.T) {
	loginUser1()
	defer proc.RPC_Logout(kUser1)

	var bin []byte = proc.RPC_GetFriendList(kUser1)
	if bin == nil {
		return
	}

	infos := []chatapi.FriendInfo{}
	err := json.Unmarshal(bin, &infos)
	assert.Nil(t, err)
}
