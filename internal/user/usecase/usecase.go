package usecase

import (
	"context"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/internal/entity"
	"github.com/binarstrike/magic-url/internal/model"
	"github.com/binarstrike/magic-url/internal/model/converter"
	"github.com/binarstrike/magic-url/internal/user"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/binarstrike/magic-url/pkg/merror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const prefixKey = "auth-user_"

type userUseCase struct {
	userRepo  user.UserRepository
	redisRepo user.RedisRepository
	log       lojer.Logger
	cfg       *config.Config
}

func NewUserUseCase(userRepo user.UserRepository, redisRepo user.RedisRepository, logger lojer.Logger, config *config.Config) user.UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		redisRepo: redisRepo,
		cfg:       config,
		log:       logger,
	}
}

func (uc *userUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, merror.New(merror.InternalDatabaseError).WithHTTPCode(fiber.StatusInternalServerError)
	} else if user != nil {
		return nil, merror.New(merror.EmailExistsError).WithHTTPCode(fiber.StatusConflict)
	}

	newUser := new(entity.User)

	newUser.UserId, err = uuid.NewV7()
	if err != nil {
		uc.log.Error(merror.UUIDGenError, zap.Error(err))
		return nil, merror.New(merror.UUIDGenError, err).WithHTTPCode(fiber.StatusInternalServerError)
	}

	newUser.Password = request.Password
	err = newUser.HashPassword()
	if err != nil {
		uc.log.Error(merror.BcryptPassGenError, zap.Error(err))
		return nil, merror.New(merror.BcryptPassGenError, err).WithHTTPCode(fiber.StatusInternalServerError)
	}

	newUser.Username = request.Username
	newUser.Email = request.Email

	_, err = uc.userRepo.Register(ctx, newUser)
	if err != nil {
		return nil, merror.New(merror.InternalDatabaseError).WithHTTPCode(fiber.StatusInternalServerError)
	}

	return converter.UserToResponse(newUser), nil
}

func (uc *userUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, merror.New(merror.InternalDatabaseError).WithHTTPCode(fiber.StatusInternalServerError)
	} else if user == nil {
		return nil, merror.New(merror.UserNotFound).WithHTTPCode(fiber.StatusNotFound)
	}

	if err = user.ComparePassword(request.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, merror.New(merror.CredentialError, err).WithHTTPCode(fiber.StatusUnauthorized)
		}
		uc.log.Error(merror.BcryptPassCompareError, zap.Error(err))
		return nil, merror.New(merror.BcryptPassCompareError, err).WithHTTPCode(fiber.StatusInternalServerError)
	}

	return converter.UserToResponse(user), nil
}

func (uc *userUseCase) Delete(ctx context.Context, request *model.DeleteUserRequest) error {
	user, err := uc.userRepo.GetById(ctx, request.UserId)
	if err != nil {
		return merror.New(merror.InternalDatabaseError).WithHTTPCode(fiber.StatusInternalServerError)
	} else if user == nil {
		return merror.New(merror.UserNotFound).WithHTTPCode(fiber.StatusNotFound)
	}

	err = uc.userRepo.Delete(ctx, request.UserId)
	if err != nil {
		return merror.New(merror.InternalDatabaseError).WithHTTPCode(fiber.StatusInternalServerError)
	}

	_ = uc.redisRepo.DeleteUser(ctx, createCacheKey(request.UserId))

	return nil
}

func (uc *userUseCase) GetById(ctx context.Context, userId string) (*model.UserResponse, error) {
	key := createCacheKey(userId)
	cachedUser, _ := uc.redisRepo.GetUser(ctx, key)
	if cachedUser != nil {
		return converter.UserToResponse(cachedUser), nil
	}

	user, err := uc.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, merror.New(merror.InternalDatabaseError).WithHTTPCode(fiber.StatusInternalServerError)
	} else if user == nil {
		return nil, merror.New(merror.UserNotFound).WithHTTPCode(fiber.StatusNotFound)
	}

	_ = uc.redisRepo.SetUser(ctx, key, user)

	return converter.UserToResponse(user), nil
}

func createCacheKey(userId string) string {
	return prefixKey + userId
}

// func (s userUseCase) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
// 	updatedUser, err := s.userRepo.Update(ctx, user)
// 	if err != nil {
// 		return nil, merror.New(merror.DatabaseRepositoryError, err).WithHTTPCode(fiber.StatusInternalServerError)
// 	}

// 	updatedUser.SanitizePassword()

// 	err = s.redisRepo.DeleteUser(ctx, createCacheKey(user.UserId.String()))
// 	if err != nil {
// 		s.logger.Error(merror.RedisRepositoryError, err)
// 	}

// 	return updatedUser, nil
// }
