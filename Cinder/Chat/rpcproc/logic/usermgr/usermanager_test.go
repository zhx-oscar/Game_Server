package usermgr

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/bc"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type _MockedGroupMgr struct {
	mock.Mock
}

func (m *_MockedGroupMgr) LoginGroups(userID UserID, groups []GroupID) {
}

func (m *_MockedGroupMgr) LogoutGroups(userID UserID, groups []GroupID) {
}

func (m *_MockedGroupMgr) GetGroupOfflineMessageMap(g2s map[GroupID]SequenceID) map[string][]*chatapi.ChatMessage {
	return nil
}

func (m *_MockedGroupMgr) GetOrLoadGroup(groupID types.GroupID) (types.IGroup, error) {
	return nil, fmt.Errorf("test mock")
}

func init() {
	mockedGroupMgr := &_MockedGroupMgr{}
	bc.SetUserMgr(GetUserMgr())
	bc.SetGroupMgr(mockedGroupMgr)
	user.SetUserMgr(GetUserMgr())
	user.SetGroupMgr(mockedGroupMgr)
}

func TestLoginUser(t *testing.T) {
	u, err := GetUserMgr().LoginUser("test_add_user", "srvID")
	assert.NotNil(t, u)
	assert.Nil(t, err)
}
