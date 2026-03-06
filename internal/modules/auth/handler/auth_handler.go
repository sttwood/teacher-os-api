package handler

import (
	"net/http"

	"teacher-os-api/internal/modules/auth/dto"
	authMiddleware "teacher-os-api/internal/modules/auth/middleware"
	"teacher-os-api/internal/modules/auth/service"
	"teacher-os-api/internal/shared/errs"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.Register(req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.Login(req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.Logout(req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	currentUser, exists := c.Get(authMiddleware.CurrentUserKey)
	if !exists {
		errs.WriteError(c, errs.ErrUnauthorized)
		return
	}

	user, ok := currentUser.(*dto.UserResponse)
	if !ok {
		errs.WriteError(c, errs.ErrInvalidCurrentUserContext)
		return
	}

	c.JSON(http.StatusOK, dto.MeResponse{
		User: *user,
	})
}

func (h *AuthHandler) ResendVerifyEmail(c *gin.Context) {
	var req dto.ResendVerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.ResendVerifyEmail(req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) ConfirmVerifyEmail(c *gin.Context) {
	var req dto.ConfirmVerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.ConfirmVerifyEmail(req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.ForgotPassword(req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.ResetPassword(req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.Refresh(req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
