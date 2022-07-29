// mockcore 仅用于测试时替换 Core.Inst
package mockcore

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Core"
	"Cinder/Base/Core/mocks"
	"errors"
)

// MockCore 用来替换 Core.Inst
/* 使用示例
func TestRpc(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()

	mockCore.On("RpcByID", "srvID", "methodName", "arg1").Return(mockCore.Ch())
	ret := Rpc("srvID", "methodName", "arg1")
	require.Error(t, ret.Err)
}
*/
type MockCore struct {
	mocks.ICore

	// 用于恢复 Core.Inst
	originalCore Core.ICore
}

func (m *MockCore) SetupMock() {
	m.originalCore = Core.Inst
	Core.Inst = m
}

func (m *MockCore) TearDownMock() {
	Core.Inst = m.originalCore
}

func (m *MockCore) Ch() chan *CRpc.RpcRet {
	ch := make(chan *CRpc.RpcRet, 1)
	ch <- &CRpc.RpcRet{
		Err: errors.New("mocked rpc return"),
	}
	close(ch)
	return ch
}
