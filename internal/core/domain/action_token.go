package domain

import "time"

type ActionType string

const (
	PasswordReset       ActionType = "PASSWORD_RESET"
	EmailConfirmation   ActionType = "EMAIL_CONFIRMATION"
	AccountVerification ActionType = "ACCOUNT_VERIFICATION"
)

type ActionToken struct {
	Token     string
	UserID    string
	Type      ActionType
	Payload   string
	ExpiresAt time.Time
}
