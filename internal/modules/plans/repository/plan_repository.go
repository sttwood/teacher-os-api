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

func (r *PlanRepository) FindPlanByIDAndOwner(planID uuid.UUID, ownerID uuid.UUID) (*model.Plan, error) {
	var plan model.Plan
	err := r.db.
		Where("id = ? AND owner_id = ?", planID, ownerID).
		First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) ListWidgetsByPlanID(planID uuid.UUID) ([]model.PlanWidget, error) {
	var widgets []model.PlanWidget
	err := r.db.Where("plan_id = ?", planID).Order("order_index ASC").Find(&widgets).Error
	if err != nil {
		return nil, err
	}
	return widgets, nil
}

func (r *PlanRepository) ReplaceWidgets(planID uuid.UUID, widgets []model.PlanWidget) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Where("plan_id = ?", planID).
			Delete(&model.PlanWidget{}).Error; err != nil {
			return err
		}

		if len(widgets) == 0 {
			return nil
		}

		return tx.Create(&widgets).Error
	})
}
