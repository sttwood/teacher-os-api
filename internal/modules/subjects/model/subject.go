package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectStatus string

const (
	SubjectStatusActive   SubjectStatus = "active"
	SubjectStatusArchived SubjectStatus = "archived"
)

type Subject struct {
	ID           uuid.UUID     `gorm:"type:uuid;primaryKey"`
	OwnerID      uuid.UUID     `gorm:"type:uuid;not null;index"`
	SubjectGroup string        `gorm:"type:varchar(255);not null"`
	SubjectName  string        `gorm:"type:varchar(255);not null"`
	SubjectCode  string        `gorm:"type:varchar(100)"`
	GradeLevel   string        `gorm:"type:varchar(100);not null"`
	Semester     string        `gorm:"type:varchar(50)"`
	AcademicYear string        `gorm:"type:varchar(50)"`
	SchoolName   string        `gorm:"type:varchar(255)"`
	TeacherName  string        `gorm:"type:varchar(255)"`
	Status       SubjectStatus `gorm:"type:varchar(50);not null;default:'active'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Subject) TableName() string {
	return "subjects"
}
