// Code generated by MockGen. DO NOT EDIT.
// Source: internal/gordon/env_gen.go

// Package gordon_test is a generated GoMock package.
package gordon_test

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockEnv is a mock of Env interface
type MockEnv struct {
	ctrl     *gomock.Controller
	recorder *MockEnvMockRecorder
}

// MockEnvMockRecorder is the mock recorder for MockEnv
type MockEnvMockRecorder struct {
	mock *MockEnv
}

// NewMockEnv creates a new mock instance
func NewMockEnv(ctrl *gomock.Controller) *MockEnv {
	mock := &MockEnv{ctrl: ctrl}
	mock.recorder = &MockEnvMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEnv) EXPECT() *MockEnvMockRecorder {
	return m.recorder
}

// Architecture mocks base method
func (m *MockEnv) Architecture() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Architecture")
	ret0, _ := ret[0].(string)
	return ret0
}

// Architecture indicates an expected call of Architecture
func (mr *MockEnvMockRecorder) Architecture() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Architecture", reflect.TypeOf((*MockEnv)(nil).Architecture))
}

// Bin mocks base method
func (m *MockEnv) Bin() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bin")
	ret0, _ := ret[0].(string)
	return ret0
}

// Bin indicates an expected call of Bin
func (mr *MockEnvMockRecorder) Bin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bin", reflect.TypeOf((*MockEnv)(nil).Bin))
}

// Cache mocks base method
func (m *MockEnv) Cache() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cache")
	ret0, _ := ret[0].(string)
	return ret0
}

// Cache indicates an expected call of Cache
func (mr *MockEnvMockRecorder) Cache() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cache", reflect.TypeOf((*MockEnv)(nil).Cache))
}

// GithubHost mocks base method
func (m *MockEnv) GithubHost() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GithubHost")
	ret0, _ := ret[0].(string)
	return ret0
}

// GithubHost indicates an expected call of GithubHost
func (mr *MockEnvMockRecorder) GithubHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GithubHost", reflect.TypeOf((*MockEnv)(nil).GithubHost))
}

// GithubUser mocks base method
func (m *MockEnv) GithubUser() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GithubUser")
	ret0, _ := ret[0].(string)
	return ret0
}

// GithubUser indicates an expected call of GithubUser
func (mr *MockEnvMockRecorder) GithubUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GithubUser", reflect.TypeOf((*MockEnv)(nil).GithubUser))
}

// Hooks mocks base method
func (m *MockEnv) Hooks() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hooks")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Hooks indicates an expected call of Hooks
func (mr *MockEnvMockRecorder) Hooks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hooks", reflect.TypeOf((*MockEnv)(nil).Hooks))
}

// Man mocks base method
func (m *MockEnv) Man() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Man")
	ret0, _ := ret[0].(string)
	return ret0
}

// Man indicates an expected call of Man
func (mr *MockEnvMockRecorder) Man() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Man", reflect.TypeOf((*MockEnv)(nil).Man))
}

// OS mocks base method
func (m *MockEnv) OS() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OS")
	ret0, _ := ret[0].(string)
	return ret0
}

// OS indicates an expected call of OS
func (mr *MockEnvMockRecorder) OS() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OS", reflect.TypeOf((*MockEnv)(nil).OS))
}

// Root mocks base method
func (m *MockEnv) Root() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Root")
	ret0, _ := ret[0].(string)
	return ret0
}

// Root indicates an expected call of Root
func (mr *MockEnvMockRecorder) Root() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Root", reflect.TypeOf((*MockEnv)(nil).Root))
}
