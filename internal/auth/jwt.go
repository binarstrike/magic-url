package auth

import (
	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserId string
	jwt.RegisteredClaims
}

func GenerateJWTToken(cfg *config.Config, user *entity.User) (string, error) {
	claims := JWTClaims{UserId: user.UserId.String()}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
