package repository

import (
	"context"
	"errors"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/entity"
	"github.com/binarstrike/magic-url/internal/user"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/binarstrike/magic-url/pkg/merror"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type redisRepository struct {
	client *redis.Client
	log    lojer.Logger
	cfg    *config.Config
}

func NewRedisRepository(client *redis.Client, config *config.Config, logger lojer.Logger) user.RedisRepository {
	return &redisRepository{client: client, cfg: config, log: logger}
}

func (rr redisRepository) SetUser(ctx context.Context, key string, user *entity.User) error {
	if user == nil {
		err := errors.New("user object is nil")
		rr.log.Error(merror.InternalRedisError, zap.Error(err))
		return err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		rr.log.Error("marshaling error", zap.Error(err))
		return err
	}

	err = rr.client.Set(ctx, key, userBytes, rr.cfg.Consts.UserExpirationTime).Err()
	if err != nil {
		rr.log.Error(merror.InternalRedisError, zap.Error(err))
		return err
	}

	return nil
}

func (rr redisRepository) GetUser(ctx context.Context, key string) (*entity.User, error) {
	s := rr.client.Get(ctx, key)
	if result, err := s.Result(); err == redis.Nil || result == "" {
		return nil, nil
	} else if err != nil {
		rr.log.Error(merror.InternalRedisError, zap.Error(err))
		return nil, err
	}

	user := new(entity.User)

	err := s.Scan(user)
	if err != nil {
		rr.log.Error(merror.InternalRedisError, zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (rr redisRepository) DeleteUser(ctx context.Context, key string) error {
	err := rr.client.Del(ctx, key).Err()
	if err != nil {
		rr.log.Error(merror.InternalRedisError, zap.Error(err))
		return err
	}

	return nil
}
