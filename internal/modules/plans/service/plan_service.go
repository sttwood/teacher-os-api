package service

import (
	"errors"
	"strings"
	"time"

	authDto "teacher-os-api/internal/modules/auth/dto"
	planDto "teacher-os-api/internal/modules/plans/dto"
	planModel "teacher-os-api/internal/modules/plans/model"
	planRepository "teacher-os-api/internal/modules/plans/repository"
	subjectRepository "teacher-os-api/internal/modules/subjects/repository"
	"teacher-os-api/internal/shared/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanService struct {
	repo        *planRepository.PlanRepository
	subjectRepo *subjectRepository.SubjectRepository
}

func NewPlanService(
	repo *planRepository.PlanRepository,
	subjectRepo *subjectRepository.SubjectRepository,
) *PlanService {
	return &PlanService{
		repo:        repo,
		subjectRepo: subjectRepo,
	}
}

func (s *PlanService) CreatePlan(
	currentUser *authDto.UserResponse,
	input planDto.CreatePlanRequest,
) (*planDto.PlanResponse, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	subjectID, err := uuid.Parse(input.SubjectID)
	if err != nil {
		return nil, errs.New("SUBJECT_INVALID_ID", "invalid subject id", 400)
	}

	learningUnitID, err := uuid.Parse(input.LearningUnitID)
	if err != nil {
		return nil, errs.New("LEARNING_UNIT_INVALID_ID", "invalid learning unit id", 400)
	}

	subject, err := s.subjectRepo.FindSubjectByIDAndOwner(subjectID, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New("SUBJECT_NOT_FOUND", "subject not found", 404)
		}
		return nil, errs.Internal(err)
	}

	_, err = s.subjectRepo.FindLearningUnitByIDAndSubjectID(learningUnitID, subjectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New("LEARNING_UNIT_NOT_FOUND", "learning unit not found", 404)
		}
		return nil, errs.Internal(err)
	}

	plan := &planModel.Plan{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		SubjectID:      &subjectID,
		LearningUnitID: &learningUnitID,
		LessonNo:       input.LessonNo,
		LessonTitle:    strings.TrimSpace(input.LessonTitle),
		LessonHours:    input.LessonHours,

		Title:        strings.TrimSpace(input.Title),
		SubjectGroup: subject.SubjectGroup,
		GradeLevel:   subject.GradeLevel,
		Semester:     subject.Semester,
		AcademicYear: subject.AcademicYear,
		SchoolName:   subject.SchoolName,
		TeacherName:  subject.TeacherName,
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

func (s *PlanService) GetPlanByID(
	currentUser *authDto.UserResponse,
	planID string,
) (*planDto.PlanResponse, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	parsedPlanID, err := uuid.Parse(planID)
	if err != nil {
		return nil, errs.New("PLAN_INVALID_ID", "invalid plan id", 400)
	}

	plan, err := s.repo.FindPlanByIDAndOwner(parsedPlanID, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New("PLAN_NOT_FOUND", "plan not found", 404)
		}
		return nil, errs.Internal(err)
	}

	result := mapPlanToResponse(*plan)
	return &result, nil
}

func mapPlanToResponse(plan planModel.Plan) planDto.PlanResponse {
	subjectID := ""
	if plan.SubjectID != nil {
		subjectID = plan.SubjectID.String()
	}

	learningUnitID := ""
	if plan.LearningUnitID != nil {
		learningUnitID = plan.LearningUnitID.String()
	}

	return planDto.PlanResponse{
		ID:             plan.ID.String(),
		OwnerID:        plan.OwnerID.String(),
		SubjectID:      subjectID,
		LearningUnitID: learningUnitID,
		LessonNo:       plan.LessonNo,
		LessonTitle:    plan.LessonTitle,
		LessonHours:    plan.LessonHours,

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
