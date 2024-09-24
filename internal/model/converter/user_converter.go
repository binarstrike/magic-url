package converter

import (
	"github.com/binarstrike/magic-url/internal/entity"
	"github.com/binarstrike/magic-url/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		UserId:    user.UserId,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
