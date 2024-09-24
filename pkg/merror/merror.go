package merror

import (
	"errors"
	"fmt"

	"github.com/binarstrike/magic-url/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

var errEmpty = errors.New("")

type Merror struct {
	Cause    string
	Detail   error
	httpCode int
}

type httpErrorResponse struct {
	Cause  string `json:"cause"`
	Detail string `json:"detail"`
}

func New(cause string, detail ...error) *Merror {
	merror := new(Merror)
	merror.Detail = errEmpty

	if len(detail) > 0 {
		merror.Detail = detail[0]
	}

	merror.Cause = cause

	return merror
}

func (me *Merror) Error() string {
	return fmt.Sprintf("cause, detail: %s, %v", me.Cause, me.Detail)
}

func (me *Merror) WithHTTPCode(code int) *Merror {
	me.httpCode = code
	return me
}

func (me *Merror) FiberError(c *fiber.Ctx) error {
	return c.Status(utils.If(me.httpCode <= 0, fiber.StatusInternalServerError, me.httpCode)).
		JSON(httpErrorResponse{
			Cause:  me.Cause,
			Detail: me.Detail.Error(),
		})
}

func FromError(err error) *Merror {
	merror, ok := err.(*Merror)
	if ok {
		return merror
	}

	return New(UnknownError, err)
}

const (
	UnknownError           = "unknown error"
	InvalidRequestBody     = "invalid request body"
	ValidationError        = "validation error"
	InternalDatabaseError  = "internal database service error"
	InternalRedisError     = "internal redis service error"
	SessionServiceError    = "session service error"
	AuthServiceError       = "authentication service error"
	EmailExistsError       = "user or email already exists"
	CredentialError        = "wrong or invalid user and password"
	UserNotFound           = "user not found"
	SessionExpires         = "session not found or already expired"
	SessionCookieNotFound  = "session id cookie not found or empty"
	JWTTokenCookieNotFound = "jwt token cookie not found"
	BcryptPassCompareError = "bcrypt password comparation error"
	BcryptPassGenError     = "bcrypt password hash generation error"
	UUIDGenError           = "uuid generation error"
	UUIDParsingError       = "uuid parsing error"
	JWTGenError            = "jwt token generation error"
	JWTParsingError        = "jwt token parsing error"
	PaylodParsingError     = "request payload parsing error"
	TypeAssertionError     = "type assertion error"
)
