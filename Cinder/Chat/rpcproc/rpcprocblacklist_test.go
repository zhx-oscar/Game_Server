package rpcproc

import (
	"Cinder/Chat/chatapi"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RPC_AddFriendToBlacklist(t *testing.T) {
	loginUser1()
	defer proc.RPC_Logout(kUser1)

	proc.RPC_AddFriendToBlacklist(kUser1, "blacklistedID")
}

func Test_RPC_RemoveFriendFromBlacklist(t *testing.T) {
	loginUser1()
	defer proc.RPC_Logout(kUser1)

	proc.RPC_RemoveFriendFromBlacklist(kUser1, "blacklistedID")
	proc.RPC_AddFriendToBlacklist(kUser1, "blacklistedID")
	proc.RPC_RemoveFriendFromBlacklist(kUser1, "blacklistedID")
}

func Test_RPC_GetFriendBlacklist(t *testing.T) {
	loginUser1()
	defer proc.RPC_Logout(kUser1)

	var bin []byte = proc.RPC_GetFriendBlacklist(kUser1)
	if bin == nil {
		return
	}

	infos := []chatapi.FriendInfo{}
	err := json.Unmarshal(bin, &infos)
	assert.Nil(t, err)
}
