// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/s21platform/optionhub-service/internal/model"
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

// GetAttributeValueById mocks base method.
func (m *MockDBRepo) GetAttributeValueById(ctx context.Context, ids []int64) ([]model.Attribute, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttributeValueById", ctx, ids)
	ret0, _ := ret[0].([]model.Attribute)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttributeValueById indicates an expected call of GetAttributeValueById.
func (mr *MockDBRepoMockRecorder) GetAttributeValueById(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttributeValueById", reflect.TypeOf((*MockDBRepo)(nil).GetAttributeValueById), ctx, ids)
}

// GetOptionRequests mocks base method.
func (m *MockDBRepo) GetOptionRequests(ctx context.Context) (model.OptionRequestList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOptionRequests", ctx)
	ret0, _ := ret[0].(model.OptionRequestList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOptionRequests indicates an expected call of GetOptionRequests.
func (mr *MockDBRepoMockRecorder) GetOptionRequests(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOptionRequests", reflect.TypeOf((*MockDBRepo)(nil).GetOptionRequests), ctx)
}
