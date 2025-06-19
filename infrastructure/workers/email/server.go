package email

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/muriloFlores/StoreManager/infrastructure/queue"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

func RunTaskServer(
	redisOpt asynq.RedisClientOpt,
	workerOpt ports.Worker,
) {
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(
		queue.TaskTypeEmailPasswordReset,
		func(ctx context.Context, t *asynq.Task) error {
			return workerOpt.HandlePasswordReset(ctx, t.Payload())
		},
	)

	mux.HandleFunc(
		queue.TaskTypeEmailChange,
		func(ctx context.Context, t *asynq.Task) error {
			return workerOpt.HandleEmailChangeTask(ctx, t.Payload())
		},
	)

	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}
