// Code generated by MockGen. DO NOT EDIT.
// Source: repo.go

// Package posts is a generated GoMock package.
package posts

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockFindInterface is a mock of FindInterface interface
type MockFindInterface struct {
	ctrl     *gomock.Controller
	recorder *MockFindInterfaceMockRecorder
}

// MockFindInterfaceMockRecorder is the mock recorder for MockFindInterface
type MockFindInterfaceMockRecorder struct {
	mock *MockFindInterface
}

// NewMockFindInterface creates a new mock instance
func NewMockFindInterface(ctrl *gomock.Controller) *MockFindInterface {
	mock := &MockFindInterface{ctrl: ctrl}
	mock.recorder = &MockFindInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFindInterface) EXPECT() *MockFindInterfaceMockRecorder {
	return m.recorder
}

// One mocks base method
func (m *MockFindInterface) One(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "One", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// One indicates an expected call of One
func (mr *MockFindInterfaceMockRecorder) One(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "One", reflect.TypeOf((*MockFindInterface)(nil).One), arg0)
}

// All mocks base method
func (m *MockFindInterface) All(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// All indicates an expected call of All
func (mr *MockFindInterfaceMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockFindInterface)(nil).All), arg0)
}

// MockPostRepositoryDBInterface is a mock of PostRepositoryDBInterface interface
type MockPostRepositoryDBInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepositoryDBInterfaceMockRecorder
}

// MockPostRepositoryDBInterfaceMockRecorder is the mock recorder for MockPostRepositoryDBInterface
type MockPostRepositoryDBInterfaceMockRecorder struct {
	mock *MockPostRepositoryDBInterface
}

// NewMockPostRepositoryDBInterface creates a new mock instance
func NewMockPostRepositoryDBInterface(ctrl *gomock.Controller) *MockPostRepositoryDBInterface {
	mock := &MockPostRepositoryDBInterface{ctrl: ctrl}
	mock.recorder = &MockPostRepositoryDBInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPostRepositoryDBInterface) EXPECT() *MockPostRepositoryDBInterfaceMockRecorder {
	return m.recorder
}

// Find mocks base method
func (m *MockPostRepositoryDBInterface) Find(arg0 interface{}) FindInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0)
	ret0, _ := ret[0].(FindInterface)
	return ret0
}

// Find indicates an expected call of Find
func (mr *MockPostRepositoryDBInterfaceMockRecorder) Find(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockPostRepositoryDBInterface)(nil).Find), arg0)
}

// Insert mocks base method
func (m *MockPostRepositoryDBInterface) Insert(arg0 ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Insert", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (mr *MockPostRepositoryDBInterfaceMockRecorder) Insert(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockPostRepositoryDBInterface)(nil).Insert), arg0...)
}

// Update mocks base method
func (m *MockPostRepositoryDBInterface) Update(arg0, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockPostRepositoryDBInterfaceMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPostRepositoryDBInterface)(nil).Update), arg0, arg1)
}

// Remove mocks base method
func (m *MockPostRepositoryDBInterface) Remove(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockPostRepositoryDBInterfaceMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockPostRepositoryDBInterface)(nil).Remove), arg0)
}
