package repository

import (
	"teacher-os-api/internal/modules/plans/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) CreatePlan(plan *model.Plan) error {
	return r.db.Create(plan).Error
}

func (r *PlanRepository) ListPlansByOwner(
	ownerID uuid.UUID,
	status string,
	subjectGroup string,
	gradeLevel string,
) ([]model.Plan, error) {
	var plans []model.Plan

	query := r.db.Where("owner_id = ?", ownerID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if subjectGroup != "" {
		query = query.Where("subject_group = ?", subjectGroup)
	}

	if gradeLevel != "" {
		query = query.Where("grade_level = ?", gradeLevel)
	}

	err := query.Order("updated_at DESC").Find(&plans).Error
	if err != nil {
		return nil, err
	}

	return plans, nil
}
