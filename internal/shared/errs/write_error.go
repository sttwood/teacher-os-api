package errs

import (
	"teacher-os-api/internal/shared/httpx"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	appErr, ok := err.(*AppError)
	if !ok {
		httpx.Error(c, 500, "INTERNAL_ERROR", "internal server error")
		return
	}

	httpx.Error(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
}
