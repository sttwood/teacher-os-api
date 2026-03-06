package errs

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

func WriteError(c *gin.Context, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = Internal(err)
	}

	c.JSON(appErr.HTTPStatus, ErrorResponse{
		Error: ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}
