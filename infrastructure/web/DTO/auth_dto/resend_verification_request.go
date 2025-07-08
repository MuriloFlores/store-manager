package authDTO

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}
