package ports

import (
	"github.com/muriloFlores/StoreManager/internal/core/domain"
)

type TaskEnqueuer interface {
	EnqueuePasswordReset(data *domain.PasswordResetJobData) error
	EnqueueEmailChangeConfirmation(data *domain.EmailChangeConfirmationJobData) error
	EnqueueSecurityNotification(data *domain.SecurityNotificationJobData) error
}
