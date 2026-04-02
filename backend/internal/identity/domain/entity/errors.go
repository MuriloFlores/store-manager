package entity

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserIsDeactivated  = errors.New("user is deactivated")
	ErrInvalidOldPassword = errors.New("invalid old password")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionNotFound    = errors.New("session not found")
	ErrOTPNotFound        = errors.New("otp not found")
)
