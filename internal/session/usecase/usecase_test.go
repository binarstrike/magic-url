package usecase

import (
	"context"
	"testing"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/model"
	session_mock "github.com/binarstrike/magic-url/internal/session/mock"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var (
	cfg    *config.Config
	logger lojer.Logger
)

func TestMain(m *testing.M) {
	cfg, _ = config.InitConfig(true)
	logger = lojer.NewFromZap(zap.NewNop())

	m.Run()
}

func TestSessionUseCase_CreateSession(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := session_mock.NewMockSessionRepository(ctrl)
	_sessionUseCase := NewSessionUseCase(mockSessionRepo, logger, cfg)

	userId := uuid.NewString()
	sessionId := uuid.NewString()

	mockSessionRepo.EXPECT().CreateSession(gomock.Any(), gomock.Eq(userId)).Return(sessionId, nil)

	_sessionId, err := _sessionUseCase.CreateSession(context.Background(), userId)
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotEmpty(t, _sessionId)
}

func TestSessionUseCase_GetSessionById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := session_mock.NewMockSessionUseCase(ctrl)
	_sessionUseCase := NewSessionUseCase(mockSessionRepo, logger, cfg)

	sessionId := uuid.NewString()
	session := &model.SessionData{
		UserId:        uuid.NewString(),
		Authenticated: true,
	}

	// * fungsi gomock.Eq pada baris dibawah di maksudkan agar saat memanggil fungsi dari layer service paramter pada fungsi harus
	// * sesuai nilai dan tipe data pada saat memanggil fungsi pada layer repository
	// * juga fungsi .Return pada fungsi mock di gunakan untuk memberi nilai kembalian pada saat memanggil fungsi dari layer repository
	// * di layer service sehingga test dapat di jalankan sesuai yang di harapkan
	mockSessionRepo.EXPECT().GetSessionById(gomock.Any(), gomock.Eq(sessionId)).Return(session, nil)

	_session, err := _sessionUseCase.GetSessionById(context.Background(), sessionId)
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, _session)
	assert.True(t, _session.Authenticated)
}

func TestSessionUseCase_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := session_mock.NewMockSessionRepository(ctrl)
	_sessionService := NewSessionUseCase(mockSessionRepo, logger, cfg)

	sessionId := uuid.NewString()

	mockSessionRepo.EXPECT().DeleteById(gomock.Any(), gomock.Eq(sessionId)).Return(nil)

	err := _sessionService.DeleteById(context.Background(), sessionId)
	assert.NoError(t, err)
	assert.Nil(t, err)
}
