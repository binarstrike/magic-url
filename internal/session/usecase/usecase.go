package usecase

import (
	"context"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/model"
	"github.com/binarstrike/magic-url/internal/session"
	"github.com/binarstrike/magic-url/pkg/lojer"
)

type sessionUseCase struct {
	redisRepository session.SessionRepository
	cfg             *config.Config
}

func NewSessionUseCase(redisRepository session.SessionRepository, logger lojer.Logger, config *config.Config) session.SessionUseCase {
	return &sessionUseCase{redisRepository: redisRepository, cfg: config}
}

func (s sessionUseCase) CreateSession(ctx context.Context, userId string) (string, error) {
	return s.redisRepository.CreateSession(ctx, userId)
}

func (s sessionUseCase) GetSessionById(ctx context.Context, sessionId string) (*model.SessionData, error) {
	return s.redisRepository.GetSessionById(ctx, sessionId)
}

func (s sessionUseCase) DeleteById(ctx context.Context, sessionId string) error {
	return s.redisRepository.DeleteById(ctx, sessionId)
}
