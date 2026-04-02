package ports

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
)

type NotificationService interface {
	SendChangePasswordEmail(ctx context.Context, toEmail vo.Email, resetToken vo.OTP) error
	SendForgotPasswordEmail(ctx context.Context, toEmail vo.Email, resetToken vo.OTP) error
}
