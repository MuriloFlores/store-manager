package jobs

type PasswordResetJobData struct {
	UserName  string
	UserEmail string
	ResetLink string
}
