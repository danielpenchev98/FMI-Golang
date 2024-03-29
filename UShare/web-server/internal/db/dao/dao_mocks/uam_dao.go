// Code generated by MockGen. DO NOT EDIT.
// Source: uam_dao.go

// Package dao_mocks is a generated GoMock package.
package dao_mocks

import (
	models "github.com/danielpenchev98/FMI-Golang/UShare/web-server/internal/db/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUamDAO is a mock of UamDAO interface
type MockUamDAO struct {
	ctrl     *gomock.Controller
	recorder *MockUamDAOMockRecorder
}

// MockUamDAOMockRecorder is the mock recorder for MockUamDAO
type MockUamDAOMockRecorder struct {
	mock *MockUamDAO
}

// NewMockUamDAO creates a new mock instance
func NewMockUamDAO(ctrl *gomock.Controller) *MockUamDAO {
	mock := &MockUamDAO{ctrl: ctrl}
	mock.recorder = &MockUamDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUamDAO) EXPECT() *MockUamDAOMockRecorder {
	return m.recorder
}

// Migrate mocks base method
func (m *MockUamDAO) Migrate() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Migrate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Migrate indicates an expected call of Migrate
func (mr *MockUamDAOMockRecorder) Migrate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockUamDAO)(nil).Migrate))
}

// CreateUser mocks base method
func (m *MockUamDAO) CreateUser(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockUamDAOMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUamDAO)(nil).CreateUser), arg0, arg1)
}

// GetUser mocks base method
func (m *MockUamDAO) GetUser(arg0 string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockUamDAOMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUamDAO)(nil).GetUser), arg0)
}

// DeleteUser mocks base method
func (m *MockUamDAO) DeleteUser(arg0 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockUamDAOMockRecorder) DeleteUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUamDAO)(nil).DeleteUser), arg0)
}

// CreateGroup mocks base method
func (m *MockUamDAO) CreateGroup(arg0 uint, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateGroup indicates an expected call of CreateGroup
func (mr *MockUamDAOMockRecorder) CreateGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockUamDAO)(nil).CreateGroup), arg0, arg1)
}

// AddUserToGroup mocks base method
func (m *MockUamDAO) AddUserToGroup(arg0 uint, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToGroup", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToGroup indicates an expected call of AddUserToGroup
func (mr *MockUamDAOMockRecorder) AddUserToGroup(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToGroup", reflect.TypeOf((*MockUamDAO)(nil).AddUserToGroup), arg0, arg1, arg2)
}

// RemoveUserFromGroup mocks base method
func (m *MockUamDAO) RemoveUserFromGroup(arg0 uint, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserFromGroup", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserFromGroup indicates an expected call of RemoveUserFromGroup
func (mr *MockUamDAOMockRecorder) RemoveUserFromGroup(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserFromGroup", reflect.TypeOf((*MockUamDAO)(nil).RemoveUserFromGroup), arg0, arg1, arg2)
}

// MemberExists mocks base method
func (m *MockUamDAO) MemberExists(arg0, arg1 uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MemberExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MemberExists indicates an expected call of MemberExists
func (mr *MockUamDAOMockRecorder) MemberExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MemberExists", reflect.TypeOf((*MockUamDAO)(nil).MemberExists), arg0, arg1)
}

// DeactivateGroup mocks base method
func (m *MockUamDAO) DeactivateGroup(arg0 uint, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeactivateGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeactivateGroup indicates an expected call of DeactivateGroup
func (mr *MockUamDAOMockRecorder) DeactivateGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeactivateGroup", reflect.TypeOf((*MockUamDAO)(nil).DeactivateGroup), arg0, arg1)
}

// GetGroup mocks base method
func (m *MockUamDAO) GetGroup(arg0 string) (models.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", arg0)
	ret0, _ := ret[0].(models.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup
func (mr *MockUamDAOMockRecorder) GetGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockUamDAO)(nil).GetGroup), arg0)
}

// GetDeactivatedGroupNames mocks base method
func (m *MockUamDAO) GetDeactivatedGroupNames() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeactivatedGroupNames")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeactivatedGroupNames indicates an expected call of GetDeactivatedGroupNames
func (mr *MockUamDAOMockRecorder) GetDeactivatedGroupNames() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeactivatedGroupNames", reflect.TypeOf((*MockUamDAO)(nil).GetDeactivatedGroupNames))
}

// EraseDeactivatedGroups mocks base method
func (m *MockUamDAO) EraseDeactivatedGroups(arg0 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EraseDeactivatedGroups", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// EraseDeactivatedGroups indicates an expected call of EraseDeactivatedGroups
func (mr *MockUamDAOMockRecorder) EraseDeactivatedGroups(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EraseDeactivatedGroups", reflect.TypeOf((*MockUamDAO)(nil).EraseDeactivatedGroups), arg0)
}

// GetAllGroups mocks base method
func (m *MockUamDAO) GetAllGroups() ([]models.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllGroups")
	ret0, _ := ret[0].([]models.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllGroups indicates an expected call of GetAllGroups
func (mr *MockUamDAOMockRecorder) GetAllGroups() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllGroups", reflect.TypeOf((*MockUamDAO)(nil).GetAllGroups))
}

// GetAllUsers mocks base method
func (m *MockUamDAO) GetAllUsers() ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers")
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers
func (mr *MockUamDAOMockRecorder) GetAllUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockUamDAO)(nil).GetAllUsers))
}

// GetAllUsersInGroup mocks base method
func (m *MockUamDAO) GetAllUsersInGroup(arg0 uint, arg1 string) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsersInGroup", arg0, arg1)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsersInGroup indicates an expected call of GetAllUsersInGroup
func (mr *MockUamDAOMockRecorder) GetAllUsersInGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsersInGroup", reflect.TypeOf((*MockUamDAO)(nil).GetAllUsersInGroup), arg0, arg1)
}
