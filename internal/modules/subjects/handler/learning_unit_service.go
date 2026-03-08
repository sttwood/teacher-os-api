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

type LearningUnitHandler struct {
	service *subjectService.LearningUnitService
}

func NewLearningUnitHandler(service *subjectService.LearningUnitService) *LearningUnitHandler {
	return &LearningUnitHandler{service: service}
}

func (h *LearningUnitHandler) CreateLearningUnit(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	var req subjectDto.CreateLearningUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.CreateLearningUnit(currentUser, c.Param("id"), req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.Success(c, http.StatusCreated, result)
}

func (h *LearningUnitHandler) ListLearningUnits(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	result, err := h.service.ListLearningUnits(currentUser, c.Param("id"))
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.Success(c, http.StatusOK, result)
}

func (h *LearningUnitHandler) getCurrentUser(c *gin.Context) (*authDto.UserResponse, bool) {
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
