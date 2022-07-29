package internal

import (
	"Cinder/Matcher/rpcmsg"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsRequest(t *testing.T) {
	assert := require.New(t)
	assert.False(rpcmsg.IsRequest(struct{}{}))
	assert.True(rpcmsg.IsRequest(&rpcmsg.BroadcastRoomReq{}))

	assert.True(rpcmsg.IsRequest(rpcmsg.BroadcastRoomReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.CreateTeamReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.GetRoomInfoReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.GetRoomListReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.RoleCreateRoomReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.RoleJoinRandomRoomReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.RoleJoinRoomReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.JoinTeamReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.RoleCreateRoomReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.RoleJoinRandomRoomReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.RoleJoinRoomReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.SetRoomDataReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.UpdateRoomRoleDataReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.TeamCreateRoomReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.TeamJoinRandomRoomReq{}))
	assert.True(rpcmsg.IsRequest(rpcmsg.TeamJoinRoomReq{}))
}
