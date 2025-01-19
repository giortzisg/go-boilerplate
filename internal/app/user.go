package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/giortzisg/go-boilerplate/api/v1"
	"github.com/giortzisg/go-boilerplate/internal/model"
	"github.com/giortzisg/go-boilerplate/internal/repository"
	e "github.com/giortzisg/go-boilerplate/pkg/error"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var (
	ErrUserNotFound = e.NewStatusError(errors.New("user not found"), http.StatusNotFound)
	ErrUserExists   = e.NewStatusError(errors.New("user already exists"), http.StatusBadRequest)
)

type UserService interface {
	GetByEmail(ctx context.Context, req *v1.GetUserByEmailRequest) (*v1.GetUserResponse, error)
	Create(ctx context.Context, user *v1.CreateUserRequest) error
	Update(ctx context.Context, user *v1.UpdateUserRequest) error
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository,
	}
}

type userService struct {
	userRepo repository.UserRepository
}

func (u *userService) getUserModelByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, e.NewStatusError(err, http.StatusInternalServerError)
	}

	return user, nil
}

func (u *userService) GetByEmail(ctx context.Context, req *v1.GetUserByEmailRequest) (*v1.GetUserResponse, error) {
	user, err := u.getUserModelByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, e.NewStatusError(err, http.StatusInternalServerError)
	}

	return &v1.GetUserResponse{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userService) Create(ctx context.Context, user *v1.CreateUserRequest) error {
	modelUser, err := u.getUserModelByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return e.NewStatusError(err, http.StatusInternalServerError)
	}

	if modelUser != nil {
		return ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return e.NewStatusError(fmt.Errorf("failed to hash password: %e", err), http.StatusInternalServerError)
	}

	if err = u.userRepo.Create(ctx, &model.User{
		Id:        uuid.New(),
		Name:      user.Name,
		Email:     user.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return e.NewStatusError(err, http.StatusInternalServerError)
	}

	return nil
}

func (u *userService) Update(ctx context.Context, user *v1.UpdateUserRequest) error {
	modelUser, err := u.getUserModelByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrUserNotFound
		}
		return e.NewStatusError(err, http.StatusInternalServerError)
	}

	modelUser.Name = user.Name
	modelUser.UpdatedAt = time.Now()
	if err = u.userRepo.Update(ctx, modelUser); err != nil {
		return e.NewStatusError(err, http.StatusInternalServerError)
	}

	return nil
}
