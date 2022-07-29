package rpcproc

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/mockcore"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_RPC_FollowUnFollowFriendReq(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()
	a := mock.Anything
	mockCore.On("RpcByID", a, a, a, a, a, a, a, a).Return(mockCore.Ch())

	loginUser1()
	defer proc.RPC_Logout(kUser1)
	proc.RPC_FollowFriendReq(kUser1, "followedID")
	proc.RPC_UnFollowFriendReq(kUser1, "unfollowedID")

	loginUser2()
	defer proc.RPC_Logout(kUser2)
	proc.RPC_FollowFriendReq(kUser1, kUser2)
	proc.RPC_UnFollowFriendReq(kUser1, kUser2)
}

func Test_RPC_GetFollowingList(t *testing.T) {
	loginUser1()
	defer proc.RPC_Logout(kUser1)

	var bin []byte = proc.RPC_GetFollowingList(kUser1)
	if bin == nil {
		return
	}
	infos := []chatapi.FriendInfo{}
	err := json.Unmarshal(bin, &infos)
	assert.Nil(t, err)
}

func Test_RPC_GetFollowerList(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()
	a := mock.Anything
	mockCore.On("RpcByID", a, a, a, a, a, a, a, a).Return(mockCore.Ch())

	loginUser1()
	defer proc.RPC_Logout(kUser1)
	loginUser2()
	defer proc.RPC_Logout(kUser2)
	proc.RPC_FollowFriendReq(kUser2, kUser1)
	defer proc.RPC_UnFollowFriendReq(kUser2, kUser1)

	var bin []byte = proc.RPC_GetFollowerList(kUser1)
	if bin == nil {
		return
	}

	infos := []chatapi.FollowerInfo{}
	err := json.Unmarshal(bin, &infos)
	assert.Nil(t, err)
}
