// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/entity"
)

type RepositoryInterface interface {
	GetUserById(ctx context.Context, request entity.GetUserByIdRequest) (entity.User, error)
	GetUserByPhoneNumber(ctx context.Context, request entity.GetUserByPhoneNumberRequest) (entity.User, error)
	CreateUser(ctx context.Context, request entity.CreateUserRequest) (int, error)
	UpdateUser(ctx context.Context, request entity.UpdateUserRequest) error
	SuccesLoginIncrement(ctx context.Context, requeset entity.SuccessLoginIncrementRequest) error
}
