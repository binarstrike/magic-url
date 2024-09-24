package model

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	UserId    uuid.UUID `json:"-"`
	Username  string    `json:"name,omitempty"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,gte=5,lte=16"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

type DeleteUserRequest struct {
	UserId string `json:"user_id" validate:"uuid"`
}

type GetUserRequest struct {
	UserId string `json:"user_id" validate:"uuid"`
	Email  string `json:"email" validate:"email"`
}

type VerifyUserRequest struct {
	SessionId string `json:"session_id" validate:"required,uuid"`
}

type UpdateUserRequest struct {
	UserName string `json:"username" validate:"required,gte=5,lte=16"`
}
