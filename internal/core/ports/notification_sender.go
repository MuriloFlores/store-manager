package ports

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
)

type NotificationSender interface {
	Send(ctx context.Context, data domain.EmailData) error
}
