package repository

import (
	"teacher-os-api/internal/modules/auth/model"
	"time"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *AuthRepository) CreateEmailVerificationToken(token *model.EmailVerificationToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) FindEmailVerificationTokenByHash(tokenHash string) (*model.EmailVerificationToken, error) {
	var token model.EmailVerificationToken
	err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *AuthRepository) MarkEmailVerificationTokenUsed(id string, usedAt time.Time) error {
	return r.db.Model(&model.EmailVerificationToken{}).
		Where("id = ?", id).
		Update("used_at", usedAt).Error
}

func (r *AuthRepository) CreatePasswordResetToken(token *model.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) FindPasswordResetTokenByHash(tokenHash string) (*model.PasswordResetToken, error) {
	var token model.PasswordResetToken
	err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *AuthRepository) MarkPasswordResetTokenUsed(id string, usedAt time.Time) error {
	return r.db.Model(&model.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", usedAt).Error
}

func (r *AuthRepository) CreateRefreshToken(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) FindRefreshTokenByHash(tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *AuthRepository) RevokeRefreshToken(id string, revokedAt time.Time) error {
	return r.db.Model(&model.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked_at", revokedAt).Error
}
