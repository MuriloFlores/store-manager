package jobs

type EmailChangeConfirmationJobData struct {
	UserName         string
	ConfirmationLink string
	ToEmail          string
}

type SecurityNotificationJobData struct {
	UserName string
	ToEmail  string
	Message  string
}
