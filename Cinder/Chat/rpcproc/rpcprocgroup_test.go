package rpcproc

import (
	"Cinder/Chat/mockcore"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createGroup(members []string) error {
	binMembers, err := json.Marshal(members)
	if err != nil {
		return errors.Wrap(err, "marshal error")
	}
	var errStr string = proc.RPC_CreateGroup(kGroup, binMembers)
	if errStr != "" {
		return fmt.Errorf("RPC_CreateGroup returns error: %s", errStr)
	}
	return nil
}

func Test_RPC_CreateGroup(t *testing.T) {
	proc.RPC_DeleteGroup(kGroup)
	err := createGroup(nil)
	if err != nil {
		t.Errorf("RPC_CreateGroup error: %v", err)
		return
	}
	err = createGroup([]string{kUser1})
	if err != nil {
		t.Errorf("RPC_CreateGroup error: %v", err)
		return
	}
	proc.RPC_DeleteGroup(kGroup)
}

func Test_RPC_DeleteGroup(t *testing.T) {
	proc.RPC_DeleteGroup(kGroup)
}

func Test_RPC_JoinGroup(t *testing.T) {
	proc.RPC_DeleteGroup(kGroup)
	proc.RPC_JoinGroup(kUser1, kGroup) // group not exist: group

	createGroup(nil)
	proc.RPC_JoinGroup(kUser1, kGroup)
	proc.RPC_DeleteGroup(kGroup)
}

func Test_RPC_LeaveGroup(t *testing.T) {
	proc.RPC_DeleteGroup(kGroup)
	proc.RPC_LeaveGroup(kUser1, kGroup) // group not exist

	createGroup([]string{kUser1})
	proc.RPC_LeaveGroup(kUser1, kGroup)
	proc.RPC_DeleteGroup(kGroup)
}

func Test_RPC_GetGroupMembers(t *testing.T) {
	var binMembers []byte = proc.RPC_GetGroupMembers(kGroup)
	members := []string{}
	err := json.Unmarshal(binMembers, &members)
	if err != nil {
		t.Errorf("RPC_GetGroupMembers result unmarshal error: %v", err)
		return
	}
}

func Test_RPC_SendGroupMessage(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()
	a := mock.Anything
	mockCore.On("RpcByID", a, a, a, a, a, a, a, a, a, a).Return(mockCore.Ch())

	proc.RPC_DeleteGroup(kGroup)
	proc.RPC_SendGroupMessage(kUser1, kGroup, []byte("123")) // group not exist

	createGroup([]string{kUser1, kUser2})
	proc.RPC_SendGroupMessage(kUser1, kGroup, []byte("111"))
	proc.RPC_LoginWithNickData(kUser1, "nick", []byte("111"), "peerSrvID")
	proc.RPC_SendGroupMessage(kUser1, kGroup, []byte("111"))
	proc.RPC_LoginWithNickData(kUser2, "nick", []byte("111"), "peerSrvID")
	proc.RPC_SendGroupMessage(kUser1, kGroup, []byte("111"))
	proc.RPC_Logout(kUser2)
	proc.RPC_Logout(kUser1)
	proc.RPC_DeleteGroup(kGroup)
}

func Test_RPC_GetGroupMemberCounts(t *testing.T) {
	createGroup([]string{"11", "2", "3"})
	grps := []string{kGroup}
	bin, err := json.Marshal(grps)
	assert.Nil(t, err)
	proc.RPC_GetGroupMemberCounts(bin)
	proc.RPC_DeleteGroup(kGroup)
}

func Test_RPC_GetGroupHistoryMessages(t *testing.T) {
	createGroup(nil)
	proc.RPC_GetGroupHistoryMessages(kGroup, 123)
	proc.RPC_DeleteGroup(kGroup)
}
