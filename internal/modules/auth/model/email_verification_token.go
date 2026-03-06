package model

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerificationToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"type:text;not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	UsedAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}
