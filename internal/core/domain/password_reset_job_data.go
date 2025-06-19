package domain

type PasswordChangeJobData struct {
	UserName  string
	UserEmail string
	ResetLink string
}
