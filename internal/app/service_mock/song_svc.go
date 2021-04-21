// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/typical-go/typical-rest-server/internal/app/service (interfaces: SongSvc)

// Package service_mock is a generated GoMock package.
package service_mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entity "github.com/typical-go/typical-rest-server/internal/app/entity"
	service "github.com/typical-go/typical-rest-server/internal/app/service"
	reflect "reflect"
)

// MockSongSvc is a mock of SongSvc interface
type MockSongSvc struct {
	ctrl     *gomock.Controller
	recorder *MockSongSvcMockRecorder
}

// MockSongSvcMockRecorder is the mock recorder for MockSongSvc
type MockSongSvcMockRecorder struct {
	mock *MockSongSvc
}

// NewMockSongSvc creates a new mock instance
func NewMockSongSvc(ctrl *gomock.Controller) *MockSongSvc {
	mock := &MockSongSvc{ctrl: ctrl}
	mock.recorder = &MockSongSvcMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSongSvc) EXPECT() *MockSongSvcMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockSongSvc) Create(arg0 context.Context, arg1 *entity.Song) (*entity.Song, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*entity.Song)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockSongSvcMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSongSvc)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockSongSvc) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockSongSvcMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSongSvc)(nil).Delete), arg0, arg1)
}

// Find mocks base method
func (m *MockSongSvc) Find(arg0 context.Context, arg1 *service.FindSongReq) (*service.FindSongResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*service.FindSongResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockSongSvcMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockSongSvc)(nil).Find), arg0, arg1)
}

// FindOne mocks base method
func (m *MockSongSvc) FindOne(arg0 context.Context, arg1 string) (*entity.Song, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", arg0, arg1)
	ret0, _ := ret[0].(*entity.Song)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne
func (mr *MockSongSvcMockRecorder) FindOne(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockSongSvc)(nil).FindOne), arg0, arg1)
}

// Patch mocks base method
func (m *MockSongSvc) Patch(arg0 context.Context, arg1 string, arg2 *entity.Song) (*entity.Song, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Patch", arg0, arg1, arg2)
	ret0, _ := ret[0].(*entity.Song)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Patch indicates an expected call of Patch
func (mr *MockSongSvcMockRecorder) Patch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockSongSvc)(nil).Patch), arg0, arg1, arg2)
}

// Update mocks base method
func (m *MockSongSvc) Update(arg0 context.Context, arg1 string, arg2 *entity.Song) (*entity.Song, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*entity.Song)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockSongSvcMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSongSvc)(nil).Update), arg0, arg1, arg2)
}
