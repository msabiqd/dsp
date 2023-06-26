package usecase

import (
	"context"

	"github.com/SawitProRecruitment/UserService/entity"
)

type UseCaseInterface interface {
	GetUser(ctx context.Context, request entity.GetUserByIdRequest) (output entity.UserResponse, err error)
	Login(ctx context.Context, request entity.GetUserByPhoneNumberRequest, password string) (entity.Token, error)
	CreateUser(ctx context.Context, request entity.CreateUserRequest) (int, error)
	UpdateUser(ctx context.Context, request entity.UpdateUserRequest) error
}
