// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	model "optionhub-service/internal/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDBRepo is a mock of DBRepo interface.
type MockDBRepo struct {
	ctrl     *gomock.Controller
	recorder *MockDBRepoMockRecorder
}

// MockDBRepoMockRecorder is the mock recorder for MockDBRepo.
type MockDBRepoMockRecorder struct {
	mock *MockDBRepo
}

// NewMockDBRepo creates a new mock instance.
func NewMockDBRepo(ctrl *gomock.Controller) *MockDBRepo {
	mock := &MockDBRepo{ctrl: ctrl}
	mock.recorder = &MockDBRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBRepo) EXPECT() *MockDBRepoMockRecorder {
	return m.recorder
}

// AddOS mocks base method.
func (m *MockDBRepo) AddOS(ctx context.Context, name, uuid string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOS", ctx, name, uuid)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOS indicates an expected call of AddOS.
func (mr *MockDBRepoMockRecorder) AddOS(ctx, name, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOS", reflect.TypeOf((*MockDBRepo)(nil).AddOS), ctx, name, uuid)
}

// GetAllOs mocks base method.
func (m *MockDBRepo) GetAllOs() (*model.OSList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllOs")
	ret0, _ := ret[0].(*model.OSList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllOs indicates an expected call of GetAllOs.
func (mr *MockDBRepoMockRecorder) GetAllOs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllOs", reflect.TypeOf((*MockDBRepo)(nil).GetAllOs))
}

// GetOsByID mocks base method.
func (m *MockDBRepo) GetOsByID(ctx context.Context, id int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOsByID", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOsByID indicates an expected call of GetOsByID.
func (mr *MockDBRepoMockRecorder) GetOsByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOsByID", reflect.TypeOf((*MockDBRepo)(nil).GetOsByID), ctx, id)
}

// GetOsBySearchName mocks base method.
func (m *MockDBRepo) GetOsBySearchName(ctx context.Context, name string) (*model.OSList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOsBySearchName", ctx, name)
	ret0, _ := ret[0].(*model.OSList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOsBySearchName indicates an expected call of GetOsBySearchName.
func (mr *MockDBRepoMockRecorder) GetOsBySearchName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOsBySearchName", reflect.TypeOf((*MockDBRepo)(nil).GetOsBySearchName), ctx, name)
}
