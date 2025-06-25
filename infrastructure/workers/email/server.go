package email

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/muriloFlores/StoreManager/infrastructure/queue"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"log"
)

func RunTaskServer(
	redisOpt asynq.RedisClientOpt,
	workerOpt ports.Worker,
	logger ports.Logger,
) {
	logger.InfoLevel("Starting email task server...")

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.ErrorLevel("Error processing task", err, map[string]interface{}{"task_payload": task.Payload(), "task_type": task.Type()})
			}),
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(
		queue.TaskTypeEmailPasswordReset,
		func(ctx context.Context, t *asynq.Task) error {
			logger.InfoLevel("Processing task", map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})

			err := workerOpt.HandlePasswordReset(ctx, t.Payload())

			if err != nil {
				logger.ErrorLevel("Error processing task", err, map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})
				return err
			}

			logger.InfoLevel("Task processed successfully", map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})

			return nil
		},
	)

	mux.HandleFunc(
		queue.TaskTypeEmailChange,
		func(ctx context.Context, t *asynq.Task) error {
			logger.InfoLevel("Processing task", map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})

			err := workerOpt.HandleEmailChangeTask(ctx, t.Payload())

			if err != nil {
				logger.ErrorLevel("Error processing task", err, map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})
				return err
			}

			logger.InfoLevel("Task processed successfully", map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})

			return nil
		},
	)

	mux.HandleFunc(
		queue.TaskTypeEmailAccountVerification,
		func(ctx context.Context, t *asynq.Task) error {
			logger.InfoLevel("Processing task", map[string]interface{}{"task_type": t.Type()})

			err := workerOpt.HandleAccountVerification(ctx, t.Payload())

			if err != nil {
				logger.ErrorLevel("Error processing task", err, map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})
				return err
			}

			logger.InfoLevel("Task processed successfully", map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})
			return nil
		},
	)

	mux.HandleFunc(
		queue.TaskTypeEmailPromotionNotification,
		func(ctx context.Context, t *asynq.Task) error {
			logger.InfoLevel("Processing task", map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})

			err := workerOpt.HandlePromotionNotification(ctx, t.Payload())

			if err != nil {
				logger.ErrorLevel("Error processing task", err, map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})
				return err
			}

			logger.InfoLevel("Task processed successfully", map[string]interface{}{"task_payload": t.Payload(), "task_type": t.Type()})
			return nil
		},
	)

	if err := srv.Run(mux); err != nil {
		log.Fatal("Could not run task server", err)
	}
}
