package handler

import (
	"net/http"

	authDto "teacher-os-api/internal/modules/auth/dto"
	authMiddleware "teacher-os-api/internal/modules/auth/middleware"
	exportService "teacher-os-api/internal/modules/export/service"
	"teacher-os-api/internal/shared/errs"
	"teacher-os-api/internal/shared/httpx"

	"github.com/gin-gonic/gin"
)

type ExportHandler struct {
	service *exportService.ExportService
}

func NewExportHandler(service *exportService.ExportService) *ExportHandler {
	return &ExportHandler{service: service}
}

func (h *ExportHandler) PreviewLessonPlan(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	result, err := h.service.GetLessonPlanPreview(currentUser, c.Param("id"))
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	httpx.Success(c, http.StatusOK, result)
}

func (h *ExportHandler) getCurrentUser(c *gin.Context) (*authDto.UserResponse, bool) {
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

func (h *ExportHandler) ExportLessonPlanDOCX(c *gin.Context) {
	currentUser, ok := h.getCurrentUser(c)
	if !ok {
		return
	}

	content, filename, err := h.service.ExportLessonPlanDOCX(currentUser, c.Param("id"))
	if err != nil {
		errs.WriteError(c, err)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", content)
}
