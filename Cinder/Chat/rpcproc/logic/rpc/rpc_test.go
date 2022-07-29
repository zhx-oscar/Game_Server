package rpc

import (
	"Cinder/Chat/mockcore"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRpc(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()

	mockCore.On("RpcByID", "srvID", "methodName", "arg1").Return(mockCore.Ch())
	ret := Rpc("srvID", "methodName", "arg1")
	require.Error(t, ret.Err)
}
