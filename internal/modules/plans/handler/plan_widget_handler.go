package handler

import (
	"net/http"

	authDto "teacher-os-api/internal/modules/auth/dto"
	authMiddleware "teacher-os-api/internal/modules/auth/middleware"
	planDto "teacher-os-api/internal/modules/plans/dto"
	planService "teacher-os-api/internal/modules/plans/service"
	"teacher-os-api/internal/shared/errs"
	"teacher-os-api/internal/shared/httpx"

	"github.com/gin-gonic/gin"
)

type PlanWidgetHandler struct {
	service *planService.PlanWidgetService
}

func NewPlanWidgetHandler(service *planService.PlanWidgetService) *PlanWidgetHandler {
	return &PlanWidgetHandler{service: service}
}

func (h *PlanWidgetHandler) GetWidgets(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	result, err := h.service.GetWidgets(currentUser, c.Param("id"))
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.Success(c, http.StatusOK, result)
}

func (h *PlanWidgetHandler) SaveWidgets(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	var req planDto.UpdatePlanWidgetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New("VALIDATION_ERROR", "invalid request", http.StatusBadRequest))
		return
	}

	result, err := h.service.SaveWidgets(currentUser, c.Param("id"), req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.Success(c, http.StatusOK, result)
}

func (h *PlanWidgetHandler) getCurrentUser(c *gin.Context) (*authDto.UserResponse, bool) {
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
