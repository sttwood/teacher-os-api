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

type PlanHandler struct {
	service *planService.PlanService
}

func NewPlanHandler(service *planService.PlanService) *PlanHandler {
	return &PlanHandler{service: service}
}

func (h *PlanHandler) CreatePlan(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	var req planDto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid request",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.service.CreatePlan(currentUser, req)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.Success(c, http.StatusCreated, result)
}

func (h *PlanHandler) ListPlans(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	var query planDto.ListPlansQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errs.WriteError(c, errs.New(
			"VALIDATION_ERROR",
			"invalid query",
			http.StatusBadRequest,
		))
		return
	}

	items, meta, err := h.service.ListPlans(currentUser, query)
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.SuccessWithMeta(c, http.StatusOK, items, meta)
}

func (h *PlanHandler) GetPlanByID(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	result, err := h.service.GetPlanByID(currentUser, c.Param("id"))
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.Success(c, http.StatusOK, result)
}

func (h *PlanHandler) getCurrentUser(c *gin.Context) (*authDto.UserResponse, bool) {
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
