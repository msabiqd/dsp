package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetByUser(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
		output    entity.UserResponse
		request   entity.GetUserByIdRequest
	}{
		{
			name:      "Success",
			expectErr: false,
			request:   entity.GetUserByIdRequest{Id: 1},
			output:    entity.UserResponse{FullName: "testuser"},
		},
		{
			name:      "Repository Return Error",
			expectErr: true,
			request:   entity.GetUserByIdRequest{Id: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now())
			defer cancel()
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			mockRepo := repository.NewMockRepositoryInterface(mockCtl)

			switch tt.name {
			case "Success":
				gomock.InOrder(
					mockRepo.EXPECT().GetUserById(ctx, tt.request).Return(entity.User{Id: 1, FullName: "testuser"}, nil),
				)
			case "Repository Return Error":
				gomock.InOrder(
					mockRepo.EXPECT().GetUserById(ctx, tt.request).Return(entity.User{}, errors.New("unit test")),
				)
			}

			uc := NewUseCase(mockRepo)
			result, err := uc.GetUser(ctx, entity.GetUserByIdRequest{Id: 1})
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.output, result)
			}
		})

	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name          string
		expectErr     bool
		output        entity.UserResponse
		request       entity.GetUserByPhoneNumberRequest
		expectedCalls []*gomock.Call
	}{
		{
			name:      "Success",
			expectErr: false,
			request:   entity.GetUserByPhoneNumberRequest{PhoneNumber: "62821111111"},
			output:    entity.UserResponse{FullName: "testuser"},
		},
		{
			name:      "GetUserByPhoneNumber Return Error",
			expectErr: true,
			request:   entity.GetUserByPhoneNumberRequest{PhoneNumber: "62821111111"},
			output:    entity.UserResponse{},
		},
		{
			name:      "SuccesLoginIncrement Return Error",
			expectErr: true,
			request:   entity.GetUserByPhoneNumberRequest{PhoneNumber: "62821111111"},
			output:    entity.UserResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now())
			defer cancel()
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			mockRepo := repository.NewMockRepositoryInterface(mockCtl)

			switch tt.name {
			case "Success":
				gomock.InOrder(
					mockRepo.EXPECT().GetUserByPhoneNumber(ctx, tt.request).Return(entity.User{Id: 1, PhoneNumber: "62821111111", Password: "$2a$10$mlfTmjvJI/vbe.Su8.z0ju6vCNbmm4LpRB8NXJ/yKpkyDy1tSoXrO"}, nil),
					mockRepo.EXPECT().SuccesLoginIncrement(ctx, entity.SuccessLoginIncrementRequest{Id: 1, SuccessfulLogin: 1}).Return(nil),
				)
			case "GetUserByPhoneNumber Return Error":
				gomock.InOrder(
					mockRepo.EXPECT().GetUserByPhoneNumber(ctx, tt.request).Return(entity.User{}, errors.New("unit test")),
				)
			case "SuccesLoginIncrement Return Error":
				gomock.InOrder(
					mockRepo.EXPECT().GetUserByPhoneNumber(ctx, tt.request).Return(entity.User{Id: 1, PhoneNumber: "62821111111", Password: "$2a$10$mlfTmjvJI/vbe.Su8.z0ju6vCNbmm4LpRB8NXJ/yKpkyDy1tSoXrO"}, nil),
					mockRepo.EXPECT().SuccesLoginIncrement(ctx, entity.SuccessLoginIncrementRequest{Id: 1, SuccessfulLogin: 1}).Return(errors.New("unit test")))
			}

			uc := NewUseCase(mockRepo)
			result, err := uc.Login(ctx, entity.GetUserByPhoneNumberRequest{PhoneNumber: "62821111111"}, "Password1!")
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
			}
		})

	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
		output    int
		input     entity.CreateUserRequest
	}{
		{
			name:      "Success",
			expectErr: false,
			output:    1,
			input:     entity.CreateUserRequest{FullName: "testing", PhoneNumber: "628211111111", Password: "Password1!"},
		},
		{
			name:      "Invalid Payload",
			expectErr: true,
			output:    1,
			input:     entity.CreateUserRequest{FullName: "a", PhoneNumber: "08", Password: "asd"},
		},
		{
			name:      "Repository Return Error",
			expectErr: true,
			input:     entity.CreateUserRequest{FullName: "testing", PhoneNumber: "628211111111", Password: "Password1!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now())
			defer cancel()
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			mockRepo := repository.NewMockRepositoryInterface(mockCtl)

			switch tt.name {
			case "Success":
				gomock.InOrder(
					mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(1, nil),
				)
			case "Repository Return Error":
				gomock.InOrder(
					mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(0, errors.New("unit test")),
				)
			}

			uc := NewUseCase(mockRepo)
			result, err := uc.CreateUser(ctx, tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.output, result)
			}
		})

	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
		request   entity.UpdateUserRequest
	}{
		{
			name:      "Success",
			expectErr: false,
			request:   entity.UpdateUserRequest{Id: 1, FullName: "usertesting"},
		},
		{
			name:      "GetUserById Return Error",
			expectErr: true,
			request:   entity.UpdateUserRequest{Id: 1, FullName: "usertesting"},
		},
		{
			name:      "UpdateUser Return Error",
			expectErr: true,
			request:   entity.UpdateUserRequest{Id: 1, FullName: "usertesting"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now())
			defer cancel()
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			mockRepo := repository.NewMockRepositoryInterface(mockCtl)

			switch tt.name {
			case "Success":
				gomock.InOrder(
					mockRepo.EXPECT().GetUserById(ctx, entity.GetUserByIdRequest{Id: 1}).Return(entity.User{Id: 1, FullName: "usertesting"}, nil),
					mockRepo.EXPECT().UpdateUser(ctx, entity.UpdateUserRequest{Id: 1, FullName: "usertesting"}).Return(nil),
				)
			case "GetUserById Return Error":
				gomock.InOrder(
					mockRepo.EXPECT().GetUserById(ctx, entity.GetUserByIdRequest{Id: 1}).Return(entity.User{}, errors.New("unit test")),
				)
			case "UpdateUser Return Error":
				gomock.InOrder(
					mockRepo.EXPECT().GetUserById(ctx, entity.GetUserByIdRequest{Id: 1}).Return(entity.User{Id: 1, FullName: "usertesting"}, nil),
					mockRepo.EXPECT().UpdateUser(ctx, entity.UpdateUserRequest{Id: 1, FullName: "usertesting"}).Return(errors.New("unit test")),
				)
			}

			uc := NewUseCase(mockRepo)
			err := uc.UpdateUser(ctx, entity.UpdateUserRequest{Id: 1})
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})

	}
}
