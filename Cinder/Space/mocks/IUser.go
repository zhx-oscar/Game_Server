// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import CRpc "Cinder/Base/CRpc"
import Core "Cinder/Base/Core"
import Message "Cinder/Base/Message"
import Prop "Cinder/Base/Prop"
import Space "Cinder/Space"
import User "Cinder/Base/User"
import mock "github.com/stretchr/testify/mock"

// IUser is an autogenerated mock type for the IUser type
type IUser struct {
	mock.Mock
}

// GetID provides a mock function with given fields:
func (_m *IUser) GetID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMgr provides a mock function with given fields:
func (_m *IUser) GetMgr() User.IUserMgr {
	ret := _m.Called()

	var r0 User.IUserMgr
	if rf, ok := ret.Get(0).(func() User.IUserMgr); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(User.IUserMgr)
		}
	}

	return r0
}

// GetPeerServerID provides a mock function with given fields: srvType
func (_m *IUser) GetPeerServerID(srvType string) (string, error) {
	ret := _m.Called(srvType)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(srvType)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(srvType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProp provides a mock function with given fields:
func (_m *IUser) GetProp() Prop.IProp {
	ret := _m.Called()

	var r0 Prop.IProp
	if rf, ok := ret.Get(0).(func() Prop.IProp); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Prop.IProp)
		}
	}

	return r0
}

// GetPropType provides a mock function with given fields:
func (_m *IUser) GetPropType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetRealPtr provides a mock function with given fields:
func (_m *IUser) GetRealPtr() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// GetSpace provides a mock function with given fields:
func (_m *IUser) GetSpace() Space.ISpace {
	ret := _m.Called()

	var r0 Space.ISpace
	if rf, ok := ret.Get(0).(func() Space.ISpace); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Space.ISpace)
		}
	}

	return r0
}

// GetSrvInst provides a mock function with given fields:
func (_m *IUser) GetSrvInst() Core.ICore {
	ret := _m.Called()

	var r0 Core.ICore
	if rf, ok := ret.Get(0).(func() Core.ICore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Core.ICore)
		}
	}

	return r0
}

// GetType provides a mock function with given fields:
func (_m *IUser) GetType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetUserData provides a mock function with given fields:
func (_m *IUser) GetUserData() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Offline provides a mock function with given fields:
func (_m *IUser) Offline() {
	_m.Called()
}

// Rpc provides a mock function with given fields: srvType, methodName, args
func (_m *IUser) Rpc(srvType string, methodName string, args ...interface{}) chan *CRpc.RpcRet {
	var _ca []interface{}
	_ca = append(_ca, srvType, methodName)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 chan *CRpc.RpcRet
	if rf, ok := ret.Get(0).(func(string, string, ...interface{}) chan *CRpc.RpcRet); ok {
		r0 = rf(srvType, methodName, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *CRpc.RpcRet)
		}
	}

	return r0
}

// SendToAllClient provides a mock function with given fields: msg
func (_m *IUser) SendToAllClient(msg Message.IMessage) {
	_m.Called(msg)
}

// SendToAllClientExceptMe provides a mock function with given fields: msg
func (_m *IUser) SendToAllClientExceptMe(msg Message.IMessage) {
	_m.Called(msg)
}

// SendToClient provides a mock function with given fields: msg
func (_m *IUser) SendToClient(msg Message.IMessage) error {
	ret := _m.Called(msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(Message.IMessage) error); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendToPeerServer provides a mock function with given fields: srvType, msg
func (_m *IUser) SendToPeerServer(srvType string, msg Message.IMessage) error {
	ret := _m.Called(srvType, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, Message.IMessage) error); ok {
		r0 = rf(srvType, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendToPeerUser provides a mock function with given fields: srvType, msg
func (_m *IUser) SendToPeerUser(srvType string, msg Message.IMessage) error {
	ret := _m.Called(srvType, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, Message.IMessage) error); ok {
		r0 = rf(srvType, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
