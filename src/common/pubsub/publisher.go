package pubsub

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

type Publisher struct {
	client asynq.Client
}

func NewPublisher(client asynq.Client) *Publisher {
	return &Publisher{
		client: client,
	}
}

func (p *Publisher) Publish(channelKey string, payload interface{}, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(channelKey, payloadBytes, opts...)
	return p.client.Enqueue(task)
}
