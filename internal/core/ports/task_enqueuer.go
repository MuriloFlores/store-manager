package ports

import (
	"github.com/muriloFlores/StoreManager/internal/core/domain/jobs"
)

type TaskEnqueuer interface {
	EnqueuePasswordReset(data *jobs.PasswordResetJobData) error
	EnqueueEmailChangeConfirmation(data *jobs.EmailChangeConfirmationJobData) error
	EnqueueSecurityNotification(data *jobs.SecurityNotificationJobData) error
	EnqueueAccountVerification(data *jobs.AccountVerificationJobData) error
	EnqueuePromotionNotification(data *jobs.PromotionNotificationJobData) error
}
