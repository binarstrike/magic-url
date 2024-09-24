package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/binarstrike/magic-url/internal/entity"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestUserRepository_Register(t *testing.T) {
	t.Parallel()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	userRepo := NewUserRepository(sqlxDB, lojer.NewFromZap(zap.NewNop()))

	user := &entity.User{
		UserId:   uuid.New(),
		Username: "ucup64",
		Email:    "example128@xyz.com",
		Password: "Wkwkwk123",
	}

	// buat kolom dan baris data yang akan di kembalikan ketika menjalankan perintah query ke database
	rows := sqlmock.NewRows([]string{"user_id", "username", "email", "hashed_password"}).
		AddRow(user.UserId, user.Username, user.Email, user.Password)

	// fungsi .ExpectQuery(queryString) berfungsi untuk memastikan perintah query yang di jalankan ke database sesuai dengan yang di harapkan
	// fungsi .WithArgs(...parameter) adalah untuk memastikan parameter dari perintah query ke database sesuai dengan parameter yang di harapkan
	// fungsi .WillReturnRows(...rows) di gunakan untuk mengembalikan data berupa kolom dan baris dari hasil perintah query ke database
	mock.ExpectQuery(createUser).WithArgs(user.UserId, user.Username, user.Email, user.Password).WillReturnRows(rows)

	createdUser, err := userRepo.Register(context.Background(), user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, createdUser.Username, user.Username)
}

func TestUserRepository_Update(t *testing.T) {
	t.Parallel()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	userRepo := NewUserRepository(sqlxDB, lojer.NewFromZap(zap.NewNop()))

	user := &entity.User{
		UserId:   uuid.New(),
		Username: "otong32",
	}

	rows := sqlmock.NewRows([]string{"user_id", "username"}).AddRow(user.UserId, user.Username)
	ctx := context.Background()

	mock.ExpectQuery(updateUser).WithArgs(user.Username, user.UserId).WillReturnRows(rows)

	updatedUser, err := userRepo.Update(ctx, user)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, updatedUser.Username, user.Username)
}

func TestUserRepository_Delete(t *testing.T) {
	t.Parallel()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	userRepo := NewUserRepository(sqlxDB, lojer.NewFromZap(zap.NewNop()))

	userId := uuid.New()
	ctx := context.Background()

	mock.ExpectExec(deleteUser).WithArgs(userId).WillReturnResult(sqlmock.NewResult(1, 1))

	err := userRepo.Delete(ctx, userId.String())
	assert.NoError(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	t.Parallel()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	userRepo := NewUserRepository(sqlxDB, lojer.NewFromZap(zap.NewNop()))

	userEmail := "example128@xyz.com"
	userName := "otong123"

	rows := sqlmock.NewRows([]string{"username", "email"}).AddRow(userName, userEmail)
	ctx := context.Background()

	mock.ExpectQuery(getUserByEmail).WithArgs(userEmail).WillReturnRows(rows)

	user, err := userRepo.FindByEmail(ctx, userEmail)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Username, userName)
}
func TestUserRepository_GetById(t *testing.T) {
	t.Parallel()

	// it should successfully query the user by id
	t.Run("test_1", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		sqlxDB := sqlx.NewDb(db, "sqlmock")

		userRepo := NewUserRepository(sqlxDB, lojer.NewFromZap(zap.NewNop()))

		userMock := &entity.User{
			UserId:   uuid.New(),
			Username: "ujang512",
		}

		rows := sqlmock.NewRows([]string{"username"}).AddRow(userMock.Username)
		mock.ExpectQuery(getUserById).WithArgs(userMock.UserId).WillReturnRows(rows)

		user, err := userRepo.GetById(context.Background(), userMock.UserId.String())
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, user.Username, userMock.Username)
	})

	// it should return an error when the database query fails
	t.Run("test_2", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		sqlxDB := sqlx.NewDb(db, "sqlmock")

		userRepo := NewUserRepository(sqlxDB, lojer.NewFromZap(zap.NewNop()))

		userMock := &entity.User{UserId: uuid.New()}

		mock.ExpectQuery(getUserById).WithArgs(userMock.UserId).WillReturnError(sql.ErrConnDone)

		user, err := userRepo.GetById(context.Background(), userMock.UserId.String())
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrConnDone)
		assert.Nil(t, user)
	})

	// it should return nil if returned row is empty
	t.Run("test_3", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		sqlxDB := sqlx.NewDb(db, "sqlmock")

		userRepo := NewUserRepository(sqlxDB, lojer.NewFromZap(zap.NewNop()))

		userMock := &entity.User{UserId: uuid.New()}

		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery(getUserById).WithArgs(userMock.UserId).WillReturnRows(rows)

		user, err := userRepo.GetById(context.Background(), userMock.UserId.String())
		assert.NoError(t, err)
		assert.Nil(t, user)
	})
}
