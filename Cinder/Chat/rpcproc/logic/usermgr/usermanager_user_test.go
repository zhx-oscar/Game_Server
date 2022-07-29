package usermgr

import (
	"Cinder/Chat/mockcore"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendChatMessage(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()
	a := mock.Anything
	mockCore.On("RpcByID", a, a, a, a, a, a, a, a).Return(mockCore.Ch())

	const kUser1 = "test_user1"
	const kUser2 = "test_user2"
	assert := assert.New(t)

	mgr := GetUserMgr()
	u, err := mgr.LoginUser(kUser1, "srvID")
	assert.NotNil(u)
	assert.Nil(err)
	u, err = mgr.LoginUser(kUser2, "srvID")
	assert.NotNil(u)
	assert.Nil(err)
	user1 := mgr.GetUser(kUser1)
	assert.NotNil(user1)
	user1.SendChatMessage(kUser2, []byte("1->2"))
	user1.SendChatMessage("user_offline", []byte("1->offline"))
	mgr.RemoveUser(kUser1)
	mgr.RemoveUser(kUser2)
}
