package repository

import (
	subjectModel "teacher-os-api/internal/modules/subjects/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

func (r *SubjectRepository) CreateSubject(subject *subjectModel.Subject) error {
	return r.db.Create(subject).Error
}

func (r *SubjectRepository) ListSubjectsByOwner(ownerID uuid.UUID) ([]subjectModel.Subject, error) {
	var subjects []subjectModel.Subject

	err := r.db.
		Where("owner_id = ?", ownerID).
		Order("updated_at DESC").
		Find(&subjects).Error
	if err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *SubjectRepository) FindSubjectByIDAndOwner(subjectID uuid.UUID, ownerID uuid.UUID) (*subjectModel.Subject, error) {
	var subject subjectModel.Subject

	err := r.db.
		Where("id = ? AND owner_id = ?", subjectID, ownerID).
		First(&subject).Error
	if err != nil {
		return nil, err
	}

	return &subject, nil
}

func (r *SubjectRepository) CreateLearningUnit(unit *subjectModel.LearningUnit) error {
	return r.db.Create(unit).Error
}

func (r *SubjectRepository) ListLearningUnitsBySubjectID(subjectID uuid.UUID) ([]subjectModel.LearningUnit, error) {
	var units []subjectModel.LearningUnit

	err := r.db.
		Where("subject_id = ?", subjectID).
		Order("order_index ASC, unit_no ASC, created_at ASC").
		Find(&units).Error
	if err != nil {
		return nil, err
	}

	return units, nil
}

func (r *SubjectRepository) FindLearningUnitByIDAndSubjectID(
	unitID uuid.UUID,
	subjectID uuid.UUID,
) (*subjectModel.LearningUnit, error) {
	var unit subjectModel.LearningUnit

	err := r.db.
		Where("id = ? AND subject_id = ?", unitID, subjectID).
		First(&unit).Error
	if err != nil {
		return nil, err
	}

	return &unit, nil
}
