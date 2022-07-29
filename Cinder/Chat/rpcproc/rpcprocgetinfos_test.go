package rpcproc

import (
	"Cinder/Chat/chatapi"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RPC_GetFriendInfos(t *testing.T) {
	assert := require.New(t)

	loginUser1()
	proc.RPC_Logout(kUser1)
	loginUser2()
	proc.RPC_Logout(kUser2)

	ids := []string{kUser1, kUser2}
	binIDs, errJson := json.Marshal(ids)
	assert.NoError(errJson)

	buf := proc.RPC_GetFriendInfos(binIDs)
	infos := []chatapi.FriendInfo{}
	err := json.Unmarshal(buf, &infos)
	assert.NoError(err)
	assert.Equal(2, len(infos))
}
