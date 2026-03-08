package service

import (
	"encoding/json"
	"time"

	authDto "teacher-os-api/internal/modules/auth/dto"
	planDto "teacher-os-api/internal/modules/plans/dto"
	planModel "teacher-os-api/internal/modules/plans/model"
	planRepository "teacher-os-api/internal/modules/plans/repository"
	"teacher-os-api/internal/shared/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanWidgetService struct {
	repo *planRepository.PlanRepository
}

func NewPlanWidgetService(repo *planRepository.PlanRepository) *PlanWidgetService {
	return &PlanWidgetService{repo: repo}
}

func (s *PlanWidgetService) GetWidgets(
	currentUser *authDto.UserResponse,
	planID string,
) ([]planDto.PlanWidgetResponse, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	parsedPlanID, err := uuid.Parse(planID)
	if err != nil {
		return nil, errs.New("PLAN_INVALID_ID", "invalid plan id", 400)
	}

	_, err = s.repo.FindPlanByIDAndOwner(parsedPlanID, ownerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.New("PLAN_NOT_FOUND", "plan not found", 404)
		}
		return nil, errs.Internal(err)
	}

	widgets, err := s.repo.ListWidgetsByPlanID(parsedPlanID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	result := make([]planDto.PlanWidgetResponse, 0, len(widgets))
	for _, widget := range widgets {
		var content map[string]interface{}
		if len(widget.ContentJSON) > 0 {
			_ = json.Unmarshal(widget.ContentJSON, &content)
		}
		if content == nil {
			content = map[string]interface{}{}
		}

		result = append(result, planDto.PlanWidgetResponse{
			ID:          widget.ID.String(),
			PlanID:      widget.PlanID.String(),
			WidgetType:  string(widget.WidgetType),
			OrderIndex:  widget.OrderIndex,
			Title:       widget.Title,
			ContentJSON: content,
			IsCollapsed: widget.IsCollapsed,
			CreatedAt:   widget.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   widget.UpdatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

func (s *PlanWidgetService) SaveWidgets(
	currentUser *authDto.UserResponse,
	planID string,
	input planDto.UpdatePlanWidgetsRequest,
) ([]planDto.PlanWidgetResponse, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	parsedPlanID, err := uuid.Parse(planID)
	if err != nil {
		return nil, errs.New("PLAN_INVALID_ID", "invalid plan id", 400)
	}

	_, err = s.repo.FindPlanByIDAndOwner(parsedPlanID, ownerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.New("PLAN_NOT_FOUND", "plan not found", 404)
		}
		return nil, errs.Internal(err)
	}

	models := make([]planModel.PlanWidget, 0, len(input.Widgets))
	for _, item := range input.Widgets {
		contentBytes, err := json.Marshal(item.ContentJSON)
		if err != nil {
			return nil, errs.New("PLAN_WIDGET_INVALID_CONTENT", "invalid widget content", 400)
		}

		widgetID := uuid.New()
		if item.ID != "" {
			if parsedID, parseErr := uuid.Parse(item.ID); parseErr == nil {
				widgetID = parsedID
			}
		}

		models = append(models, planModel.PlanWidget{
			ID:          widgetID,
			PlanID:      parsedPlanID,
			WidgetType:  planModel.PlanWidgetType(item.WidgetType),
			OrderIndex:  item.OrderIndex,
			Title:       item.Title,
			ContentJSON: contentBytes,
			IsCollapsed: item.IsCollapsed,
		})
	}

	if err := s.repo.ReplaceWidgets(parsedPlanID, models); err != nil {
		return nil, errs.Internal(err)
	}

	return s.GetWidgets(currentUser, planID)
}
