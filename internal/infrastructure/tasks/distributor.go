package tasks

import (
	"context"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type TaskDistributor interface {
	DistributeTaskEmbedDocument(ctx context.Context, payload *PayloadEmbedDocument, opt TaskProcessingOption) error
}

type AsynqTaskDistributor struct {
	logger *zap.Logger
	client *asynq.Client
}

func NewAsynqTaskDistributor(opt asynq.RedisClientOpt, logger *zap.Logger) TaskDistributor {
	client := asynq.NewClient(opt)
	return &AsynqTaskDistributor{
		logger: logger,
		client: client,
	}
}
