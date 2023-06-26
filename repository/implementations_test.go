package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetUserById(t *testing.T) {

	tests := []struct {
		name      string
		expectErr bool
		output    entity.User
		request   entity.GetUserByIdRequest
		rows      *sqlmock.Rows
	}{
		{
			name:      "Success",
			expectErr: false,
			output:    entity.User{Id: 1, FullName: "asd", PhoneNumber: "0821", Password: "pw", SuccessfulLogin: 1},
			request:   entity.GetUserByIdRequest{Id: 1},
			rows:      sqlmock.NewRows([]string{"id", "full_name", "phone_number", "password", "successful_login"}).AddRow(1, "asd", "0821", "pw", 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
			defer cancel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(getUserById)).WillReturnRows(tt.rows)

			repo := Repository{Db: db}
			result, err := repo.GetUserById(ctx, tt.request)

			assert.Nil(t, err)
			assert.Equal(t, tt.output, result)
		})
	}
}

func TestGetUserByPhoneNumber(t *testing.T) {

	tests := []struct {
		name      string
		expectErr bool
		output    entity.User
		request   entity.GetUserByPhoneNumberRequest
		rows      *sqlmock.Rows
	}{
		{
			name:      "Success",
			expectErr: false,
			output:    entity.User{Id: 1, FullName: "asd", PhoneNumber: "0821", Password: "pw", SuccessfulLogin: 1},
			request:   entity.GetUserByPhoneNumberRequest{PhoneNumber: "0821"},
			rows:      sqlmock.NewRows([]string{"id", "full_name", "phone_number", "password", "successful_login"}).AddRow(1, "asd", "0821", "pw", 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
			defer cancel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(getUserByPhoneNumber)).WillReturnRows(tt.rows)

			repo := Repository{Db: db}
			result, err := repo.GetUserByPhoneNumber(ctx, tt.request)

			assert.Nil(t, err)
			assert.Equal(t, tt.output, result)
		})
	}

}

func TestCreateUser(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(createUser)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	repo := Repository{Db: db}
	result, err := repo.CreateUser(ctx, entity.CreateUserRequest{FullName: "test", PhoneNumber: "0821"})

	assert.Nil(t, err)
	assert.Equal(t, 1, result)
}

func TestUpdateUser(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(updateUser)).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := Repository{Db: db}
	err = repo.UpdateUser(ctx, entity.UpdateUserRequest{FullName: "test", PhoneNumber: "0821"})

	assert.Nil(t, err)
}

func TestSuccesLoginIncrement(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(succesLoginIncrement)).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := Repository{Db: db}
	err = repo.SuccesLoginIncrement(ctx, entity.SuccessLoginIncrementRequest{Id: 1, SuccessfulLogin: 2})

	assert.Nil(t, err)
}
