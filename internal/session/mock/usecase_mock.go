// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source usecase.go -destination mock/usecase_mock.go -package session_mock
//

// Package session_mock is a generated GoMock package.
package session_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/binarstrike/magic-url/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockSessionUseCase is a mock of SessionUseCase interface.
type MockSessionUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSessionUseCaseMockRecorder
}

// MockSessionUseCaseMockRecorder is the mock recorder for MockSessionUseCase.
type MockSessionUseCaseMockRecorder struct {
	mock *MockSessionUseCase
}

// NewMockSessionUseCase creates a new mock instance.
func NewMockSessionUseCase(ctrl *gomock.Controller) *MockSessionUseCase {
	mock := &MockSessionUseCase{ctrl: ctrl}
	mock.recorder = &MockSessionUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionUseCase) EXPECT() *MockSessionUseCaseMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSessionUseCase) CreateSession(ctx context.Context, userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSessionUseCaseMockRecorder) CreateSession(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessionUseCase)(nil).CreateSession), ctx, userId)
}

// DeleteById mocks base method.
func (m *MockSessionUseCase) DeleteById(ctx context.Context, sessionId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, sessionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockSessionUseCaseMockRecorder) DeleteById(ctx, sessionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockSessionUseCase)(nil).DeleteById), ctx, sessionId)
}

// GetSessionById mocks base method.
func (m *MockSessionUseCase) GetSessionById(ctx context.Context, sessionId string) (*model.SessionData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionById", ctx, sessionId)
	ret0, _ := ret[0].(*model.SessionData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionById indicates an expected call of GetSessionById.
func (mr *MockSessionUseCaseMockRecorder) GetSessionById(ctx, sessionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionById", reflect.TypeOf((*MockSessionUseCase)(nil).GetSessionById), ctx, sessionId)
}
