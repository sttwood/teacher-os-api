package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"teacher-os-api/internal/modules/auth/dto"
	"teacher-os-api/internal/modules/auth/model"
	"teacher-os-api/internal/modules/auth/repository"
	"teacher-os-api/internal/shared/errs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const emailVerificationTokenTTL = 24 * time.Hour
const passwordResetTokenTTL = 1 * time.Hour
const accessTokenTTL = 15 * time.Minute
const refreshTokenTTL = 7 * 24 * time.Hour

type AuthService struct {
	repo      *repository.AuthRepository
	jwtSecret string
}

func NewAuthService(repo *repository.AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(input dto.RegisterRequest) (*dto.RegisterResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	fullName := strings.TrimSpace(input.FullName)
	schoolName := strings.TrimSpace(input.SchoolName)

	existingUser, err := s.repo.FindUserByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errs.ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Internal(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.Internal(err)
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		SchoolName:   schoolName,
		IsVerified:   false,
		IsActive:     true,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, errs.Internal(err)
	}

	plainToken, err := generateRandomToken(32)
	if err != nil {
		return nil, errs.Internal(err)
	}

	tokenRecord := &model.EmailVerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hashToken(plainToken),
		ExpiresAt: time.Now().Add(emailVerificationTokenTTL),
	}

	if err := s.repo.CreateEmailVerificationToken(tokenRecord); err != nil {
		return nil, errs.Internal(err)
	}

	// DEV ONLY:
	// ตอนนี้ยังไม่ส่ง email จริง ให้ print token/link ใน console ก่อน
	// ภายหลัง Step mail integration จะย้ายไป mail service
	println("VERIFY EMAIL TOKEN:", plainToken)
	println("VERIFY EMAIL LINK: http://localhost:3000/verify-email?token=" + plainToken)

	return &dto.RegisterResponse{
		Message: "register success",
		User: dto.UserResponse{
			ID:         user.ID.String(),
			Email:      user.Email,
			FullName:   user.FullName,
			SchoolName: user.SchoolName,
			IsVerified: user.IsVerified,
			IsActive:   user.IsActive,
		},
	}, nil
}

func (s *AuthService) Login(input dto.LoginRequest) (*dto.LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, errs.Internal(err)
	}

	if !user.IsActive {
		return nil, errs.ErrAccountDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, errs.ErrEmailNotVerified
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errs.Internal(err)
	}

	refreshToken, err := s.createRefreshToken(user)
	if err != nil {
		return nil, errs.Internal(err)
	}

	return &dto.LoginResponse{
		Message:      "login success",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:         user.ID.String(),
			Email:      user.Email,
			FullName:   user.FullName,
			SchoolName: user.SchoolName,
			IsVerified: user.IsVerified,
			IsActive:   user.IsActive,
		},
	}, nil
}

func (s *AuthService) Logout(input dto.LogoutRequest) (*dto.LogoutResponse, error) {
	token := strings.TrimSpace(input.RefreshToken)
	if token == "" {
		return nil, errs.ErrRefreshTokenInvalid
	}

	tokenHash := hashToken(token)

	tokenRecord, err := s.repo.FindRefreshTokenByHash(tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.LogoutResponse{
				Message: "logout success",
			}, nil
		}
		return nil, errs.Internal(err)
	}

	if tokenRecord.RevokedAt == nil {
		now := time.Now()
		if err := s.repo.RevokeRefreshToken(tokenRecord.ID.String(), now); err != nil {
			return nil, errs.Internal(err)
		}
	}

	return &dto.LogoutResponse{
		Message: "logout success",
	}, nil
}

func (s *AuthService) GetUserByID(id string) (*dto.UserResponse, error) {
	user, err := s.repo.FindUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, errs.Internal(err)
	}

	if !user.IsActive {
		return nil, errs.ErrAccountDisabled
	}

	if !user.IsVerified {
		return nil, errs.ErrEmailNotVerified
	}

	return &dto.UserResponse{
		ID:         user.ID.String(),
		Email:      user.Email,
		FullName:   user.FullName,
		SchoolName: user.SchoolName,
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
	}, nil
}

func (s *AuthService) ResendVerifyEmail(input dto.ResendVerifyEmailRequest) (*dto.ResendVerifyEmailResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, errs.Internal(err)
	}

	if !user.IsActive {
		return nil, errs.ErrAccountDisabled
	}

	if user.IsVerified {
		return nil, errs.ErrUserAlreadyVerified
	}

	plainToken, err := generateRandomToken(32)
	if err != nil {
		return nil, errs.Internal(err)
	}

	tokenRecord := &model.EmailVerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hashToken(plainToken),
		ExpiresAt: time.Now().Add(emailVerificationTokenTTL),
	}

	if err := s.repo.CreateEmailVerificationToken(tokenRecord); err != nil {
		return nil, errs.Internal(err)
	}

	// DEV ONLY
	println("RESEND VERIFY EMAIL TOKEN:", plainToken)
	println("RESEND VERIFY EMAIL LINK: http://localhost:3000/verify-email?token=" + plainToken)

	return &dto.ResendVerifyEmailResponse{
		Message: "verification email sent",
	}, nil
}

func (s *AuthService) ConfirmVerifyEmail(input dto.ConfirmVerifyEmailRequest) (*dto.ConfirmVerifyEmailResponse, error) {
	// !TODO production กว่าเดิมควรใช้ transaction
	token := strings.TrimSpace(input.Token)
	if token == "" {
		return nil, errs.ErrVerificationTokenInvalid
	}

	tokenHash := hashToken(token)

	tokenRecord, err := s.repo.FindEmailVerificationTokenByHash(tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrVerificationTokenInvalid
		}
		return nil, errs.Internal(err)
	}

	if tokenRecord.UsedAt != nil {
		return nil, errs.ErrVerificationTokenUsed
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, errs.ErrVerificationTokenExpired
	}

	user, err := s.repo.FindUserByID(tokenRecord.UserID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, errs.Internal(err)
	}

	if user.IsVerified {
		return nil, errs.ErrUserAlreadyVerified
	}

	user.IsVerified = true
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, errs.Internal(err)
	}

	now := time.Now()
	if err := s.repo.MarkEmailVerificationTokenUsed(tokenRecord.ID.String(), now); err != nil {
		return nil, errs.Internal(err)
	}

	return &dto.ConfirmVerifyEmailResponse{
		Message: "email verified successfully",
	}, nil
}

func (s *AuthService) ForgotPassword(input dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.ForgotPasswordResponse{
				Message: "if the account exists, a reset link has been sent",
			}, nil
		}
		return nil, errs.Internal(err)
	}

	if !user.IsActive {
		return &dto.ForgotPasswordResponse{
			Message: "if the account exists, a reset link has been sent",
		}, nil
	}

	plainToken, err := generateRandomToken(32)
	if err != nil {
		return nil, errs.Internal(err)
	}

	tokenRecord := &model.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hashToken(plainToken),
		ExpiresAt: time.Now().Add(passwordResetTokenTTL),
	}

	if err := s.repo.CreatePasswordResetToken(tokenRecord); err != nil {
		return nil, errs.Internal(err)
	}

	// DEV ONLY
	println("PASSWORD RESET TOKEN:", plainToken)
	println("PASSWORD RESET LINK: http://localhost:3000/reset-password?token=" + plainToken)

	return &dto.ForgotPasswordResponse{
		Message: "if the account exists, a reset link has been sent",
	}, nil
}

func (s *AuthService) ResetPassword(input dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	token := strings.TrimSpace(input.Token)
	if token == "" {
		return nil, errs.ErrPasswordResetTokenInvalid
	}

	tokenHash := hashToken(token)

	tokenRecord, err := s.repo.FindPasswordResetTokenByHash(tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPasswordResetTokenInvalid
		}
		return nil, errs.Internal(err)
	}

	if tokenRecord.UsedAt != nil {
		return nil, errs.ErrPasswordResetTokenUsed
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, errs.ErrPasswordResetTokenExpired
	}

	user, err := s.repo.FindUserByID(tokenRecord.UserID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, errs.Internal(err)
	}

	if !user.IsActive {
		return nil, errs.ErrAccountDisabled
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.Internal(err)
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, errs.Internal(err)
	}

	now := time.Now()
	if err := s.repo.MarkPasswordResetTokenUsed(tokenRecord.ID.String(), now); err != nil {
		return nil, errs.Internal(err)
	}

	return &dto.ResetPasswordResponse{
		Message: "password reset successfully",
	}, nil
}

func (s *AuthService) Refresh(input dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	/*
		!TODO ควรใช้ transaction ครอบขั้นตอน:
		- revoke old refresh token
		- create new refresh token
		อย่างน้อยส่วน refresh token rotation ควร atomic
	*/
	token := strings.TrimSpace(input.RefreshToken)
	if token == "" {
		return nil, errs.ErrRefreshTokenInvalid
	}

	tokenHash := hashToken(token)

	tokenRecord, err := s.repo.FindRefreshTokenByHash(tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRefreshTokenInvalid
		}
		return nil, errs.Internal(err)
	}

	if tokenRecord.RevokedAt != nil {
		return nil, errs.ErrRefreshTokenRevoked
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, errs.ErrRefreshTokenExpired
	}

	user, err := s.repo.FindUserByID(tokenRecord.UserID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, errs.Internal(err)
	}

	if !user.IsActive {
		return nil, errs.ErrAccountDisabled
	}

	if !user.IsVerified {
		return nil, errs.ErrEmailNotVerified
	}

	now := time.Now()
	if err := s.repo.RevokeRefreshToken(tokenRecord.ID.String(), now); err != nil {
		return nil, errs.Internal(err)
	}

	newAccessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.createRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		Message:      "token refreshed",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Helpers
// ===================================================
func (s *AuthService) ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS512 {
			return nil, errs.ErrUnauthorized
		}

		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errs.ErrUnauthorized
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", errs.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errs.ErrUnauthorized
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", errs.ErrUnauthorized
	}

	return sub, nil
}

func (s *AuthService) createRefreshToken(user *model.User) (string, error) {
	plainToken, err := generateRandomToken(32)
	if err != nil {
		return "", errs.Internal(err)
	}

	tokenRecord := &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hashToken(plainToken),
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}

	if err := s.repo.CreateRefreshToken(tokenRecord); err != nil {
		return "", errs.Internal(err)
	}

	return plainToken, nil
}

func (s *AuthService) generateAccessToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"exp":   time.Now().Add(accessTokenTTL).Unix(),
		"iat":   time.Now().Unix(),
	})

	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errs.Internal(err)
	}

	return accessToken, nil
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
