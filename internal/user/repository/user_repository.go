package repository

import (
	"context"
	"database/sql"

	"github.com/binarstrike/magic-url/internal/entity"
	"github.com/binarstrike/magic-url/internal/user"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/binarstrike/magic-url/pkg/merror"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type userRepository struct {
	db  *sqlx.DB
	log lojer.Logger
}

func NewUserRepository(db *sqlx.DB, logger lojer.Logger) user.UserRepository {
	return &userRepository{db: db, log: logger}
}

func (ur userRepository) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := ur.db.QueryRowxContext(ctx, createUser, user.UserId, user.Username, user.Email, user.Password).StructScan(user)
	if err != nil {
		ur.log.Error(merror.InternalDatabaseError, zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (ur userRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := ur.db.QueryRowxContext(ctx, updateUser, user.Username, user.UserId).StructScan(user)
	if err != nil {
		ur.log.Error(merror.InternalDatabaseError, zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (ur userRepository) Delete(ctx context.Context, userId string) error {
	_, err := ur.db.ExecContext(ctx, deleteUser, userId)
	if err != nil {
		ur.log.Error(merror.InternalDatabaseError, zap.Error(err))
		return err
	}

	return nil
}

func (ur userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := new(entity.User)
	err := ur.db.QueryRowxContext(ctx, getUserByEmail, email).StructScan(user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		ur.log.Error(merror.InternalDatabaseError, zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (ur userRepository) GetById(ctx context.Context, userId string) (*entity.User, error) {
	user := new(entity.User)
	err := ur.db.QueryRowxContext(ctx, getUserById, userId).StructScan(user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		ur.log.Error(merror.InternalDatabaseError, zap.Error(err))
		return nil, err
	}

	return user, nil
}
