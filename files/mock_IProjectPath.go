// Code generated by mockery v2.14.1. DO NOT EDIT.

package files

import (
	fs "io/fs"

	mock "github.com/stretchr/testify/mock"
)

// MockIProjectPath is an autogenerated mock type for the IProjectPath type
type MockIProjectPath struct {
	mock.Mock
}

// FileDoesNotExist provides a mock function with given fields:
func (_m *MockIProjectPath) FileDoesNotExist() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetFileStringWithoutExt provides a mock function with given fields:
func (_m *MockIProjectPath) GetFileStringWithoutExt() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetPath provides a mock function with given fields:
func (_m *MockIProjectPath) GetPath() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// HasParentPath provides a mock function with given fields: _a0
func (_m *MockIProjectPath) HasParentPath(_a0 IProjectPath) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(IProjectPath) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsPhpFile provides a mock function with given fields:
func (_m *MockIProjectPath) IsPhpFile() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsTestFile provides a mock function with given fields:
func (_m *MockIProjectPath) IsTestFile() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Join provides a mock function with given fields: _a0
func (_m *MockIProjectPath) Join(_a0 ...string) IProjectPath {
	_va := make([]interface{}, len(_a0))
	for _i := range _a0 {
		_va[_i] = _a0[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 IProjectPath
	if rf, ok := ret.Get(0).(func(...string) IProjectPath); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(IProjectPath)
		}
	}

	return r0
}

// MakeAbsolute provides a mock function with given fields:
func (_m *MockIProjectPath) MakeAbsolute() IProjectPath {
	ret := _m.Called()

	var r0 IProjectPath
	if rf, ok := ret.Get(0).(func() IProjectPath); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(IProjectPath)
		}
	}

	return r0
}

// MakeRelative provides a mock function with given fields:
func (_m *MockIProjectPath) MakeRelative() IProjectPath {
	ret := _m.Called()

	var r0 IProjectPath
	if rf, ok := ret.Get(0).(func() IProjectPath); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(IProjectPath)
		}
	}

	return r0
}

// ReadFile provides a mock function with given fields: _a0
func (_m *MockIProjectPath) ReadFile(_a0 fs.FS) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(fs.FS) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(fs.FS) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockIProjectPath interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockIProjectPath creates a new instance of MockIProjectPath. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockIProjectPath(t mockConstructorTestingTNewMockIProjectPath) *MockIProjectPath {
	mock := &MockIProjectPath{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
