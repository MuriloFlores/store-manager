package queue

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

const (
	TaskTypeEmailPasswordReset = "email:password_reset"
	TaskTypeEmailChange        = "email:change_confirmation"
	TaskTypeEmailSecurityAlert = "email:security_alert"
)

type taskEnqueuer struct {
	client *asynq.Client
}

func NewTaskEnqueuer(redisOpt asynq.RedisClientOpt) ports.TaskEnqueuer {
	return &taskEnqueuer{
		client: asynq.NewClient(redisOpt),
	}
}

func (t *taskEnqueuer) EnqueuePasswordReset(data *domain.PasswordResetJobData) error {
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

func (t *taskEnqueuer) EnqueueEmailChangeConfirmation(data *domain.EmailChangeConfirmationJobData) error {
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

func (t *taskEnqueuer) EnqueueSecurityNotification(data *domain.SecurityNotificationJobData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("fail in marshal data: %w", err)
	}

	task := asynq.NewTask(TaskTypeEmailSecurityAlert, payload)

	_, err = t.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("fail in enqueue security notification payload: %w", err)
	}

	return nil
}
