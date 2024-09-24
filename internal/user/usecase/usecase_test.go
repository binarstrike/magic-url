package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/entity"
	"github.com/binarstrike/magic-url/internal/model"
	user_mock "github.com/binarstrike/magic-url/internal/user/mock"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user_mock.NewMockUserRepository(ctrl)
	mockRedisRepo := user_mock.NewMockRedisRepository(ctrl)
	_userUseCase := NewUserUseCase(mockUserRepo, mockRedisRepo, logger, cfg)

	userRequest := &model.RegisterUserRequest{
		Username: "ucup123",
		Email:    "ucup.128@xyz.com",
		Password: "foo123456",
	}

	userMock := &entity.User{
		Username: "ucup123",
		Email:    "ucup.128@xyz.com",
		Password: "foo123456",
	}

	mockUserRepo.EXPECT().FindByEmail(gomock.Any(), gomock.Eq(userRequest.Email)).Return(nil, nil)
	mockUserRepo.EXPECT().Register(gomock.Any(), gomock.AssignableToTypeOf(userMock)).Return(userMock, nil)

	userResponse, err := _userUseCase.Register(context.Background(), userRequest)
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, userResponse)
	assert.Equal(t, userResponse.Username, userRequest.Username)
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user_mock.NewMockUserRepository(ctrl)
	mockRedisRepo := user_mock.NewMockRedisRepository(ctrl)
	_userUseCase := NewUserUseCase(mockUserRepo, mockRedisRepo, logger, cfg)

	userRequest := &model.LoginUserRequest{
		Email:    "budi.256@xyz.com",
		Password: "bar123456",
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	now := time.Now()

	mockUser := &entity.User{
		UserId:    uuid.New(),
		Username:  "budi128",
		Email:     "budi.256@xyz.com",
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUserRepo.EXPECT().FindByEmail(gomock.Any(), gomock.Eq(userRequest.Email)).Return(mockUser, nil)

	userResponse, err := _userUseCase.Login(context.Background(), userRequest)
	assert.NoError(t, err)
	require.NotNil(t, userResponse)
}

func TestAuthService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user_mock.NewMockUserRepository(ctrl)
	mockRedisRepo := user_mock.NewMockRedisRepository(ctrl)
	_userUseCase := NewUserUseCase(mockUserRepo, mockRedisRepo, logger, cfg)

	userId := uuid.New()

	userRequest := &model.DeleteUserRequest{
		UserId: userId.String(),
	}

	userMock := &entity.User{UserId: userId}

	mockUserRepo.EXPECT().GetById(gomock.Any(), userRequest.UserId).Return(userMock, nil)
	mockUserRepo.EXPECT().Delete(gomock.Any(), userRequest.UserId).Return(nil)
	mockRedisRepo.EXPECT().DeleteUser(gomock.Any(), createCacheKey(userRequest.UserId)).Return(nil)

	err := _userUseCase.Delete(context.Background(), userRequest)
	assert.NoError(t, err)
	assert.Nil(t, err)
}

func TestAuthService_GetById(t *testing.T) {
	// it should get user from redis
	t.Run("test_1", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := user_mock.NewMockUserRepository(ctrl)
		mockRedisRepo := user_mock.NewMockRedisRepository(ctrl)
		_userUseCase := NewUserUseCase(mockUserRepo, mockRedisRepo, logger, cfg)

		userMock := &entity.User{UserId: uuid.New()}

		mockRedisRepo.EXPECT().GetUser(gomock.Any(), createCacheKey(userMock.UserId.String())).Return(userMock, nil)

		userResponse, err := _userUseCase.GetById(context.Background(), userMock.UserId.String())
		assert.NoError(t, err)
		assert.NotNil(t, userResponse)
	})

	// it should get user from database if redis cache is empty
	t.Run("test_2", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := user_mock.NewMockUserRepository(ctrl)
		mockRedisRepo := user_mock.NewMockRedisRepository(ctrl)
		_userUseCase := NewUserUseCase(mockUserRepo, mockRedisRepo, logger, cfg)

		userMock := &entity.User{UserId: uuid.New()}

		cacheKey := createCacheKey(userMock.UserId.String())

		mockRedisRepo.EXPECT().GetUser(gomock.Any(), cacheKey).Return(nil, nil)
		mockUserRepo.EXPECT().GetById(gomock.Any(), userMock.UserId.String()).Return(userMock, nil)
		mockRedisRepo.EXPECT().SetUser(gomock.Any(), cacheKey, gomock.Eq(userMock)).Return(nil)

		userResponse, err := _userUseCase.GetById(context.Background(), userMock.UserId.String())
		assert.NoError(t, err)
		assert.NotNil(t, userResponse)
	})
}
