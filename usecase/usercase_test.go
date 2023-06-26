package usecase

import (
	"context"

	"github.com/SawitProRecruitment/UserService/entity"
)

type MockRepository struct {
	user entity.User
	err  error
}

func (r MockRepository) GetUserById(ctx context.Context, request entity.GetUserByIdRequest) (entity.User, error) {
	return r.user, r.err
}

func (r MockRepository) GetUserByPhoneNumber(ctx context.Context, request entity.GetUserByPhoneNumberRequest) (entity.User, error) {
	return r.user, r.err
}

func (r MockRepository) CreateUser(ctx context.Context, request entity.CreateUserRequest) (int, error) {
	return r.user.Id, r.err
}

func (r MockRepository) UpdateUser(ctx context.Context, request entity.UpdateUserRequest) error {
	return r.err
}

func (r MockRepository) SuccesLoginIncrement(ctx context.Context, requeset entity.SuccessLoginIncrementRequest) error {
	return r.err
}
