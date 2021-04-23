// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/typical-go/typical-rest-server/internal/generated/dbrepo (interfaces: SongRepo)

// Package dbrepo_mock is a generated GoMock package.
package dbrepo_mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entity "github.com/typical-go/typical-rest-server/internal/app/entity"
	sqkit "github.com/typical-go/typical-rest-server/pkg/sqkit"
	reflect "reflect"
)

// MockSongRepo is a mock of SongRepo interface
type MockSongRepo struct {
	ctrl     *gomock.Controller
	recorder *MockSongRepoMockRecorder
}

// MockSongRepoMockRecorder is the mock recorder for MockSongRepo
type MockSongRepoMockRecorder struct {
	mock *MockSongRepo
}

// NewMockSongRepo creates a new mock instance
func NewMockSongRepo(ctrl *gomock.Controller) *MockSongRepo {
	mock := &MockSongRepo{ctrl: ctrl}
	mock.recorder = &MockSongRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSongRepo) EXPECT() *MockSongRepoMockRecorder {
	return m.recorder
}

// BulkInsert mocks base method
func (m *MockSongRepo) BulkInsert(arg0 context.Context, arg1 ...*entity.Song) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BulkInsert", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BulkInsert indicates an expected call of BulkInsert
func (mr *MockSongRepoMockRecorder) BulkInsert(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkInsert", reflect.TypeOf((*MockSongRepo)(nil).BulkInsert), varargs...)
}

// Count mocks base method
func (m *MockSongRepo) Count(arg0 context.Context, arg1 ...sqkit.SelectOption) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Count", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockSongRepoMockRecorder) Count(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockSongRepo)(nil).Count), varargs...)
}

// Delete mocks base method
func (m *MockSongRepo) Delete(arg0 context.Context, arg1 sqkit.DeleteOption) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete
func (mr *MockSongRepoMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSongRepo)(nil).Delete), arg0, arg1)
}

// Find mocks base method
func (m *MockSongRepo) Find(arg0 context.Context, arg1 ...sqkit.SelectOption) ([]*entity.Song, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Find", varargs...)
	ret0, _ := ret[0].([]*entity.Song)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockSongRepoMockRecorder) Find(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockSongRepo)(nil).Find), varargs...)
}

// Insert mocks base method
func (m *MockSongRepo) Insert(arg0 context.Context, arg1 *entity.Song) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert
func (mr *MockSongRepoMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockSongRepo)(nil).Insert), arg0, arg1)
}

// Patch mocks base method
func (m *MockSongRepo) Patch(arg0 context.Context, arg1 *entity.Song, arg2 sqkit.UpdateOption) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Patch", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Patch indicates an expected call of Patch
func (mr *MockSongRepoMockRecorder) Patch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockSongRepo)(nil).Patch), arg0, arg1, arg2)
}

// Update mocks base method
func (m *MockSongRepo) Update(arg0 context.Context, arg1 *entity.Song, arg2 sqkit.UpdateOption) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockSongRepoMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSongRepo)(nil).Update), arg0, arg1, arg2)
}