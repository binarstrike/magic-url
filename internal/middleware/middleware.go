package middleware

import (
	"time"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/session"
	"github.com/binarstrike/magic-url/internal/user"
	"github.com/binarstrike/magic-url/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

type Middleware struct {
	userUseCase    user.UserUseCase
	sessionUseCase session.SessionUseCase
	validator      *validator.Validator
	cfg            *config.Config
	CSRFMiddleware fiber.Handler
}

func NewMiddleware(userUseCase user.UserUseCase, sessionUseCase session.SessionUseCase, validator *validator.Validator, cfg *config.Config) *Middleware {
	m := new(Middleware)
	m.CSRFMiddleware = csrf.New(csrf.Config{
		KeyLookup:      "header:" + cfg.Consts.CSRFHeaderName,
		CookieName:     cfg.Consts.CSRFCookieName,
		ContextKey:     cfg.Consts.CSRFContextKey,
		CookieSameSite: "Strict",
		CookieHTTPOnly: false,
		CookieSecure:   config.APP_ENV == "production",
		Expiration:     time.Hour * 1,
	})
	m.userUseCase = userUseCase
	m.sessionUseCase = sessionUseCase
	m.validator = validator
	m.cfg = cfg

	return m
}
