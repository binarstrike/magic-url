package session

import (
	"context"

	"github.com/binarstrike/magic-url/internal/model"
)

type SessionUseCase interface {
	GetSessionById(ctx context.Context, sessionId string) (*model.SessionData, error)
	CreateSession(ctx context.Context, userId string) (string, error)
	DeleteById(ctx context.Context, sessionId string) error
}

//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package session_mock
