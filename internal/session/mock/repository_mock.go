// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source repository.go -destination mock/repository_mock.go -package session_mock
//

// Package session_mock is a generated GoMock package.
package session_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/binarstrike/magic-url/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockSessionRepository is a mock of SessionRepository interface.
type MockSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepositoryMockRecorder
}

// MockSessionRepositoryMockRecorder is the mock recorder for MockSessionRepository.
type MockSessionRepositoryMockRecorder struct {
	mock *MockSessionRepository
}

// NewMockSessionRepository creates a new mock instance.
func NewMockSessionRepository(ctrl *gomock.Controller) *MockSessionRepository {
	mock := &MockSessionRepository{ctrl: ctrl}
	mock.recorder = &MockSessionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRepository) EXPECT() *MockSessionRepositoryMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSessionRepository) CreateSession(ctx context.Context, userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSessionRepositoryMockRecorder) CreateSession(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessionRepository)(nil).CreateSession), ctx, userId)
}

// DeleteById mocks base method.
func (m *MockSessionRepository) DeleteById(ctx context.Context, sessionId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, sessionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockSessionRepositoryMockRecorder) DeleteById(ctx, sessionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockSessionRepository)(nil).DeleteById), ctx, sessionId)
}

// GetSessionById mocks base method.
func (m *MockSessionRepository) GetSessionById(ctx context.Context, sessionId string) (*model.SessionData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionById", ctx, sessionId)
	ret0, _ := ret[0].(*model.SessionData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionById indicates an expected call of GetSessionById.
func (mr *MockSessionRepositoryMockRecorder) GetSessionById(ctx, sessionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionById", reflect.TypeOf((*MockSessionRepository)(nil).GetSessionById), ctx, sessionId)
}
