package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

const (
	TaskTypeEmailPasswordReset = "email:password_reset"
	TaskTypeEmailChange        = "email:change_confirmation"
)

type taskEnqueuer struct {
	client *asynq.Client
}

func NewTaskEnqueuer(redisOpt asynq.RedisClientOpt) ports.TaskEnqueuer {
	return &taskEnqueuer{
		client: asynq.NewClient(redisOpt),
	}
}

func (t *taskEnqueuer) EnqueuePasswordReset(data *domain.PasswordChangeJobData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("fail in marshal data: %w", err)
	}

	task := asynq.NewTask(TaskTypeEmailPasswordReset, payload)

	_, err = t.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("fail in enqueue password rest payload: %w", err)
	}

	return nil
}

func (t *taskEnqueuer) EnqueueEmailChangeConfirmation(ctx context.Context, data *domain.EmailChangeConfirmationJobData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("fail in marshal data: %w", err)
	}

	task := asynq.NewTask(TaskTypeEmailChange, payload)

	_, err = t.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("fail in enqueue email rest payload: %w", err)
	}

	return nil
}

func (t *taskEnqueuer) EnqueueSecurityNotification(ctx context.Context, data *domain.SecurityNotificationJobData) error {
	//TODO implement me
	panic("implement me")
}
