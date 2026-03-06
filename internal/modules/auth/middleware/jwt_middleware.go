package middleware

import (
	"strings"

	"teacher-os-api/internal/modules/auth/service"
	"teacher-os-api/internal/shared/errs"

	"github.com/gin-gonic/gin"
)

const CurrentUserKey = "currentUser"

type JWTMiddleware struct {
	authService *service.AuthService
}

func NewJWTMiddleware(authService *service.AuthService) *JWTMiddleware {
	return &JWTMiddleware{
		authService: authService,
	}
}

func (m *JWTMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errs.WriteError(c, errs.ErrMissingAuthorizationHeader)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
			errs.WriteError(c, errs.ErrInvalidAuthorizationHeader)
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		userID, err := m.authService.ParseAccessToken(tokenString)
		if err != nil {
			errs.WriteError(c, err)
			c.Abort()
			return
		}

		user, err := m.authService.GetUserByID(userID)
		if err != nil {
			errs.WriteError(c, err)
			c.Abort()
			return
		}

		c.Set(CurrentUserKey, user)
		c.Next()
	}
}
