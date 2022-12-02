package handlers

import (
	"context"
	"github.com/hibiken/asynq"
)

type statusProcessor struct {
}

func NewStatusProcessor() *statusProcessor {
	return &statusProcessor{}
}

func (sp *statusProcessor) ProcessTask(_ context.Context, t *asynq.Task) error {
	return nil
}
