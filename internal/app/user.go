package app

import (
	"context"
	v1 "github.com/giortzisg/go-boilerplate/api/v1"
)

type UserService interface {
	GetByEmail(ctx context.Context, req *v1.GetUserByEmailRequest) (*v1.GetUserResponse, error)
	Create(ctx context.Context, user *v1.CreateUserRequest) error
	Update(ctx context.Context, user *v1.UpdateUserRequest) error
}
