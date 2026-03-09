package service

import (
	"errors"
	"strings"
	"time"

	authDto "teacher-os-api/internal/modules/auth/dto"
	subjectDto "teacher-os-api/internal/modules/subjects/dto"
	subjectModel "teacher-os-api/internal/modules/subjects/model"
	subjectRepository "teacher-os-api/internal/modules/subjects/repository"
	"teacher-os-api/internal/shared/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LearningUnitService struct {
	repo *subjectRepository.SubjectRepository
}

func NewLearningUnitService(repo *subjectRepository.SubjectRepository) *LearningUnitService {
	return &LearningUnitService{repo: repo}
}

func (s *LearningUnitService) CreateLearningUnit(
	currentUser *authDto.UserResponse,
	subjectID string,
	input subjectDto.CreateLearningUnitRequest,
) (*subjectDto.LearningUnitResponse, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	parsedSubjectID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errs.New("SUBJECT_INVALID_ID", "invalid subject id", 400)
	}

	_, err = s.repo.FindSubjectByIDAndOwner(parsedSubjectID, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New("SUBJECT_NOT_FOUND", "subject not found", 404)
		}
		return nil, errs.Internal(err)
	}

	unit := &subjectModel.LearningUnit{
		ID:          uuid.New(),
		SubjectID:   parsedSubjectID,
		UnitNo:      input.UnitNo,
		UnitTitle:   strings.TrimSpace(input.UnitTitle),
		Description: strings.TrimSpace(input.Description),
		TotalHours:  input.TotalHours,
		OrderIndex:  input.OrderIndex,
	}

	if err := s.repo.CreateLearningUnit(unit); err != nil {
		return nil, errs.Internal(err)
	}

	result := mapLearningUnitToResponse(*unit)
	return &result, nil
}

func (s *LearningUnitService) ListLearningUnits(
	currentUser *authDto.UserResponse,
	subjectID string,
) ([]subjectDto.LearningUnitResponse, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	parsedSubjectID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errs.New("SUBJECT_INVALID_ID", "invalid subject id", 400)
	}

	_, err = s.repo.FindSubjectByIDAndOwner(parsedSubjectID, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New("SUBJECT_NOT_FOUND", "subject not found", 404)
		}
		return nil, errs.Internal(err)
	}

	units, err := s.repo.ListLearningUnitsBySubjectID(parsedSubjectID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	result := make([]subjectDto.LearningUnitResponse, 0, len(units))
	for _, unit := range units {
		result = append(result, mapLearningUnitToResponse(unit))
	}

	return result, nil
}

func mapLearningUnitToResponse(unit subjectModel.LearningUnit) subjectDto.LearningUnitResponse {
	return subjectDto.LearningUnitResponse{
		ID:          unit.ID.String(),
		SubjectID:   unit.SubjectID.String(),
		UnitNo:      unit.UnitNo,
		UnitTitle:   unit.UnitTitle,
		Description: unit.Description,
		TotalHours:  unit.TotalHours,
		OrderIndex:  unit.OrderIndex,
		CreatedAt:   unit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   unit.UpdatedAt.Format(time.RFC3339),
	}
}
