// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import CRpc "Cinder/Base/CRpc"
import Core "Cinder/Base/Core"
import Message "Cinder/Base/Message"
import Prop "Cinder/Base/Prop"
import SrvNet "Cinder/Base/SrvNet"
import mock "github.com/stretchr/testify/mock"

// ICore is an autogenerated mock type for the ICore type
type ICore struct {
	mock.Mock
}

// Broadcast provides a mock function with given fields: srvType, msg
func (_m *ICore) Broadcast(srvType string, msg Message.IMessage) error {
	ret := _m.Called(srvType, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, Message.IMessage) error); ok {
		r0 = rf(srvType, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CallRpcToAllUsers provides a mock function with given fields: srvType, methodName, args
func (_m *ICore) CallRpcToAllUsers(srvType string, methodName string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, srvType, methodName)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// CallRpcToUser provides a mock function with given fields: userID, srvType, methodName, args
func (_m *ICore) CallRpcToUser(userID string, srvType string, methodName string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, userID, srvType, methodName)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// CallRpcToUsers provides a mock function with given fields: userIDS, srvType, methodName, args
func (_m *ICore) CallRpcToUsers(userIDS []string, srvType string, methodName string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, userIDS, srvType, methodName)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// CreateProp provides a mock function with given fields: propType
func (_m *ICore) CreateProp(propType string) (Prop.IProp, error) {
	ret := _m.Called(propType)

	var r0 Prop.IProp
	if rf, ok := ret.Get(0).(func(string) Prop.IProp); ok {
		r0 = rf(propType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Prop.IProp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(propType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePropObject provides a mock function with given fields: propObjectType, id, propData, userData
func (_m *ICore) CreatePropObject(propObjectType string, id string, propData []byte, userData interface{}) (Prop.IPropObject, error) {
	ret := _m.Called(propObjectType, id, propData, userData)

	var r0 Prop.IPropObject
	if rf, ok := ret.Get(0).(func(string, string, []byte, interface{}) Prop.IPropObject); ok {
		r0 = rf(propObjectType, id, propData, userData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Prop.IPropObject)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, []byte, interface{}) error); ok {
		r1 = rf(propObjectType, id, propData, userData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Debug provides a mock function with given fields: v
func (_m *ICore) Debug(v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Destroy provides a mock function with given fields:
func (_m *ICore) Destroy() {
	_m.Called()
}

// DestroyPropObject provides a mock function with given fields: id
func (_m *ICore) DestroyPropObject(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Error provides a mock function with given fields: v
func (_m *ICore) Error(v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// GetNetNode provides a mock function with given fields:
func (_m *ICore) GetNetNode() SrvNet.INode {
	ret := _m.Called()

	var r0 SrvNet.INode
	if rf, ok := ret.Get(0).(func() SrvNet.INode); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(SrvNet.INode)
		}
	}

	return r0
}

// GetServiceID provides a mock function with given fields:
func (_m *ICore) GetServiceID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetServiceType provides a mock function with given fields:
func (_m *ICore) GetServiceType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetSrvIDByType provides a mock function with given fields: srvType
func (_m *ICore) GetSrvIDByType(srvType string) (string, error) {
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

// GetSrvIDSByType provides a mock function with given fields: srvType
func (_m *ICore) GetSrvIDSByType(srvType string) ([]string, error) {
	ret := _m.Called(srvType)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(srvType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(srvType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSrvTypeByID provides a mock function with given fields: srvID
func (_m *ICore) GetSrvTypeByID(srvID string) (string, error) {
	ret := _m.Called(srvID)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(srvID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(srvID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Info provides a mock function with given fields: v
func (_m *ICore) Info(v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Init provides a mock function with given fields: info
func (_m *ICore) Init(info *Core.Info) error {
	ret := _m.Called(info)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Core.Info) error); ok {
		r0 = rf(info)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPropOwner provides a mock function with given fields: realPtr
func (_m *ICore) NewPropOwner(realPtr interface{}) Prop.IPropOwner {
	ret := _m.Called(realPtr)

	var r0 Prop.IPropOwner
	if rf, ok := ret.Get(0).(func(interface{}) Prop.IPropOwner); ok {
		r0 = rf(realPtr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Prop.IPropOwner)
		}
	}

	return r0
}

// RegisterProp provides a mock function with given fields: propType, propProto
func (_m *ICore) RegisterProp(propType string, propProto Prop.IProp) {
	_m.Called(propType, propProto)
}

// RegisterPropObject provides a mock function with given fields: propObjectType, object
func (_m *ICore) RegisterPropObject(propObjectType string, object Prop.IPropObject) {
	_m.Called(propObjectType, object)
}

// RpcByID provides a mock function with given fields: srvID, methodName, args
func (_m *ICore) RpcByID(srvID string, methodName string, args ...interface{}) chan *CRpc.RpcRet {
	var _ca []interface{}
	_ca = append(_ca, srvID, methodName)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 chan *CRpc.RpcRet
	if rf, ok := ret.Get(0).(func(string, string, ...interface{}) chan *CRpc.RpcRet); ok {
		r0 = rf(srvID, methodName, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *CRpc.RpcRet)
		}
	}

	return r0
}

// RpcByType provides a mock function with given fields: srvType, methodName, args
func (_m *ICore) RpcByType(srvType string, methodName string, args ...interface{}) chan *CRpc.RpcRet {
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

// Send provides a mock function with given fields: srvID, msg
func (_m *ICore) Send(srvID string, msg Message.IMessage) error {
	ret := _m.Called(srvID, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, Message.IMessage) error); ok {
		r0 = rf(srvID, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendMessageToAllUsers provides a mock function with given fields: srvType, message
func (_m *ICore) SendMessageToAllUsers(srvType string, message Message.IMessage) {
	_m.Called(srvType, message)
}

// SendMessageToUser provides a mock function with given fields: userID, srvType, message
func (_m *ICore) SendMessageToUser(userID string, srvType string, message Message.IMessage) {
	_m.Called(userID, srvType, message)
}

// SendMessageToUsers provides a mock function with given fields: userIDS, srvType, message
func (_m *ICore) SendMessageToUsers(userIDS []string, srvType string, message Message.IMessage) {
	_m.Called(userIDS, srvType, message)
}

// Warning provides a mock function with given fields: v
func (_m *ICore) Warning(v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}
