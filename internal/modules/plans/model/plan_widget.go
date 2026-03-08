package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PlanWidgetType string

const (
	PlanWidgetTypeBasicInfo         PlanWidgetType = "basicInfo"
	PlanWidgetTypeObjective         PlanWidgetType = "objective"
	PlanWidgetTypeStandardIndicator PlanWidgetType = "standardIndicator"
	PlanWidgetTypeActivity          PlanWidgetType = "activity"
	PlanWidgetTypeMediaMaterial     PlanWidgetType = "mediaMaterial"
	PlanWidgetTypeAssessment        PlanWidgetType = "assessment"
	PlanWidgetTypeHomework          PlanWidgetType = "homework"
)

type PlanWidget struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	PlanID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	WidgetType  PlanWidgetType `gorm:"type:varchar(100);not null;index"`
	OrderIndex  int            `gorm:"not null;index"`
	Title       string         `gorm:"type:varchar(255)"`
	ContentJSON datatypes.JSON `gorm:"type:jsonb;not null"`
	IsCollapsed bool           `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (PlanWidget) TableName() string {
	return "plan_widgets"
}
