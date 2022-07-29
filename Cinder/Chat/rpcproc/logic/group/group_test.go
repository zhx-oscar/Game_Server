package group

import (
	"Cinder/Chat/mockcore"
	"Cinder/Chat/rpcproc/logic/bc"
	"Cinder/Chat/rpcproc/logic/user"
	"Cinder/Chat/rpcproc/logic/usermgr"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	bc.SetUserMgr(usermgr.GetUserMgr())
	bc.SetGroupMgr(GetGroupMgr())
	user.SetUserMgr(usermgr.GetUserMgr())
}

func TestSendGroupMessage(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()
	a := mock.Anything
	mockCore.On("RpcByID", a, a, a, a, a, a, a, a, a).Return(mockCore.Ch())

	const kGroup = "test_group"
	const kUser1 = "test_user1"
	const kUser2 = "test_user2"

	assert := assert.New(t)
	umgr := usermgr.GetUserMgr()

	GetGroupMgr().InsertMembersToGroup(kGroup, []UserID{kUser1, kUser2, "some_user_offline"})
	u, err := umgr.LoginUser(kUser1, "srvID")
	assert.NotNil(u)
	assert.Nil(err)
	u, err = umgr.LoginUser(kUser2, "srvID")
	assert.NotNil(u)
	assert.Nil(err)
	group, errG := GetGroupMgr().GetOrLoadGroup(kGroup)
	assert.Nil(errG)
	assert.NotNil(group)
	group.SendGroupMessage(kUser1, []byte("1->g"))
	err = GetGroupMgr().DeleteGroup(kGroup)
	assert.Nil(err)
	umgr.RemoveUser(kUser1)
	umgr.RemoveUser(kUser2)
}
