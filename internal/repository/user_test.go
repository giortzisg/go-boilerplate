package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giortzisg/go-boilerplate/internal/model"
	"github.com/giortzisg/go-boilerplate/test/mock/any"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupRepository(t *testing.T) (UserRepository, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	repo := NewRepository(slog.New(slog.NewJSONHandler(os.Stdout, nil)), db)
	userRepo := NewUserRepository(repo)

	return userRepo, mock
}

func TestUserRepository_Create(t *testing.T) {
	userRepo, mock := setupRepository(t)
	defer func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}()

	ctx := context.Background()

	testUser := &model.User{
		Id:        uuid.New(),
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WithArgs(testUser.Id, testUser.Name, testUser.Password, testUser.Email, testUser.CreatedAt, testUser.UpdatedAt, nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := userRepo.Create(ctx, testUser)
		assert.NoError(t, err)
	})

	t.Run("creation error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WithArgs(testUser.Id, testUser.Name, testUser.Password, testUser.Email, testUser.CreatedAt, testUser.UpdatedAt, nil).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := userRepo.Create(ctx, testUser)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	userRepo, mock := setupRepository(t)
	defer func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}()

	ctx := context.Background()

	testEmail := "test@example.com"
	expectedUser := &model.User{
		Id:        uuid.New(),
		Name:      "Test User",
		Email:     testEmail,
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
			AddRow(expectedUser.Id, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(testEmail, 1).
			WillReturnRows(rows)

		user, err := userRepo.GetByEmail(ctx, testEmail)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, expectedUser.Name, user.Name)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(testEmail, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := userRepo.GetByEmail(ctx, testEmail)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_Update(t *testing.T) {
	userRepo, mock := setupRepository(t)
	defer func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}()

	ctx := context.Background()

	testUser := &model.User{
		Id:        uuid.New(),
		Name:      "Updated User",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"password"=$2,"email"=$3,"created_at"=$4,"updated_at"=$5,"deleted_at"=$6 WHERE "users"."deleted_at" IS NULL AND "id" = $7`)).
			WithArgs(testUser.Name, testUser.Password, testUser.Email, any.Time{}, any.Time{}, nil, testUser.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := userRepo.Update(ctx, testUser)
		assert.NoError(t, err)
	})

	t.Run("update error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"password"=$2,"email"=$3,"created_at"=$4,"updated_at"=$5,"deleted_at"=$6 WHERE "users"."deleted_at" IS NULL AND "id" = $7`)).
			WithArgs(testUser.Name, testUser.Password, testUser.Email, any.Time{}, any.Time{}, nil, testUser.Id).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := userRepo.Update(ctx, testUser)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	userRepo, mock := setupRepository(t)
	defer func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}()

	ctx := context.Background()

	testID := uuid.New()
	expectedUser := &model.User{
		Id:        testID,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
			AddRow(expectedUser.Id, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(testID, 1).
			WillReturnRows(rows)

		user, err := userRepo.GetByID(ctx, testID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Id, user.Id)
		assert.Equal(t, expectedUser.Name, user.Name)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(testID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := userRepo.GetByID(ctx, testID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, user)
	})
}
