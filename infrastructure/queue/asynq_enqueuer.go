package queue

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/muriloFlores/StoreManager/internal/core/domain/jobs"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

const (
	TaskTypeEmailPasswordReset         = "email:password_reset"
	TaskTypeEmailChange                = "email:change_confirmation"
	TaskTypeEmailSecurityAlert         = "email:security_alert"
	TaskTypeEmailAccountVerification   = "email:account_verification"
	TaskTypeEmailPromotionNotification = "email:promotion_notification"
)

type taskEnqueuer struct {
	client *asynq.Client
	logger ports.Logger
}

func NewTaskEnqueuer(redisOpt asynq.RedisClientOpt, logger ports.Logger) ports.TaskEnqueuer {
	return &taskEnqueuer{
		client: asynq.NewClient(redisOpt),
		logger: logger,
	}
}

func (t *taskEnqueuer) EnqueuePasswordReset(data *jobs.PasswordResetJobData) error {
	t.logger.InfoLevel("Enqueuing password reset task", map[string]interface{}{"user_name": data.UserName, "email": data.UserEmail})

	payload, err := json.Marshal(data)
	if err != nil {
		t.logger.ErrorLevel("Failed to marshal password reset data", err, map[string]interface{}{"user_name": data.UserName, "email": data.UserEmail})
		return fmt.Errorf("fail in marshal data: %w", err)
	}

	task := asynq.NewTask(TaskTypeEmailPasswordReset, payload)

	_, err = t.client.Enqueue(task)
	if err != nil {
		t.logger.ErrorLevel("Failed to enqueue password reset task", err, map[string]interface{}{"user_name": data.UserName, "email": data.UserEmail})
		return fmt.Errorf("fail in enqueue password rest payload: %w", err)
	}

	t.logger.InfoLevel("Password reset task enqueued successfully", map[string]interface{}{"user_name": data.UserName, "email": data.UserEmail})
	return nil
}

func (t *taskEnqueuer) EnqueueEmailChangeConfirmation(data *jobs.EmailChangeConfirmationJobData) error {
	t.logger.InfoLevel("Enqueuing email change confirmation task", map[string]interface{}{"user_name": data.UserName, "user_email": data.ToEmail})

	payload, err := json.Marshal(data)
	if err != nil {
		t.logger.ErrorLevel("Failed to marshal email change confirmation data", err, map[string]interface{}{"user_name": data.UserName, "user_email": data.ToEmail})
		return fmt.Errorf("fail in marshal data: %w", err)
	}

	task := asynq.NewTask(TaskTypeEmailChange, payload)

	_, err = t.client.Enqueue(task)
	if err != nil {
		t.logger.ErrorLevel("Failed to enqueue email change confirmation task", err, map[string]interface{}{"user_name": data.UserName, "user_email": data.ToEmail})
		return fmt.Errorf("fail in enqueue email rest payload: %w", err)
	}

	t.logger.InfoLevel("Email change confirmation task enqueued successfully", map[string]interface{}{"user_name": data.UserName, "user_email": data.ToEmail})
	return nil
}

func (t *taskEnqueuer) EnqueueSecurityNotification(data *jobs.SecurityNotificationJobData) error {
	t.logger.InfoLevel("Enqueuing security alert task", map[string]interface{}{"user_name": data.UserName, "user_email": data.ToEmail})

	payload, err := json.Marshal(data)
	if err != nil {
		t.logger.ErrorLevel("Failed to marshal security alert data", err, map[string]interface{}{"user_name": data.UserName, "user_email": data.ToEmail})
		return fmt.Errorf("fail in marshal data: %w", err)
	}

	task := asynq.NewTask(TaskTypeEmailSecurityAlert, payload)

	_, err = t.client.Enqueue(task)
	if err != nil {
		t.logger.ErrorLevel("Failed to enqueue security notification task", err, map[string]interface{}{"user_name": data.UserName, "user_email": data.ToEmail})
		return fmt.Errorf("fail in enqueue security notification payload: %w", err)
	}

	t.logger.InfoLevel("Security notification task enqueued successfully", map[string]interface{}{"user_name": data.UserName, "user_email": data.ToEmail})
	return nil
}

func (t *taskEnqueuer) EnqueueAccountVerification(data *jobs.AccountVerificationJobData) error {
	t.logger.InfoLevel("Enqueuing account verification task", map[string]interface{}{"user_name": data.UserName, "to_email": data.ToEmail})

	payload, err := json.Marshal(data)
	if err != nil {
		t.logger.ErrorLevel("Failed to marshal account verification data", err, map[string]interface{}{"user_name": data.UserName, "to_email": data.ToEmail})
		return fmt.Errorf("fail in marshal data: %w", err)
	}

	task := asynq.NewTask(TaskTypeEmailAccountVerification, payload)

	_, err = t.client.Enqueue(task)
	if err != nil {
		t.logger.ErrorLevel("Failed to enqueue account verification task", err, map[string]interface{}{"user_name": data.UserName, "to_email": data.ToEmail})
		return fmt.Errorf("fail in enqueue account verification payload: %w", err)
	}

	t.logger.InfoLevel("Account verification task enqueued successfully", map[string]interface{}{"user_name": data.UserName, "to_email": data.ToEmail})
	return nil
}

func (t *taskEnqueuer) EnqueuePromotionNotification(data *jobs.PromotionNotificationJobData) error {
	t.logger.InfoLevel("Enqueuing promotion notification task", map[string]interface{}{"user_name": data.UserName, "to_email": data.ToEmail, "new_role": data.NewRole})

	payload, err := json.Marshal(data)
	if err != nil {
		t.logger.ErrorLevel("Failed to marshal promotion notification data", err, map[string]interface{}{"user_name": data.UserName, "to_email": data.ToEmail, "new_role": data.NewRole})
		return fmt.Errorf("fail in marshal data: %w", err)
	}

	task := asynq.NewTask("email:promotion_notification", payload)

	_, err = t.client.Enqueue(task)
	if err != nil {
		t.logger.ErrorLevel("Failed to enqueue promotion notification task", err, map[string]interface{}{"user_name": data.UserName, "to_email": data.ToEmail, "new_role": data.NewRole})
		return fmt.Errorf("fail in enqueue promotion notification payload: %w", err)
	}

	t.logger.InfoLevel("Promotion notification task enqueued successfully", map[string]interface{}{"user_name": data.UserName, "to_email": data.ToEmail, "new_role": data.NewRole})
	return nil
}
