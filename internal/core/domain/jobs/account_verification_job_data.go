package jobs

type AccountVerificationJobData struct {
	UserName         string
	ToEmail          string
	VerificationLink string
}
