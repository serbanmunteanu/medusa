package worker

import (
	"github.com/hibiken/asynq"
	"medusa/src/api/config"
)

type HandlersMap map[string]asynq.Handler

type workServer struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewWorkServer(redisOptions asynq.RedisClientOpt, config config.WorkerServerConfig) *workServer {
	var queues map[string]int

	if config.UseCustomQueue {
		for _, queue := range config.Queues {
			queues[queue.Name] = queue.Priority
		}
	}

	server := asynq.NewServer(
		redisOptions,
		asynq.Config{
			Concurrency:    config.Concurrency,
			StrictPriority: config.StrictPriority,
			Queues:         queues,
		},
	)

	serverMux := asynq.NewServeMux()

	return &workServer{
		server: server,
		mux:    serverMux,
	}
}

func (ws *workServer) RegisterHandlers(handlersMap HandlersMap) {
	for key, handler := range handlersMap {
		ws.mux.Handle(key, handler)
	}
}

func (ws *workServer) Run() error {
	return ws.server.Run(ws.mux)
}
