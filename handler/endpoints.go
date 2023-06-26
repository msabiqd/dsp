package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/pkg/apierror"
	"github.com/SawitProRecruitment/UserService/pkg/jwt"
	"github.com/SawitProRecruitment/UserService/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	publicKey = `-----BEGIN RSA PUBLIC KEY-----
	MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2FX1ZwcctMhwCeKtFIlt
	jHEBHPjlCUMoe1C/6chDDg870koOtbHM/u+j7QSDn3tLwufAbR95YtGHu+V9GgCl
	ddguV0K7RV2Vpz4LaqPWLhPafOL22OBlaW5PssLWtJUrLKUbtBeF6dml4rNJCTGN
	WKtHKZ4l9cnekcXDf30yE91QxUmExk4gv2Q8b+Ym0F6fA17kGxoCAXfZbPLPL1pJ
	oqwIzdT3PA1vcTQ2uQHuOXK4nytVgL5AgG253kFJFlSIv24UCzNks8B1CmYJvNhg
	85jLQNxu2212RBJz8APKNuD34zUcR+xOW/DWy/BcqJ+szdkZdbW+Z1Vi6YPwF9HD
	JwIDAQAB
-----END RSA PUBLIC KEY-----`
)

func (s *Server) GetUser(ctx echo.Context) error {
	headers := ctx.Request().Header
	userJWT, err := jwt.VerifyToken(strings.Split(headers.Get("Authorization"), "Bearer ")[1], publicKey)
	if err != nil {
		if apierr, ok := err.(apierror.APIError); ok {
			return ctx.JSON(apierr.HttpStatusCode, response.BuildErrorResponse([]error{apierr}))
		} else {
			return ctx.JSON(http.StatusInternalServerError, response.BuildErrorResponse([]error{err}))
		}
	}
	resp, err := s.UseCase.GetUser(ctx.Request().Context(), entity.GetUserByIdRequest{Id: userJWT.Id})
	if err != nil {
		if apierr, ok := err.(apierror.APIError); ok {
			return ctx.JSON(apierr.HttpStatusCode, response.BuildErrorResponse([]error{apierr}))
		} else {
			return ctx.JSON(http.StatusInternalServerError, response.BuildErrorResponse([]error{err}))
		}
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Registrations(ctx echo.Context) error {

	json_map := make(map[string]interface{})
	err := json.NewDecoder(ctx.Request().Body).Decode(&json_map)
	if err != nil {
		return err
	}

	id, err := s.UseCase.CreateUser(ctx.Request().Context(), entity.CreateUserRequest{
		FullName:    json_map["full_name"].(string),
		PhoneNumber: json_map["phone_number"].(string),
		Password:    json_map["password"].(string),
	})
	if err != nil {
		var errs []error
		if apierr, ok := err.(apierror.APIError); ok {
			if ves, ok := apierr.Err.(validator.ValidationErrors); ok {
				for _, ve := range ves {
					errs = append(errs, apierror.New(errors.New(apierror.Translate(ve)), apierr.HttpStatusCode))
				}
			} else {
				errs = append(errs, apierr)
			}
			return ctx.JSON(apierr.HttpStatusCode, response.BuildErrorResponse(errs))
		} else {
			errs = append(errs, err)
			return ctx.JSON(http.StatusInternalServerError, response.BuildErrorResponse(errs))
		}
	}

	return ctx.JSON(http.StatusCreated, id)
}

func (s *Server) UpdateUser(ctx echo.Context) error {

	headers := ctx.Request().Header
	userJWT, err := jwt.VerifyToken(strings.Split(headers.Get("Authorization"), "Bearer ")[1], publicKey)
	if err != nil {
		if apierr, ok := err.(apierror.APIError); ok {
			return ctx.JSON(apierr.HttpStatusCode, response.BuildErrorResponse([]error{apierr}))
		} else {
			return ctx.JSON(http.StatusInternalServerError, response.BuildErrorResponse([]error{err}))
		}
	}

	json_map := make(map[string]interface{})
	err = json.NewDecoder(ctx.Request().Body).Decode(&json_map)
	if err != nil {
		return err
	}

	fullName := ""
	phoneNumber := ""

	if json_map["full_name"] != nil {
		if fullnameJson, ok := json_map["full_name"].(string); ok {
			fullName = fullnameJson
		}
	}

	if json_map["phone_number"] != nil {
		if phoneNumberJson, ok := json_map["phone_number"].(string); ok {
			phoneNumber = phoneNumberJson
		}
	}
	err = s.UseCase.UpdateUser(ctx.Request().Context(), entity.UpdateUserRequest{
		Id:          userJWT.Id,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		var errs []error
		if apierr, ok := err.(apierror.APIError); ok {
			if ves, ok := apierr.Err.(validator.ValidationErrors); ok {
				for _, ve := range ves {
					errs = append(errs, apierror.New(errors.New(apierror.Translate(ve)), apierr.HttpStatusCode))
				}
			} else {
				errs = append(errs, apierr)
			}
			return ctx.JSON(apierr.HttpStatusCode, response.BuildErrorResponse(errs))
		} else {
			errs = append(errs, err)
			return ctx.JSON(http.StatusInternalServerError, response.BuildErrorResponse(errs))
		}
	}

	return ctx.JSON(http.StatusOK, response.BuildSuccessResponse("Success"))
}

func (s *Server) Login(ctx echo.Context) error {
	json_map := make(map[string]interface{})
	err := json.NewDecoder(ctx.Request().Body).Decode(&json_map)
	if err != nil {
		return err
	}

	jwt, err := s.UseCase.Login(ctx.Request().Context(), entity.GetUserByPhoneNumberRequest{
		PhoneNumber: json_map["phone_number"].(string),
	}, json_map["password"].(string))
	if err != nil {
		if apierr, ok := err.(apierror.APIError); ok {
			return ctx.JSON(apierr.HttpStatusCode, response.BuildErrorResponse([]error{apierr}))
		} else {
			return ctx.JSON(http.StatusInternalServerError, response.BuildErrorResponse([]error{err}))
		}
	}

	return ctx.JSON(http.StatusOK, response.BuildSuccessResponse(jwt))
}
