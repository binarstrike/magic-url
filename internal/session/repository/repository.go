package repository

import (
	"context"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/model"
	"github.com/binarstrike/magic-url/internal/session"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/binarstrike/magic-url/pkg/merror"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const prefixKey string = "session_"

type sessionRepository struct {
	redisClient *redis.Client
	log         lojer.Logger
	cfg         *config.Config
}

func NewSessionRepository(client *redis.Client, logger lojer.Logger, config *config.Config) session.SessionRepository {
	return &sessionRepository{redisClient: client, log: logger, cfg: config}
}

func (sr *sessionRepository) CreateSession(ctx context.Context, userId string) (string, error) {
	var session model.SessionData

	sessionId := uuid.NewString()

	session.UserId = userId
	session.Authenticated = true

	key := createCacheKey(sessionId)

	err := sr.redisClient.HSet(ctx, key, session).Err()
	if err != nil {
		sr.log.Error(merror.InternalRedisError, zap.Error(err))
		return "", merror.New(merror.InternalRedisError)
	}

	err = sr.redisClient.Expire(ctx, key, sr.cfg.Consts.SessionExpirationTime).Err()
	if err != nil {
		sr.log.Error(merror.InternalRedisError, zap.Error(err))
		return "", merror.New(merror.InternalRedisError)
	}

	return sessionId, nil
}

func (sr *sessionRepository) GetSessionById(ctx context.Context, sessionId string) (*model.SessionData, error) {
	m := sr.redisClient.HGetAll(ctx, createCacheKey(sessionId))
	if result, err := m.Result(); err != nil {
		if err == redis.Nil || len(result) < 1 {
			return nil, nil
		} else {
			sr.log.Error(merror.InternalRedisError, zap.Error(err))
			return nil, merror.New(merror.InternalRedisError)
		}
	}

	session := new(model.SessionData)

	if err := m.Scan(session); err != nil {
		sr.log.Error(merror.InternalRedisError, zap.Error(err))
		return nil, merror.New(merror.InternalRedisError)
	}

	return session, nil
}

func (sr *sessionRepository) DeleteById(ctx context.Context, sessionId string) error {
	if err := sr.redisClient.Del(ctx, createCacheKey(sessionId)).Err(); err != nil {
		sr.log.Error(merror.InternalRedisError, zap.Error(err))
		return merror.New(merror.InternalRedisError)
	}

	return nil
}

func createCacheKey(sessionId string) string {
	return prefixKey + sessionId
}
