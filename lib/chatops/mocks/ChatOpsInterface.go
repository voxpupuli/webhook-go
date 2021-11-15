// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	chat "github.com/pandatix/gocket-chat/api/chat"
	chatops "github.com/voxpupuli/webhook-go/lib/chatops"

	mock "github.com/stretchr/testify/mock"
)

// ChatOpsInterface is an autogenerated mock type for the ChatOpsInterface type
type ChatOpsInterface struct {
	mock.Mock
}

// PostMessage provides a mock function with given fields: code, target
func (_m *ChatOpsInterface) PostMessage(code int, target string) (*chatops.ChatOpsResponse, error) {
	ret := _m.Called(code, target)

	var r0 *chatops.ChatOpsResponse
	if rf, ok := ret.Get(0).(func(int, string) *chatops.ChatOpsResponse); ok {
		r0 = rf(code, target)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chatops.ChatOpsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string) error); ok {
		r1 = rf(code, target)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// rocketChat provides a mock function with given fields: code, target
func (_m *ChatOpsInterface) rocketChat(code int, target string) (*chat.PostMessageResponse, error) {
	ret := _m.Called(code, target)

	var r0 *chat.PostMessageResponse
	if rf, ok := ret.Get(0).(func(int, string) *chat.PostMessageResponse); ok {
		r0 = rf(code, target)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chat.PostMessageResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string) error); ok {
		r1 = rf(code, target)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// slack provides a mock function with given fields: code, target
func (_m *ChatOpsInterface) slack(code int, target string) (*string, *string, error) {
	ret := _m.Called(code, target)

	var r0 *string
	if rf, ok := ret.Get(0).(func(int, string) *string); ok {
		r0 = rf(code, target)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 *string
	if rf, ok := ret.Get(1).(func(int, string) *string); ok {
		r1 = rf(code, target)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*string)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int, string) error); ok {
		r2 = rf(code, target)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
