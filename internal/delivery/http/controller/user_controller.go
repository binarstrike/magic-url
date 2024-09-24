package controller

import (
	"time"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/model"
	"github.com/binarstrike/magic-url/internal/session"
	"github.com/binarstrike/magic-url/internal/user"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/binarstrike/magic-url/pkg/merror"
	"github.com/binarstrike/magic-url/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	BaseUserRoute     = "/api/v1/users"
	UserRegisterRoute = BaseUserRoute
	UserDeleteRoute   = BaseUserRoute + "/:userId"
	UserLoginRoute    = BaseUserRoute + "/login"
	UserLogoutRoute   = BaseUserRoute + "/logout"
	UserCurrentRoute  = BaseUserRoute + "/me"
)

type UserController struct {
	userUseCase    user.UserUseCase
	sessionUseCase session.SessionUseCase
	cfg            *config.Config
	log            lojer.Logger
	validator      *validator.Validator
}

func NewUserController(userUseCase user.UserUseCase, sessionUseCase session.SessionUseCase, logger lojer.Logger, validator *validator.Validator, config *config.Config) *UserController {
	return &UserController{
		userUseCase:    userUseCase,
		sessionUseCase: sessionUseCase,
		cfg:            config,
		log:            logger,
		validator:      validator,
	}
}

// Register
//
//	@Summary		Register a new user account
//	@Description	Register a new user account by sending a JSON payload.
//	@Tags			registration
//	@Router			/users [post]
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		model.RegisterUserRequest	true	"Register new user account"
//	@Success		201		{object}	model.UserResponse
//	@Failure		400		{object}	merror.httpErrorResponse	"data validation error or invalid request body"
//	@Failure		500		{object}	merror.httpErrorResponse
func (uc *UserController) Register(c *fiber.Ctx) error {
	userReq := new(model.RegisterUserRequest)

	err := c.BodyParser(userReq)
	if err != nil {
		uc.log.Warn(merror.InvalidRequestBody, zap.Error(err))
		return merror.New(merror.InvalidRequestBody, err).WithHTTPCode(fiber.StatusBadRequest).FiberError(c)
	}

	err = uc.validator.ValidateStruct(userReq)
	if err != nil {
		uc.log.Warn(merror.ValidationError, zap.Error(err))
		return merror.New(merror.ValidationError, err).WithHTTPCode(fiber.StatusBadRequest).FiberError(c)
	}

	ctx := c.Context()

	user, err := uc.userUseCase.Register(ctx, userReq)
	if err != nil {
		return merror.FromError(err).FiberError(c)
	}

	sessionId, err := uc.sessionUseCase.CreateSession(ctx, user.UserId.String())
	if err != nil {
		return merror.FromError(err).FiberError(c)
	}

	c.Cookie(createSessionCookie(sessionId, uc.cfg))

	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login
//
//	@Summary		User login
//	@Description	Log in an existing user account by sending a JSON payload.
//	@Tags			authentication
//	@Router			/users/login [post]
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		model.LoginUserRequest	true	"Register new user account"
//	@Success		201		{object}	model.UserResponse
//	@Failure		400		{object}	merror.httpErrorResponse	"data validation error or invalid request body"
//	@Failure		500		{object}	merror.httpErrorResponse
func (uc *UserController) Login(c *fiber.Ctx) error {
	userReq := new(model.LoginUserRequest)

	err := c.BodyParser(userReq)
	if err != nil {
		uc.log.Warn(merror.InvalidRequestBody, zap.Error(err))
		return merror.New(merror.InvalidRequestBody, err).WithHTTPCode(fiber.StatusBadRequest).FiberError(c)
	}

	err = uc.validator.ValidateStruct(userReq)
	if err != nil {
		uc.log.Warn(merror.ValidationError, zap.Error(err))
		return merror.New(merror.ValidationError, err).WithHTTPCode(fiber.StatusBadRequest).FiberError(c)
	}

	ctx := c.Context()

	user, err := uc.userUseCase.Login(ctx, userReq)
	if err != nil {
		return merror.FromError(err).FiberError(c)
	}

	sessionId, err := uc.sessionUseCase.CreateSession(ctx, user.UserId.String())
	if err != nil {
		return merror.FromError(err).FiberError(c)
	}

	c.Cookie(createSessionCookie(sessionId, uc.cfg))

	return c.Status(fiber.StatusOK).JSON(user)
}

// Logout
//
//	@Summary		User logout
//	@Description	Log out current authenticated user.
//	@Tags			authentication
//	@Router			/users/logout [post]
//	@Success		204
//	@Failure		401	{object}	merror.httpErrorResponse	"session id cookie is empty"
//	@Failure		500	{object}	merror.httpErrorResponse
func (uc *UserController) Logout(c *fiber.Ctx) error {
	sessionId := c.Cookies(uc.cfg.Consts.SessionCookieName)

	if sessionId == "" {
		return merror.New(merror.SessionCookieNotFound).WithHTTPCode(fiber.StatusUnauthorized).FiberError(c)
	}

	err := uc.sessionUseCase.DeleteById(c.Context(), sessionId)
	if err != nil {
		return merror.New(merror.SessionServiceError, err).WithHTTPCode(fiber.StatusInternalServerError).FiberError(c)
	}

	c.Cookie(&fiber.Cookie{
		Name:     uc.cfg.Consts.SessionCookieName,
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	return c.SendStatus(fiber.StatusNoContent)
}

// Delete
//
//	@Summary		Delete user
//	@Description	Delete user account by userId provided in url path.
//	@Tags			administration
//	@Router			/users/{userId} [delete]
//	@Param			userId	path	string	true	"User ID"
//	@Success		204
//	@Failure		400	{object}	merror.httpErrorResponse	"userId parameter is empty or userId is not a valid ID"
//	@Failure		500	{object}	merror.httpErrorResponse
func (uc *UserController) Delete(c *fiber.Ctx) error {
	userId := c.Params(uc.cfg.Consts.UserIdParamKey)

	if userId == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	deleteRequest := &model.DeleteUserRequest{UserId: userId}

	err := uc.validator.ValidateStruct(deleteRequest)
	if err != nil {
		uc.log.Warn(merror.ValidationError, zap.Error(err))
		return merror.New(merror.ValidationError, err).WithHTTPCode(fiber.StatusBadRequest).FiberError(c)
	}

	ctx := c.Context()

	err = uc.userUseCase.Delete(ctx, deleteRequest)
	if err != nil {
		return merror.FromError(err).FiberError(c)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Me
//
//	@Summary		Current user
//	@Description	Return current authenticated user information in JSON format.
//	@Tags			user
//	@Router			/users/me [get]
//	@Produce		json
//	@Success		200	{object}	model.UserResponse
//	@Failure		500	{object}	merror.httpErrorResponse
func (uc *UserController) Me(c *fiber.Ctx) error {
	authUser, ok := c.Locals(uc.cfg.Consts.UserContextKey).(*model.AuthUser)
	if !ok {
		return merror.New(merror.TypeAssertionError).WithHTTPCode(fiber.StatusInternalServerError).FiberError(c)
	}

	user, err := uc.userUseCase.GetById(c.Context(), authUser.UserId)
	if err != nil {
		return merror.FromError(err).FiberError(c)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func createSessionCookie(sessionId string, cfg *config.Config) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     cfg.Consts.SessionCookieName,
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(cfg.Consts.SessionExpirationTime),
		SameSite: "strict",
		HTTPOnly: true,
		Secure:   config.APP_ENV == "production",
	}
}
