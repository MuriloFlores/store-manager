package ports

import "context"

type Worker interface {
	HandlePasswordReset(ctx context.Context, taskPayload []byte) error
	HandleEmailChangeTask(ctx context.Context, taskPayload []byte) error
	HandleAccountVerification(ctx context.Context, taskPayload []byte) error
	HandlePromotionNotification(ctx context.Context, taskPayload []byte) error
}
