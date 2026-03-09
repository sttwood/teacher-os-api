package handler

import (
	"net/http"

	authDto "teacher-os-api/internal/modules/auth/dto"
	authMiddleware "teacher-os-api/internal/modules/auth/middleware"
	subjectDto "teacher-os-api/internal/modules/subjects/dto"
	subjectService "teacher-os-api/internal/modules/subjects/service"
	"teacher-os-api/internal/shared/errs"
	"teacher-os-api/internal/shared/httpx"

	"github.com/gin-gonic/gin"
)

type SubjectHandler struct {
	service *subjectService.SubjectService
}

func NewSubjectHandler(service *subjectService.SubjectService) *SubjectHandler {
	return &SubjectHandler{service: service}
}

func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	var req subjectDto.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.CreateSubject(currentUser, req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.Success(c, http.StatusCreated, result)
}

func (h *SubjectHandler) ListSubjects(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	result, err := h.service.ListSubjects(currentUser)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.Success(c, http.StatusOK, result)
}

func (h *SubjectHandler) getCurrentUser(c *gin.Context) (*authDto.UserResponse, bool) {
	currentUser, exists := c.Get(authMiddleware.CurrentUserKey)
	if !exists {
		errs.WriteError(c, errs.ErrUnauthorized)
		return nil, false
	}

	user, ok := currentUser.(*authDto.UserResponse)
	if !ok {
		errs.WriteError(c, errs.ErrInvalidCurrentUserContext)
		return nil, false
	}

	return user, true
}
