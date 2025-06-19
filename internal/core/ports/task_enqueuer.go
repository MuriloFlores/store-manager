package ports

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
)

type TaskEnqueuer interface {
	EnqueuePasswordReset(data *domain.PasswordChangeJobData) error
	EnqueueEmailChangeConfirmation(ctx context.Context, data *domain.EmailChangeConfirmationJobData) error
	EnqueueSecurityNotification(ctx context.Context, data *domain.SecurityNotificationJobData) error
}
