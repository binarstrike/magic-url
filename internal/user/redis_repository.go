package user

import (
	"context"

	"github.com/binarstrike/magic-url/internal/entity"
)

type RedisRepository interface {
	GetUser(ctx context.Context, key string) (*entity.User, error)
	SetUser(ctx context.Context, key string, user *entity.User) error
	DeleteUser(ctx context.Context, key string) error
}

//go:generate mockgen -source redis_repository.go -destination mock/redis_repository_mock.go -package user_mock
