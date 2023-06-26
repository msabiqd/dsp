package repository

import (
	"context"
	"errors"
	"net/http"

	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/pkg/apierror"
	"github.com/lib/pq"
)

func (r *Repository) GetUserById(ctx context.Context, request entity.GetUserByIdRequest) (output entity.User, err error) {
	err = r.Db.QueryRowContext(ctx, getUserById, request.Id).Scan(
		&output.Id,
		&output.FullName,
		&output.PhoneNumber,
		&output.Password,
		&output.SuccessfulLogin)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, request entity.GetUserByPhoneNumberRequest) (output entity.User, err error) {
	err = r.Db.QueryRowContext(ctx, getUserByPhoneNumber, request.PhoneNumber).Scan(
		&output.Id,
		&output.FullName,
		&output.PhoneNumber,
		&output.Password,
		&output.SuccessfulLogin)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateUser(ctx context.Context, request entity.CreateUserRequest) (res int, err error) {
	err = r.Db.QueryRowContext(ctx, createUser, request.FullName, request.PhoneNumber, request.Password).Scan(&res)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code == "23505" {
				return 0, apierror.New(errors.New("Phone number was used by other user"), http.StatusConflict)
			}
		} else {
			return 0, err
		}
	}
	return res, nil
}

func (r *Repository) UpdateUser(ctx context.Context, request entity.UpdateUserRequest) (err error) {
	_, err = r.Db.ExecContext(ctx, updateUser, request.FullName, request.PhoneNumber, request.Id)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code == "23505" {
				return apierror.New(errors.New("Phone number was used by other user"), http.StatusConflict)
			}
		} else {
			return err
		}
	}
	return nil
}

func (r *Repository) SuccesLoginIncrement(ctx context.Context, request entity.SuccessLoginIncrementRequest) (err error) {
	_, err = r.Db.ExecContext(ctx, succesLoginIncrement, request.Id, request.SuccessfulLogin)
	if err != nil {
		return err
	}
	return nil
}
