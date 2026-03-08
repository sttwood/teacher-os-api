package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LearningUnit struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	SubjectID   uuid.UUID `gorm:"type:uuid;not null;index"`
	UnitNo      int       `gorm:"not null"`
	UnitTitle   string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	TotalHours  int       `gorm:"not null;default:0"`
	OrderIndex  int       `gorm:"not null;default:0;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (LearningUnit) TableName() string {
	return "learning_units"
}
