package service

import (
	"strings"
	"time"

	authDto "teacher-os-api/internal/modules/auth/dto"
	planDto "teacher-os-api/internal/modules/plans/dto"
	planModel "teacher-os-api/internal/modules/plans/model"
	planRepository "teacher-os-api/internal/modules/plans/repository"
	"teacher-os-api/internal/shared/errs"

	"github.com/google/uuid"
)

type PlanService struct {
	repo *planRepository.PlanRepository
}

func NewPlanService(repo *planRepository.PlanRepository) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) CreatePlan(
	currentUser *authDto.UserResponse,
	input planDto.CreatePlanRequest,
) (*planDto.PlanResponse, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	plan := &planModel.Plan{
		ID:           uuid.New(),
		OwnerID:      ownerID,
		Title:        strings.TrimSpace(input.Title),
		SubjectGroup: strings.TrimSpace(input.SubjectGroup),
		GradeLevel:   strings.TrimSpace(input.GradeLevel),
		Semester:     strings.TrimSpace(input.Semester),
		AcademicYear: strings.TrimSpace(input.AcademicYear),
		SchoolName:   strings.TrimSpace(input.SchoolName),
		TeacherName:  strings.TrimSpace(input.TeacherName),
		Status:       planModel.PlanStatusDraft,
	}

	if err := s.repo.CreatePlan(plan); err != nil {
		return nil, errs.Internal(err)
	}

	result := mapPlanToResponse(*plan)
	return &result, nil
}

func (s *PlanService) ListPlans(
	currentUser *authDto.UserResponse,
	query planDto.ListPlansQuery,
) ([]planDto.PlanResponse, *planDto.ListPlansMeta, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, nil, errs.Internal(err)
	}

	plans, err := s.repo.ListPlansByOwner(
		ownerID,
		strings.TrimSpace(query.Status),
		strings.TrimSpace(query.SubjectGroup),
		strings.TrimSpace(query.GradeLevel),
	)
	if err != nil {
		return nil, nil, errs.Internal(err)
	}

	result := make([]planDto.PlanResponse, 0, len(plans))
	for _, plan := range plans {
		result = append(result, mapPlanToResponse(plan))
	}

	meta := &planDto.ListPlansMeta{
		Total: len(result),
	}

	return result, meta, nil
}

func mapPlanToResponse(plan planModel.Plan) planDto.PlanResponse {
	return planDto.PlanResponse{
		ID:           plan.ID.String(),
		OwnerID:      plan.OwnerID.String(),
		Title:        plan.Title,
		SubjectGroup: plan.SubjectGroup,
		GradeLevel:   plan.GradeLevel,
		Semester:     plan.Semester,
		AcademicYear: plan.AcademicYear,
		SchoolName:   plan.SchoolName,
		TeacherName:  plan.TeacherName,
		Status:       string(plan.Status),
		CreatedAt:    plan.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    plan.UpdatedAt.Format(time.RFC3339),
	}
}
