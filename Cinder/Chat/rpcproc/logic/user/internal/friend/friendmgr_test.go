package friend

import (
	"Cinder/Chat/rpcproc/logic/user/internal/userinfo"
	"testing"
)

func TestDeleteFriendActive(t *testing.T) {
	mockUserMgr := new(mock_IUserMgr)
	originalUserMgr := userMgr

	userMgr = mockUserMgr
	defer func() {
		userMgr = originalUserMgr
	}()

	mockUserMgr.On("GetUserFriendMgr", UserID("user2")).Return(nil)

	info := userinfo.NewUserInfo("user1", "svrID")
	f1 := NewFriendMgr("user1", info)
	f1.AddFriendWithoutDB("user2")
	f1.DeleteFriendActive("user2")
}
