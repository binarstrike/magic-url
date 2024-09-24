package repository

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/session"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var sessionRepo session.SessionRepository

func TestMain(m *testing.M) {
	mrd, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: mrd.Addr(),
	})

	cfg, _ := config.InitConfig(true)
	logger := lojer.NewFromZap(zap.NewNop())
	sessionRepo = NewSessionRepository(redisClient, logger, cfg)

	m.Run()
}

func TestSessionRedisRepository_CreateSession(t *testing.T) {
	userId := uuid.NewString()

	sessionId, err := sessionRepo.CreateSession(context.Background(), userId)
	assert.NoError(t, err)
	assert.NotEmpty(t, sessionId)
}

func TestSessionRedisRepository_GetSessionById(t *testing.T) {
	// get session by id
	t.Run("test_1", func(t *testing.T) {
		ctx := context.Background()
		userId := uuid.NewString()

		sessionId, err := sessionRepo.CreateSession(ctx, userId)
		assert.NoError(t, err)
		assert.NotEmpty(t, sessionId)

		session, err := sessionRepo.GetSessionById(ctx, sessionId)
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, userId, session.UserId)
	})

	// return empty if no session is stored
	t.Run("test_2", func(t *testing.T) {
		userId := uuid.NewString()

		session, err := sessionRepo.GetSessionById(context.Background(), userId)
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Empty(t, session)
	})
}

func TestSessionRedisRepository_DeleteById(t *testing.T) {
	ctx := context.Background()
	userId := uuid.NewString()

	sessionId, err := sessionRepo.CreateSession(ctx, userId)
	assert.NoError(t, err)
	require.NotEmpty(t, sessionId)

	err = sessionRepo.DeleteById(ctx, sessionId)
	assert.NoError(t, err)
}
