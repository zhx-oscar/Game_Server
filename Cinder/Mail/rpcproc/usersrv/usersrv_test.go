package usersrv

import (
	"Cinder/Base/Const"
	"Cinder/Mail/rpcproc/mockcore"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoadFromDB(t *testing.T) {
	err := LoadFromDB()
	require.NoError(t, err)
}

func TestGetUserSrvIDs(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()

	mockCore.On("GetSrvTypeByID", mock.AnythingOfType("string")).Return("", errors.New("no such srv ID"))
	mockCore.On("GetSrvIDSByType", Const.Mail).Return([]string{"1", "2"}, nil)
	mockCore.On("RpcByID", mock.AnythingOfType("string"), "RPC_SyncUserSrvID", mock.AnythingOfType("string")).Return(mockCore.Ch())

	const kSrvID = "no_such_srv_ID"
	InsertAndBroadcast(kSrvID)
	_ = GetUserSrvIDs()
}
