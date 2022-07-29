// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import types "Cinder/Mail/mailapi/types"

// IMailService is an autogenerated mock type for the IMailService type
type IMailService struct {
	mock.Mock
}

// BatchDelete provides a mock function with given fields: userID, mailIDs
func (_m *IMailService) BatchDelete(userID string, mailIDs []string) error {
	ret := _m.Called(userID, mailIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(userID, mailIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Broadcast provides a mock function with given fields: m
func (_m *IMailService) Broadcast(m *types.Mail) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(*types.Mail) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: userID, mailID
func (_m *IMailService) Delete(userID string, mailID string) error {
	ret := _m.Called(userID, mailID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userID, mailID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListMail provides a mock function with given fields: userID
func (_m *IMailService) ListMail(userID string) ([]*types.Mail, error) {
	ret := _m.Called(userID)

	var r0 []*types.Mail
	if rf, ok := ret.Get(0).(func(string) []*types.Mail); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Mail)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: userID
func (_m *IMailService) Login(userID string) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Logout provides a mock function with given fields: userID
func (_m *IMailService) Logout(userID string) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MarkAsRead provides a mock function with given fields: userID, mailID
func (_m *IMailService) MarkAsRead(userID string, mailID string) error {
	ret := _m.Called(userID, mailID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userID, mailID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MarkAsUnread provides a mock function with given fields: userID, mailID
func (_m *IMailService) MarkAsUnread(userID string, mailID string) error {
	ret := _m.Called(userID, mailID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userID, mailID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MarkAttachmentsAsReceived provides a mock function with given fields: userID, mailID
func (_m *IMailService) MarkAttachmentsAsReceived(userID string, mailID string) error {
	ret := _m.Called(userID, mailID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userID, mailID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MarkAttachmentsAsUnreceived provides a mock function with given fields: userID, mailID
func (_m *IMailService) MarkAttachmentsAsUnreceived(userID string, mailID string) error {
	ret := _m.Called(userID, mailID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userID, mailID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Send provides a mock function with given fields: m
func (_m *IMailService) Send(m *types.Mail) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(*types.Mail) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
