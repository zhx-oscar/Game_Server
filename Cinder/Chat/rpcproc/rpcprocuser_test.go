package rpcproc

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/mockcore"
	_ "Cinder/Chat/rpcproc/logic" // init()
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func loginUser1() {
	proc.RPC_LoginWithNickData(kUser1, "nick", []byte("111"), "peerSrvID")
}
func loginUser2() {
	proc.RPC_LoginWithNickData(kUser2, "nick", []byte("222"), "peerSrvID")
}

func Test_RPC_Login(t *testing.T) {
	proc.RPC_Login(kUser1, "peerSrvID")
	proc.RPC_Logout(kUser1)
}

func Test_RPC_LoginWithNickData(t *testing.T) {
	proc.RPC_LoginWithNickData(kUser1, "nick", []byte("111"), "peerSrvID")
	proc.RPC_LoginWithNickData(kUser1, "nick2", []byte("111"), "peerSrvID")
	proc.RPC_LoginWithNickData(kUser1, "nick2", []byte("1112"), "peerSrvID")
	proc.RPC_Logout(kUser1)
}

func Test_RPC_Logout(t *testing.T) {
	proc.RPC_Logout(kUser1)
	loginUser1()
	proc.RPC_Logout(kUser1)
}

func Test_RPC_Nick(t *testing.T) {
	assert := require.New(t)
	loginUser1()
	errStr := proc.RPC_SetNick(kUser1, "aaaabbbb")
	assert.Empty(errStr)
	nickErrStr, nick := proc.RPC_GetNick(kUser1)
	assert.Empty(nickErrStr)
	assert.Equal("aaaabbbb", nick)
	proc.RPC_Logout(kUser1)
}

func Test_RPC_Data(t *testing.T) {
	assert := require.New(t)
	loginUser1()
	errStr := proc.RPC_SetActiveData(kUser1, []byte("activeData"))
	assert.Empty(errStr)
	actErrStr, activeData := proc.RPC_GetActiveData(kUser1)
	assert.Empty(actErrStr)
	assert.Equal([]byte("activeData"), activeData)
	errStr = proc.RPC_SetPassiveData(kUser1, []byte("passiveData"))
	assert.Empty(errStr)
	pasErrStr, passiveData := proc.RPC_GetPassiveData(kUser1)
	assert.Empty(pasErrStr)
	assert.Equal([]byte("passiveData"), passiveData)
	proc.RPC_Logout(kUser1)
}

func Test_RPC_SendMessage(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()
	a := mock.Anything
	mockCore.On("RpcByID", a, a, a, a, a, a, a, a, a).Return(mockCore.Ch())

	proc.RPC_Logout(kUser1)
	proc.RPC_SendMessage(kUser1, kUser2, []byte("1to2")) // failed to get user: can't find user

	loginUser1()
	loginUser2()
	proc.RPC_SendMessage(kUser1, kUser2, []byte("1to2"))
	proc.RPC_Logout(kUser1)
	proc.RPC_Logout(kUser2)
}

func Test_RPC_GetOfflineMessage(t *testing.T) {
	proc.RPC_Logout(kUser1)

	// kUser2 发消息给离线的 kUser1
	loginUser2()
	for i := 0; i < 105; i++ {
		proc.RPC_SendMessage(kUser2, kUser1, []byte("222111"))
	}
	proc.RPC_Logout(kUser2)

	loginUser1()
	var bin []byte = proc.RPC_GetOfflineMessage(kUser1)
	assert.NotNil(t, bin)

	msg := &chatapi.OfflineMessage{}
	err := json.Unmarshal(bin, msg)
	if err != nil {
		t.Errorf("RPC_GetOfflineMessage result unmarshal error: %v", err)
	}
	assert.True(t, len(msg.P2PMessages) >= 105)
	proc.RPC_Logout(kUser1)
}
