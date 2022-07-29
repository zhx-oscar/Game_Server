package mtypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalNotifyMsg(t *testing.T) {
	assert := require.New(t)
	msg := NotifyMsg{
		RolesJoinRoomMsg: &RolesJoinRoomMsg{
			RoomInfo: RoomInfo{
				RoomID: "testRoomID",
			},
			JoinedRoleIDs: []mtypes.RoleID{"testRoleID"},
		},
	}
	buf, errJson := json.Marshal(msg)
	assert.NoError(errJson)

	msg2, err2 := UnmarshalNotifyMsg(buf)
	assert.NoError(err2)
	assert.Nil(msg2.JoinTeamMsg)
	assert.NotNil(msg2.RolesJoinRoomMsg)
	assert.Equal(msg, msg2)
	assert.Equal(*msg.RolesJoinRoomMsg, msg2.GetRolesJoinRoomMsg())
	assert.Equal(JoinTeamMsg{}, msg2.GetJoinTeamMsg())
	assert.Nil(msg2.JoinTeamMsg)
}
