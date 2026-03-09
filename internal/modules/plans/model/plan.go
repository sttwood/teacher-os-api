package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanStatus string

const (
	PlanStatusDraft     PlanStatus = "draft"
	PlanStatusCompleted PlanStatus = "completed"
	PlanStatusArchived  PlanStatus = "archived"
)

type Plan struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	OwnerID        uuid.UUID  `gorm:"type:uuid;not null;index"`
	SubjectID      *uuid.UUID `gorm:"type:uuid;index"`
	LearningUnitID *uuid.UUID `gorm:"type:uuid;index"`
	LessonNo       int        `gorm:"default:0"`
	LessonTitle    string     `gorm:"type:varchar(255)"`
	LessonHours    int        `gorm:"default:0"`

	Title        string     `gorm:"type:varchar(255);not null"`
	SubjectGroup string     `gorm:"type:varchar(255)"`
	GradeLevel   string     `gorm:"type:varchar(100)"`
	Semester     string     `gorm:"type:varchar(50)"`
	AcademicYear string     `gorm:"type:varchar(50)"`
	SchoolName   string     `gorm:"type:varchar(255)"`
	TeacherName  string     `gorm:"type:varchar(255)"`
	Status       PlanStatus `gorm:"type:varchar(50);not null;default:'draft'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Plan) TableName() string {
	return "plans"
}
