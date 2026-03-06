package errs

import "net/http"

var (
	ErrInvalidCredentials = New(
		"AUTH_INVALID_CREDENTIALS",
		"invalid credentials",
		http.StatusUnauthorized,
	)

	ErrEmailAlreadyExists = New(
		"AUTH_EMAIL_ALREADY_EXISTS",
		"email already exists",
		http.StatusConflict,
	)

	ErrAccountDisabled = New(
		"AUTH_ACCOUNT_DISABLED",
		"account disabled",
		http.StatusForbidden,
	)

	ErrUnauthorized = New(
		"AUTH_UNAUTHORIZED",
		"unauthorized",
		http.StatusUnauthorized,
	)

	ErrMissingAuthorizationHeader = New(
		"AUTH_MISSING_AUTHORIZATION_HEADER",
		"missing authorization header",
		http.StatusUnauthorized,
	)

	ErrInvalidAuthorizationHeader = New(
		"AUTH_INVALID_AUTHORIZATION_HEADER",
		"invalid authorization header",
		http.StatusUnauthorized,
	)

	ErrUserNotFound = New(
		"AUTH_USER_NOT_FOUND",
		"user not found",
		http.StatusUnauthorized,
	)

	ErrEmailNotVerified = New(
		"AUTH_EMAIL_NOT_VERIFIED",
		"email not verified",
		http.StatusForbidden,
	)

	ErrVerificationTokenInvalid = New(
		"AUTH_VERIFICATION_TOKEN_INVALID",
		"verification token invalid",
		http.StatusBadRequest,
	)

	ErrVerificationTokenExpired = New(
		"AUTH_VERIFICATION_TOKEN_EXPIRED",
		"verification token expired",
		http.StatusBadRequest,
	)

	ErrVerificationTokenUsed = New(
		"AUTH_VERIFICATION_TOKEN_USED",
		"verification token already used",
		http.StatusBadRequest,
	)

	ErrUserAlreadyVerified = New(
		"AUTH_USER_ALREADY_VERIFIED",
		"user already verified",
		http.StatusBadRequest,
	)

	ErrInvalidCurrentUserContext = New(
		"AUTH_INVALID_CURRENT_USER_CONTEXT",
		"invalid current user context",
		http.StatusInternalServerError,
	)

	ErrPasswordResetTokenInvalid = New(
		"AUTH_PASSWORD_RESET_TOKEN_INVALID",
		"password reset token invalid",
		http.StatusBadRequest,
	)

	ErrPasswordResetTokenExpired = New(
		"AUTH_PASSWORD_RESET_TOKEN_EXPIRED",
		"password reset token expired",
		http.StatusBadRequest,
	)

	ErrPasswordResetTokenUsed = New(
		"AUTH_PASSWORD_RESET_TOKEN_USED",
		"password reset token already used",
		http.StatusBadRequest,
	)

	ErrRefreshTokenInvalid = New(
		"AUTH_REFRESH_TOKEN_INVALID",
		"refresh token invalid",
		http.StatusUnauthorized,
	)

	ErrRefreshTokenExpired = New(
		"AUTH_REFRESH_TOKEN_EXPIRED",
		"refresh token expired",
		http.StatusUnauthorized,
	)

	ErrRefreshTokenRevoked = New(
		"AUTH_REFRESH_TOKEN_REVOKED",
		"refresh token revoked",
		http.StatusUnauthorized,
	)
)
