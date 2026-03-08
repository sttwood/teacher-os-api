package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PlanVersion struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	PlanID       uuid.UUID      `gorm:"type:uuid;not null;index"`
	VersionNo    int            `gorm:"not null"`
	SnapshotJSON datatypes.JSON `gorm:"type:jsonb;not null"`
	CreatedBy    uuid.UUID      `gorm:"type:uuid;not null;index"`
	CreatedAt    time.Time
}

func (PlanVersion) TableName() string {
	return "plan_versions"
}
