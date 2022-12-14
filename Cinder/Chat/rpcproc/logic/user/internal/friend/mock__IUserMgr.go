// Code generated by mockery v1.0.0. DO NOT EDIT.

package friend

import mock "github.com/stretchr/testify/mock"
import types "Cinder/Chat/rpcproc/logic/types"

// mock_IUserMgr is an autogenerated mock type for the _IUserMgr type
type mock_IUserMgr struct {
	mock.Mock
}

// GetUserFriendMgr provides a mock function with given fields: userID
func (_m *mock_IUserMgr) GetUserFriendMgr(userID types.UserID) types.IFriendMgr {
	ret := _m.Called(userID)

	var r0 types.IFriendMgr
	if rf, ok := ret.Get(0).(func(types.UserID) types.IFriendMgr); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.IFriendMgr)
		}
	}

	return r0
}
