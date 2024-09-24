package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/entity"
	"github.com/binarstrike/magic-url/internal/middleware"
	"github.com/binarstrike/magic-url/internal/model"
	"github.com/binarstrike/magic-url/internal/model/converter"
	session_mock "github.com/binarstrike/magic-url/internal/session/mock"
	user_mock "github.com/binarstrike/magic-url/internal/user/mock"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/binarstrike/magic-url/pkg/validator"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var (
	cfg        *config.Config
	logger     lojer.Logger
	_validator *validator.Validator
)

func TestMain(m *testing.M) {
	cfg, _ = config.InitConfig(true)
	logger = lojer.NewFromZap(zap.NewNop())
	_validator = validator.NewValidator()

	m.Run()
}

func TestUserController_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := user_mock.NewMockUserUseCase(ctrl)
	mockSessionUseCase := session_mock.NewMockSessionUseCase(ctrl)

	userController := NewUserController(mockUserUseCase, mockSessionUseCase, logger, _validator, cfg)

	app := fiber.New()
	app.Post(UserRegisterRoute, userController.Register)

	userMock := &entity.User{
		UserId:   uuid.New(),
		Username: "otong128",
		Email:    "example.86@xyz.com",
		Password: "Wkwkwkwkwk",
	}

	userRequest := &model.RegisterUserRequest{
		Username: userMock.Username,
		Email:    userMock.Email,
		Password: userMock.Password,
	}

	userResponse := converter.UserToResponse(userMock)

	sessionId := uuid.NewString()
	userMockEncoded, err := json.Marshal(userMock)
	assert.NoError(t, err)

	req := httptest.NewRequest(fiber.MethodPost, UserRegisterRoute, bytes.NewReader(userMockEncoded))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	mockUserUseCase.EXPECT().Register(gomock.Any(), gomock.Eq(userRequest)).Return(userResponse, nil)
	mockSessionUseCase.EXPECT().CreateSession(gomock.Any(), userMock.UserId.String()).Return(sessionId, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, fiber.StatusCreated)
	require.NotEqual(t, resp.Body, http.NoBody)

	response := new(model.UserResponse)
	err = json.NewDecoder(resp.Body).Decode(response)
	assert.NoError(t, err)
	assert.Equal(t, response.Username, userRequest.Username)
}

func TestUserController_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := user_mock.NewMockUserUseCase(ctrl)
	mockSessionUseCase := session_mock.NewMockSessionUseCase(ctrl)

	userController := NewUserController(mockUserUseCase, mockSessionUseCase, logger, _validator, cfg)

	app := fiber.New()
	app.Post(UserLoginRoute, userController.Login)

	userMock := &entity.User{
		UserId:   uuid.New(),
		Username: "otong128",
		Email:    "example.86@xyz.com",
		Password: "Wkwkwkwkwk",
	}

	userRequest := &model.LoginUserRequest{
		Email:    userMock.Email,
		Password: userMock.Password,
	}

	userResponse := converter.UserToResponse(userMock)

	sessionId := uuid.NewString()
	userMockEncoded, err := json.Marshal(userMock)
	assert.NoError(t, err)

	req := httptest.NewRequest(fiber.MethodPost, UserLoginRoute, bytes.NewReader(userMockEncoded))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	mockUserUseCase.EXPECT().Login(gomock.Any(), gomock.Eq(userRequest)).Return(userResponse, nil)
	mockSessionUseCase.EXPECT().CreateSession(gomock.Any(), userMock.UserId.String()).Return(sessionId, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, fiber.StatusOK)
	require.NotEqual(t, resp.Body, http.NoBody)

	response := new(model.UserResponse)
	err = json.NewDecoder(resp.Body).Decode(response)
	assert.NoError(t, err)
}

func TestUserController_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := user_mock.NewMockUserUseCase(ctrl)
	mockSessionUseCase := session_mock.NewMockSessionUseCase(ctrl)
	mMiddleware := middleware.NewMiddleware(mockUserUseCase, mockSessionUseCase, _validator, cfg)

	userController := NewUserController(mockUserUseCase, mockSessionUseCase, logger, _validator, cfg)

	app := fiber.New()
	app.Post(UserLogoutRoute, mMiddleware.AuthSessionMiddleware, userController.Logout)

	sessionId := uuid.NewString()
	userId := uuid.NewString()

	mockSessionData := &model.SessionData{UserId: userId, Authenticated: true}

	req := httptest.NewRequest(fiber.MethodPost, UserLogoutRoute, nil)
	req.AddCookie(&http.Cookie{
		Name:  cfg.Consts.SessionCookieName,
		Value: sessionId,
	})

	mockSessionUseCase.EXPECT().GetSessionById(gomock.Any(), sessionId).Return(mockSessionData, nil)
	mockSessionUseCase.EXPECT().DeleteById(gomock.Any(), sessionId).Return(nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, fiber.StatusNoContent)
	assert.Equal(t, resp.Body, http.NoBody)

	cookies := resp.Cookies()
	require.GreaterOrEqual(t, len(cookies), 1)
	require.NotNil(t, cookies[0])
	assert.Equal(t, cookies[0].Name, cfg.Consts.SessionCookieName)
	assert.Empty(t, cookies[0].Value)
}

func TestUserController_Me(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := user_mock.NewMockUserUseCase(ctrl)
	mockSessionUseCase := session_mock.NewMockSessionUseCase(ctrl)
	mMiddleware := middleware.NewMiddleware(mockUserUseCase, mockSessionUseCase, _validator, cfg)

	userController := NewUserController(mockUserUseCase, mockSessionUseCase, logger, _validator, cfg)

	app := fiber.New()
	app.Get(UserCurrentRoute, mMiddleware.AuthSessionMiddleware, userController.Me)

	sessionId := uuid.NewString()
	userId := uuid.New()

	mockSessionData := &model.SessionData{UserId: userId.String(), Authenticated: true}
	mockUser := &entity.User{UserId: userId}
	userResponse := converter.UserToResponse(mockUser)

	req := httptest.NewRequest(fiber.MethodGet, UserCurrentRoute, nil)
	req.AddCookie(&http.Cookie{
		Name:  cfg.Consts.SessionCookieName,
		Value: sessionId,
	})

	mockSessionUseCase.EXPECT().GetSessionById(gomock.Any(), sessionId).Return(mockSessionData, nil)
	mockUserUseCase.EXPECT().GetById(gomock.Any(), mockSessionData.UserId).Return(userResponse, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Body)

	response := new(model.UserResponse)
	err = json.NewDecoder(resp.Body).Decode(response)
	assert.NoError(t, err)
	assert.Empty(t, response.UserId)
}

func TestUserController_Delete(t *testing.T) {
	// TODO: implement me!!
}
