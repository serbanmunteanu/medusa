package worker

import (
	"github.com/hibiken/asynq"
	"medusa/src/api/config"
	"medusa/src/api/worker/handlers"
	"medusa/src/common/env"
)

func Boot(envFiles ...string) {
	env.Load(envFiles...)

	workerConfig := &config.WorkerServerConfig{}
	config.Load("worker", workerConfig)

	workServerInstance := NewWorkServer(asynq.RedisClientOpt{
		Addr:     workerConfig.RedisConfig.Addr,
		Password: workerConfig.RedisConfig.Password,
	}, workerConfig)
	handlerMap, err := getHandlers(&workerConfig.WorkerChannel)

	if err != nil {
		panic(err)
	}

	workServerInstance.RegisterHandlers(handlerMap)

	if err = workServerInstance.Run(); err != nil {
		panic(err)
	}
}

func getHandlers(config *config.WorkerChannel) (HandlerMap, error) {
	handlerMap := make(HandlerMap)

	statusProcessor := handlers.NewStatusProcessor()
	handlerMap[config.StatusChannelKey] = statusProcessor

	return handlerMap, nil
}
