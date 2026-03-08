package service

import (
	"encoding/json"

	authDto "teacher-os-api/internal/modules/auth/dto"
	exportDto "teacher-os-api/internal/modules/export/dto"
	planModel "teacher-os-api/internal/modules/plans/model"
	planRepository "teacher-os-api/internal/modules/plans/repository"
	"teacher-os-api/internal/shared/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExportService struct {
	planRepo *planRepository.PlanRepository
}

func NewExportService(planRepo *planRepository.PlanRepository) *ExportService {
	return &ExportService{
		planRepo: planRepo,
	}
}

func (s *ExportService) GetLessonPlanPreview(
	currentUser *authDto.UserResponse,
	planID string,
) (*exportDto.LessonPlanDocument, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	parsedPlanID, err := uuid.Parse(planID)
	if err != nil {
		return nil, errs.New("PLAN_INVALID_ID", "invalid plan id", 400)
	}

	plan, err := s.planRepo.FindPlanByIDAndOwner(parsedPlanID, ownerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.New("PLAN_NOT_FOUND", "plan not found", 404)
		}
		return nil, errs.Internal(err)
	}

	widgets, err := s.planRepo.ListWidgetsByPlanID(parsedPlanID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	doc := &exportDto.LessonPlanDocument{
		PlanID:       plan.ID.String(),
		Title:        plan.Title,
		SubjectGroup: plan.SubjectGroup,
		GradeLevel:   plan.GradeLevel,
		Semester:     plan.Semester,
		AcademicYear: plan.AcademicYear,
		SchoolName:   plan.SchoolName,
		TeacherName:  plan.TeacherName,

		BasicInfo:         map[string]interface{}{},
		Objective:         map[string]interface{}{},
		StandardIndicator: map[string]interface{}{},
		Activity:          map[string]interface{}{},
		MediaMaterial:     map[string]interface{}{},
		Assessment:        map[string]interface{}{},
		Homework:          map[string]interface{}{},
	}

	for _, widget := range widgets {
		content := decodeWidgetContent(widget)

		switch widget.WidgetType {
		case planModel.PlanWidgetTypeBasicInfo:
			doc.BasicInfo = content
		case planModel.PlanWidgetTypeObjective:
			doc.Objective = content
		case planModel.PlanWidgetTypeStandardIndicator:
			doc.StandardIndicator = content
		case planModel.PlanWidgetTypeActivity:
			doc.Activity = content
		case planModel.PlanWidgetTypeMediaMaterial:
			doc.MediaMaterial = content
		case planModel.PlanWidgetTypeAssessment:
			doc.Assessment = content
		case planModel.PlanWidgetTypeHomework:
			doc.Homework = content
		}
	}

	return doc, nil
}

func decodeWidgetContent(widget planModel.PlanWidget) map[string]interface{} {
	result := map[string]interface{}{}

	if len(widget.ContentJSON) == 0 {
		return result
	}

	_ = json.Unmarshal(widget.ContentJSON, &result)
	return result
}
