package apierror

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Err            error `json:"message"`
	HttpStatusCode int   `json:"-"`
}

func (e APIError) Error() string {
	return e.Err.Error()
}

func New(err error, httpStatusCode int) error {
	return APIError{
		Err:            err,
		HttpStatusCode: httpStatusCode,
	}
}

func Translate(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required",
			fe.Field())
	case "min":
		return fmt.Sprintf("%s must be contain minimum %s characters",
			fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be contain maximum %s characters",
			fe.Field(), fe.Param())
	case "password":
		return fmt.Sprintf("%s must be contain at least one capital, one number, and one special character",
			fe.Field())
	case "phonenumber":
		return fmt.Sprintf("%s must start with 62",
			fe.Field())
	}
	return ""
}
