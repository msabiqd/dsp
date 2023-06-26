package response

type SuccessResponse struct {
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Errors []ErrorInfo `json:"errors,omitempty"`
}

type ErrorInfo struct {
	Message string `json:"message"`
}

func BuildSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Data: data,
	}
}

func BuildErrorResponse(errors []error) ErrorResponse {
	if len(errors) == 0 {
		return ErrorResponse{
			[]ErrorInfo{{Message: "Internal Server Error"}},
		}
	} else {
		errInfos := []ErrorInfo{}
		for _, err := range errors {
			errInfos = append(errInfos,
				ErrorInfo{
					Message: err.Error(),
				})
		}

		return ErrorResponse{
			Errors: errInfos,
		}
	}
}
