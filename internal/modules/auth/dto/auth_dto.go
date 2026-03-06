package dto

type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	FullName   string `json:"fullName" binding:"required,min=2"`
	SchoolName string `json:"schoolName"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	SchoolName string `json:"schoolName"`
	IsVerified bool   `json:"isVerified"`
	IsActive   bool   `json:"isActive"`
	// CreatedAt  string `json:"createdAt"`
	// UpdatedAt  string `json:"updatedAt"`
}

type RegisterResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

type LoginResponse struct {
	Message      string       `json:"message"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
}

type MeResponse struct {
	User UserResponse `json:"user"`
}

type ResendVerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResendVerifyEmailResponse struct {
	Message string `json:"message"`
}

type ConfirmVerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type ConfirmVerifyEmailResponse struct {
	Message string `json:"message"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshTokenResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
