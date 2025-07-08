package auth_dto

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}
