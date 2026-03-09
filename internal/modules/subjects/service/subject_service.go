package service

import (
	"strings"
	"time"

	authDto "teacher-os-api/internal/modules/auth/dto"
	subjectDto "teacher-os-api/internal/modules/subjects/dto"
	subjectModel "teacher-os-api/internal/modules/subjects/model"
	subjectRepository "teacher-os-api/internal/modules/subjects/repository"
	"teacher-os-api/internal/shared/errs"

	"github.com/google/uuid"
)

type SubjectService struct {
	repo *subjectRepository.SubjectRepository
}

func NewSubjectService(repo *subjectRepository.SubjectRepository) *SubjectService {
	return &SubjectService{repo: repo}
}

func (s *SubjectService) CreateSubject(
	currentUser *authDto.UserResponse,
	input subjectDto.CreateSubjectRequest,
) (*subjectDto.SubjectResponse, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	subject := &subjectModel.Subject{
		ID:           uuid.New(),
		OwnerID:      ownerID,
		SubjectGroup: strings.TrimSpace(input.SubjectGroup),
		SubjectName:  strings.TrimSpace(input.SubjectName),
		SubjectCode:  strings.TrimSpace(input.SubjectCode),
		GradeLevel:   strings.TrimSpace(input.GradeLevel),
		Semester:     strings.TrimSpace(input.Semester),
		AcademicYear: strings.TrimSpace(input.AcademicYear),
		SchoolName:   strings.TrimSpace(input.SchoolName),
		TeacherName:  strings.TrimSpace(input.TeacherName),
		Status:       subjectModel.SubjectStatusActive,
	}

	if err := s.repo.CreateSubject(subject); err != nil {
		return nil, errs.Internal(err)
	}

	result := mapSubjectToResponse(*subject)
	return &result, nil
}

func (s *SubjectService) ListSubjects(
	currentUser *authDto.UserResponse,
) ([]subjectDto.SubjectResponse, error) {
	ownerID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	subjects, err := s.repo.ListSubjectsByOwner(ownerID)
	if err != nil {
		return nil, errs.Internal(err)
	}

	result := make([]subjectDto.SubjectResponse, 0, len(subjects))
	for _, subject := range subjects {
		result = append(result, mapSubjectToResponse(subject))
	}

	return result, nil
}

func mapSubjectToResponse(subject subjectModel.Subject) subjectDto.SubjectResponse {
	return subjectDto.SubjectResponse{
		ID:           subject.ID.String(),
		OwnerID:      subject.OwnerID.String(),
		SubjectGroup: subject.SubjectGroup,
		SubjectName:  subject.SubjectName,
		SubjectCode:  subject.SubjectCode,
		GradeLevel:   subject.GradeLevel,
		Semester:     subject.Semester,
		AcademicYear: subject.AcademicYear,
		SchoolName:   subject.SchoolName,
		TeacherName:  subject.TeacherName,
		Status:       string(subject.Status),
		CreatedAt:    subject.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    subject.UpdatedAt.Format(time.RFC3339),
	}
}
