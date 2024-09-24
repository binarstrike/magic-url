package entity

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId    uuid.UUID `json:"user_id,omitempty" db:"user_id" redis:"user_id"`
	Username  string    `json:"username,omitempty" db:"username" redis:"username"`
	Email     string    `json:"email,omitempty" db:"email" redis:"email"`
	Password  string    `json:"password,omitempty" db:"hashed_password"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" redis:"updated_at"`
}

type UserWithToken struct {
	*User `json:"user"`
	Token string `json:"jwt_token"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}

	return nil
}

func (u *User) SanitizePassword() {
	u.Password = ""
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for decoding a User object from redis.
func (u *User) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, u)
}
