package middleware

import (
	"fmt"
	"strings"

	"github.com/binarstrike/magic-url/internal/auth"
	"github.com/binarstrike/magic-url/internal/model"
	"github.com/binarstrike/magic-url/pkg/merror"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (m *Middleware) AuthSessionMiddleware(c *fiber.Ctx) error {
	sessionId := c.Cookies(m.cfg.Consts.SessionCookieName)

	if sessionId == "" {
		return merror.New(merror.SessionCookieNotFound).WithHTTPCode(fiber.StatusUnauthorized).FiberError(c)
	}

	err := m.validator.ValidateStruct(model.VerifyUserRequest{SessionId: sessionId})
	if err != nil {
		return merror.New(merror.ValidationError, err).WithHTTPCode(fiber.StatusBadRequest).FiberError(c)
	}

	session, err := m.sessionUseCase.GetSessionById(c.Context(), sessionId)
	if err != nil {
		// ada dua kemungkinan kesalahan yaitu internal server error dan unauthorized jika session tidak ada pada cache/redis
		return merror.FromError(err).FiberError(c)
	}

	// ? pertimbangkan apakah masih memerlukan field Authenticated pada struct SessionData
	if !session.Authenticated {
		return merror.New(merror.SessionExpires).WithHTTPCode(fiber.StatusUnauthorized).FiberError(c)
	}

	c.Locals(m.cfg.Consts.UserContextKey, &model.AuthUser{UserId: session.UserId})

	return c.Next()
}

func (m *Middleware) AuthJWTMiddleware(c *fiber.Ctx) error {
	tokenCookie := c.Cookies(m.cfg.Consts.JWTCookieName)

	if tokenCookie == "" {
		authHeader := c.Get(m.cfg.Consts.JWTHeaderName)

		h := strings.Split(authHeader, " ") // [Bearer eyJhbGcpXVCJ9.eyJzdWIiOiIxMjM0NTYwIiwibmFtZiwiaWFM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMe]
		if len(h) != 2 {
			return merror.New("no token provided").WithHTTPCode(fiber.StatusUnauthorized).FiberError(c)
		}

		tokenCookie = h[1]
	}

	claims := new(auth.JWTClaims)

	token, err := jwt.ParseWithClaims(tokenCookie, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknown signing method: %v", t.Header["alg"])
		}

		return []byte(m.cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return merror.New(merror.JWTParsingError, err).WithHTTPCode(fiber.StatusBadRequest).FiberError(c)
	}

	userId, err := uuid.Parse(claims.UserId)
	if err != nil {
		return merror.New(merror.UUIDParsingError, err).WithHTTPCode(fiber.StatusBadRequest).FiberError(c)
	}

	requestCtx := c.Context()

	user, err := m.userUseCase.GetById(requestCtx, userId.String())
	if err != nil {
		return merror.FromError(err).FiberError(c)
	}

	c.Locals(m.cfg.Consts.UserContextKey, &model.AuthUser{UserId: user.UserId.String()})

	return c.Next()
}
