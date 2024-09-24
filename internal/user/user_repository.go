package user

import (
	"context"

	"github.com/binarstrike/magic-url/internal/entity"
)

type UserRepository interface {
	Register(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userId string) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	GetById(ctx context.Context, userId string) (*entity.User, error)
}

//go:generate mockgen -source user_repository.go -destination mock/user_repository_mock.go -package user_mock
