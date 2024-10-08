// Code generated by MockGen. DO NOT EDIT.
// Source: user_usecase.go
//
// Generated by this command:
//
//	mockgen -source user_usecase.go -destination mock/user_usecase_mock.go -package user_mock
//

// Package user_mock is a generated GoMock package.
package user_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/binarstrike/magic-url/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockUserUseCase is a mock of UserUseCase interface.
type MockUserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUseCaseMockRecorder
}

// MockUserUseCaseMockRecorder is the mock recorder for MockUserUseCase.
type MockUserUseCaseMockRecorder struct {
	mock *MockUserUseCase
}

// NewMockUserUseCase creates a new mock instance.
func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
	mock := &MockUserUseCase{ctrl: ctrl}
	mock.recorder = &MockUserUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUserUseCase) Delete(ctx context.Context, request *model.DeleteUserRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserUseCaseMockRecorder) Delete(ctx, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserUseCase)(nil).Delete), ctx, request)
}

// GetById mocks base method.
func (m *MockUserUseCase) GetById(ctx context.Context, userId string) (*model.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, userId)
	ret0, _ := ret[0].(*model.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUserUseCaseMockRecorder) GetById(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUserUseCase)(nil).GetById), ctx, userId)
}

// Login mocks base method.
func (m *MockUserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, request)
	ret0, _ := ret[0].(*model.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserUseCaseMockRecorder) Login(ctx, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserUseCase)(nil).Login), ctx, request)
}

// Register mocks base method.
func (m *MockUserUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, request)
	ret0, _ := ret[0].(*model.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockUserUseCaseMockRecorder) Register(ctx, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserUseCase)(nil).Register), ctx, request)
}
