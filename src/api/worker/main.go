package worker

import (
	"github.com/hibiken/asynq"
	"medusa/src/api/config"
	"medusa/src/common/env"
)

func Boot(envFiles ...string) {
	env.Load(envFiles...)

	workerConfig := &config.WorkerServerConfig{}
	config.Load("worker", workerConfig)

	handlersMap := make(HandlersMap)
	redisOptions := asynq.RedisClientOpt{
		Addr:     workerConfig.RedisConfig.Addr,
		Password: workerConfig.RedisConfig.Password,
	}
	workServerInstance := NewWorkServer(redisOptions, config.WorkerServerConfig{})
	workServerInstance.RegisterHandlers(handlersMap)

	if err := workServerInstance.Run(); err != nil {
		panic(err)
	}
}
