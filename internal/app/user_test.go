package app

import (
	"context"
	"reflect"
	"testing"
	"time"

	"errors"

	v1 "github.com/giortzisg/go-boilerplate/api/v1"
	"github.com/giortzisg/go-boilerplate/internal/model"
	mock_repository "github.com/giortzisg/go-boilerplate/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func Test_userService_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *v1.CreateUserRequest
	}
	type mockExpect struct {
		getByEmailReturn interface{}
		createReturn     error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    mockExpect
	}{
		{
			name: "Create user successfully",
			args: args{
				ctx: context.Background(),
				user: &v1.CreateUserRequest{
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			wantErr: false,
			mock: mockExpect{
				getByEmailReturn: nil,
				createReturn:     nil,
			},
		},
		{
			name: "Create user with existing email",
			args: args{
				ctx: context.Background(),
				user: &v1.CreateUserRequest{
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			wantErr: true,
			mock: mockExpect{
				getByEmailReturn: &model.User{
					Email: "test@example.com",
				},
				createReturn: nil,
			},
		},
		{
			name: "Create user with database error",
			args: args{
				ctx: context.Background(),
				user: &v1.CreateUserRequest{
					Name:     "Test User",
					Email:    "test2@example.com",
					Password: "password123",
				},
			},
			wantErr: true,
			mock: mockExpect{
				getByEmailReturn: nil,
				createReturn:     errors.New("database error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_repository.NewMockUserRepository(gomock.NewController(t))
			mockRepo.EXPECT().GetByEmail(context.Background(), tt.args.user.Email).Return(tt.mock.getByEmailReturn, nil)
			if tt.mock.getByEmailReturn == nil {
				mockRepo.EXPECT().Create(context.Background(), gomock.AssignableToTypeOf(&model.User{
					Id:        uuid.UUID{},
					Name:      tt.args.user.Name,
					Email:     tt.args.user.Email,
					Password:  tt.args.user.Password,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				})).Return(tt.mock.createReturn)
			}
			u := NewUserService(mockRepo)
			if err := u.Create(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_GetByEmail(t *testing.T) {
	type args struct {
		ctx context.Context
		req *v1.GetUserByEmailRequest
	}
	type mockExpect struct {
		getByEmailReturn interface{}
		getByEmailError  error
	}

	tests := []struct {
		name    string
		args    args
		want    *v1.GetUserResponse
		wantErr bool
		mock    mockExpect
	}{
		{
			name: "Get user successfully",
			args: args{
				ctx: context.Background(),
				req: &v1.GetUserByEmailRequest{
					Email: "test@example.com",
				},
			},
			want: &v1.GetUserResponse{
				Name:  "Test User",
				Email: "test@example.com",
			},
			wantErr: false,
			mock: mockExpect{
				getByEmailReturn: &model.User{
					Name:  "Test User",
					Email: "test@example.com",
				},
				getByEmailError: nil,
			},
		},
		{
			name: "User not found",
			args: args{
				ctx: context.Background(),
				req: &v1.GetUserByEmailRequest{
					Email: "notfound@example.com",
				},
			},
			want:    nil,
			wantErr: true,
			mock: mockExpect{
				getByEmailReturn: nil,
				getByEmailError:  ErrUserNotFound,
			},
		},
		{
			name: "Get user with database error",
			args: args{
				ctx: context.Background(),
				req: &v1.GetUserByEmailRequest{
					Email: "error@example.com",
				},
			},
			want:    nil,
			wantErr: true,
			mock: mockExpect{
				getByEmailReturn: nil,
				getByEmailError:  errors.New("database error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_repository.NewMockUserRepository(gomock.NewController(t))
			mockRepo.EXPECT().GetByEmail(tt.args.ctx, tt.args.req.Email).Return(tt.mock.getByEmailReturn, tt.mock.getByEmailError)
			u := NewUserService(mockRepo)
			got, err := u.GetByEmail(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *v1.UpdateUserRequest
	}
	type mockExpect struct {
		getByEmailReturn interface{}
		getByEmailError  error
		updateReturn     error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    mockExpect
	}{
		{
			name: "Update user successfully",
			args: args{
				ctx: context.Background(),
				user: &v1.UpdateUserRequest{
					Email: "test@example.com",
					Name:  "Updated User",
				},
			},
			wantErr: false,
			mock: mockExpect{
				getByEmailReturn: &model.User{
					Email: "test@example.com",
					Name:  "Test User",
				},
				getByEmailError: nil,
				updateReturn:    nil,
			},
		},
		{
			name: "User not found during update",
			args: args{
				ctx: context.Background(),
				user: &v1.UpdateUserRequest{
					Email: "notfound@example.com",
					Name:  "Updated User",
				},
			},
			wantErr: true,
			mock: mockExpect{
				getByEmailReturn: nil,
				getByEmailError:  ErrUserNotFound,
			},
		},
		{
			name: "Update user with database error",
			args: args{
				ctx: context.Background(),
				user: &v1.UpdateUserRequest{
					Email: "test@example.com",
					Name:  "Updated User",
				},
			},
			wantErr: true,
			mock: mockExpect{
				getByEmailReturn: &model.User{
					Email: "test@example.com",
				},
				getByEmailError: nil,
				updateReturn:    errors.New("database error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_repository.NewMockUserRepository(gomock.NewController(t))
			mockRepo.EXPECT().GetByEmail(tt.args.ctx, tt.args.user.Email).Return(tt.mock.getByEmailReturn, tt.mock.getByEmailError)
			if tt.mock.getByEmailReturn != nil {
				mockRepo.EXPECT().Update(tt.args.ctx, gomock.AssignableToTypeOf(&model.User{})).Return(tt.mock.updateReturn)
			}
			u := NewUserService(mockRepo)
			if err := u.Update(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
