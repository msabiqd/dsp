package usecase

import (
	"context"
	"errors"
	"net/http"
	"unicode"

	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/pkg/apierror"
	"github.com/SawitProRecruitment/UserService/pkg/jwt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var (
	privateKey = `-----BEGIN PRIVATE KEY-----
	MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDYVfVnBxy0yHAJ
	4q0UiW2McQEc+OUJQyh7UL/pyEMODzvSSg61scz+76PtBIOfe0vC58BtH3li0Ye7
	5X0aAKV12C5XQrtFXZWnPgtqo9YuE9p84vbY4GVpbk+ywta0lSsspRu0F4Xp2aXi
	s0kJMY1Yq0cpniX1yd6RxcN/fTIT3VDFSYTGTiC/ZDxv5ibQXp8DXuQbGgIBd9ls
	8s8vWkmirAjN1Pc8DW9xNDa5Ae45crifK1WAvkCAbbneQUkWVIi/bhQLM2SzwHUK
	Zgm82GDzmMtA3G7bbXZEEnPwA8o24PfjNRxH7E5b8NbL8Fyon6zN2Rl1tb5nVWLp
	g/AX0cMnAgMBAAECggEAMJLkXfS32lOi7GVMDW9p/H2nDVVJP9IndcDExn4jqDV9
	bhMYyG0apSczfFVmJFnvwdx9VUMa4zg+rM6zTzJT9GjMxuUB3WpM3tdMgu40efYV
	ObNQT5Pa0VhmZrHeuX9AyW5tEPuzIrWuzH8K6BiLLxyOBucuiMiBw+NOqQJ4SwMd
	nqHhsDefKIhQR000SfxzceLy2ykiU4Yb6HH0bqGB7sySkp4roA+cgwWdHOWt+2Bc
	qQERXXwyDwrhZYIctf+h+JT/OF2cUm4APyRMBq9DiNlO590vO2TYiiBuUWV5lMop
	rgkAoDMH34EkTBG6NHgRLBcMnFrYjiyO7wyEBQii2QKBgQDuENRkiKAgDnmIpkI+
	yNblHcMzcr4fzpOniO80ly+iSu+8ZkbQaBwKmKMuLp3Ag1vXUb912yy3y/PwG/q6
	YvPWT/oT8uOST6sZhzG7S+IETOtO9oMQgq/yZ5vWPR5Aqvn391AZ26vrPnffbjbU
	hUQsYkprqzOD9pvR4usYR3sWdQKBgQDoog9YvtBWmIj8RWm8uXS9AVCpDZxiTOj4
	uy9philbQmXT7O1CDO6eKKTLVAvJEUMe3WNXsf9sOLr8OLCJlYAvvmb5KiEmimLv
	8WxxvZMbZGWjqEZr2q3BHgQ9LpygPsjyUDrV+C7rSoxnBIvTFPcXVpb//GDqB7In
	RmKCpRhXqwKBgQCRLZY4kNEFe45F9Q3k99mE83d8wnaLMxD8VBd7/M7Bq+0y+TQq
	F7MyCw9INIljQzgYwPN/Td+hXjEutgFa9Mk6Yp9g1vhM22S+NqHvVEFK2hWNm+sl
	gG0H6IMyTSdLzHiP7TPd8QaQeRHlIXMyeuquHmq/6jGKQjOX9UflEIJRmQKBgCZL
	JwOQxXK7wTDlrDYowRKruL9bQjbcOi1XgsJ4Fy6yi1iLU7LfthfK0PW3bAk5Ejdu
	cf/piQYjgIQsQMMlHOi/CuxRBwv1E7wznYpLjn+f0ytRc+YlJRz49/GqTpYCP8nD
	cyvtoquQpfP/R5UDinkJA+o3fSyI/8Z/S4/95TfpAoGATBUlHIOFuV8lOSXav39m
	sUXbfyMsyP2+gDtNUk0QjLnzj7yNeyuDwyQtTWZkWznw0KoSDAFhO2pToj8TXUFz
	H07JY1DafuQiyalyEeTIS6VGngzBw7PMTzTkSLcZ2XPCsi05A+kr/dwWpKylNIhs
	qp4Plq0ceuc/q1whN4V6ieI=
-----END PRIVATE KEY-----`
)

func (u *UseCase) GetUser(ctx context.Context, request entity.GetUserByIdRequest) (entity.UserResponse, error) {
	user, err := u.Repository.GetUserById(ctx, request)
	if err != nil {
		return entity.UserResponse{}, err
	}

	return entity.UserResponse{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}, nil
}

func (u *UseCase) Login(ctx context.Context, request entity.GetUserByPhoneNumberRequest, password string) (entity.Token, error) {
	user, err := u.Repository.GetUserByPhoneNumber(ctx, request)
	if err != nil {
		if user.Id == 0 {
			return entity.Token{}, apierror.New(errors.New("wrong username or password"), http.StatusBadRequest)
		}
		return entity.Token{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return entity.Token{}, apierror.New(errors.New("wrong username or password"), http.StatusBadRequest)
	}

	err = u.Repository.SuccesLoginIncrement(ctx, entity.SuccessLoginIncrementRequest{
		Id:              user.Id,
		SuccessfulLogin: user.SuccessfulLogin + 1})
	if err != nil {
		return entity.Token{}, err
	}

	token, err := jwt.GenerateJWT(user.Id, user.PhoneNumber, privateKey)
	if err != nil {
		return entity.Token{}, err
	}
	return entity.Token{Token: token}, nil
}

func (u *UseCase) CreateUser(ctx context.Context, request entity.CreateUserRequest) (int, error) {
	validate := validator.New()
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("phonenumber", validateIDPhoneNumber)
	err := validate.Struct(request)
	if err != nil {
		return 0, apierror.New(err, http.StatusBadRequest)
	}

	encodedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return 0, err
	}
	request.Password = string(encodedPassword)
	res, err := u.Repository.CreateUser(ctx, request)
	return res, err
}

func (u *UseCase) UpdateUser(ctx context.Context, request entity.UpdateUserRequest) error {
	user, err := u.Repository.GetUserById(ctx, entity.GetUserByIdRequest{Id: request.Id})
	if err != nil {
		return err
	}
	if request.FullName == "" {
		request.FullName = user.FullName
	}
	if request.PhoneNumber == "" {
		request.PhoneNumber = user.PhoneNumber
	}

	err = u.Repository.UpdateUser(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func validatePassword(fl validator.FieldLevel) bool {
	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	value := fl.Field().Interface().(string)
	for _, c := range value {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasNumber && hasUpperCase && hasLowercase && hasSpecial
}

func validateIDPhoneNumber(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)
	countryCode := value[:2]
	return countryCode == "62"
}
