package repository

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/entity"
	"github.com/binarstrike/magic-url/internal/user"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var redisRepo user.RedisRepository

func TestMain(m *testing.M) {
	mrd, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: mrd.Addr(),
	})

	cfg, _ := config.InitConfig(true)
	redisRepo = NewRedisRepository(redisClient, cfg, lojer.NewFromZap(zap.NewNop()))

	m.Run()
}

func TestRedisRepository_SetUser(t *testing.T) {
	ctx := context.Background()

	user := &entity.User{
		UserId: uuid.New(),
	}

	err := redisRepo.SetUser(ctx, createCacheKey(user.UserId.String()), user)
	assert.NoError(t, err)
}

func TestRedisRepository_GetUser(t *testing.T) {
	ctx := context.Background()

	user := &entity.User{
		UserId:   uuid.New(),
		Username: "budi128",
	}

	cacheKey := createCacheKey(user.UserId.String())

	err := redisRepo.SetUser(ctx, cacheKey, user)
	assert.NoError(t, err)

	_user, err := redisRepo.GetUser(ctx, cacheKey)
	assert.NoError(t, err)
	assert.NotNil(t, _user)
	assert.Equal(t, user.UserId, _user.UserId)
}

func TestRedisRepository_DeleteUser(t *testing.T) {
	ctx := context.Background()

	user := &entity.User{
		UserId:   uuid.New(),
		Username: "budi128",
	}

	key := createCacheKey(user.UserId.String())

	err := redisRepo.SetUser(ctx, key, user)
	assert.NoError(t, err)

	err = redisRepo.DeleteUser(ctx, key)
	assert.NoError(t, err)

	_user, err := redisRepo.GetUser(ctx, key)
	assert.NoError(t, err)
	assert.Nil(t, _user)
}

func createCacheKey(userId string) string {
	return "test_" + userId
}
