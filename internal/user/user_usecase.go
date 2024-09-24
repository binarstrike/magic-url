package user

import (
	"context"

	"github.com/binarstrike/magic-url/internal/model"
)

// TODO: implementasi fungsi update?

type UserUseCase interface {
	Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, request *model.DeleteUserRequest) error
	GetById(ctx context.Context, userId string) (*model.UserResponse, error)
	// Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.AuthUser, error)
	// Update(ctx context.Context, *entity.User) (*entity.User, error)
}

//go:generate mockgen -source user_usecase.go -destination mock/user_usecase_mock.go -package user_mock
