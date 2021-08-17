// Code generated by mockery (devel). DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// AuthService is an autogenerated mock type for the AuthService type
type AuthService struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: _a0, _a1
func (_m *AuthService) Authenticate(_a0 string, _a1 string) (string, string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string, string) string); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}