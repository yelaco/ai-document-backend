package tasks

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/rag"
	"go.uber.org/zap"
)

type Queue string

const (
	QueueCritical Queue = "critical"
	QueueDefault  Queue = "default"
	QueueLow      Queue = "low"
)

type TaskProcessor interface {
	Start() error
	Shutdown()

	ProcessTaskEmbedDocument(ctx context.Context, task *asynq.Task) error
}

type TaskProcessingOption struct {
	MaxRetry  int
	ProcessIn int64
	Queue     Queue
}

type AsynqProcessor struct {
	server       *asynq.Server
	logger       *zap.Logger
	documentRepo interfaces.DocumentRepository
	ragEmbedder  rag.Embedder
	ragStore     rag.Store
}

func NewAsynqProcessor(opt asynq.RedisClientOpt, logger *zap.Logger, documentRepo interfaces.DocumentRepository) TaskProcessor {
	server := asynq.NewServer(opt, asynq.Config{
		Queues: map[string]int{
			string(QueueCritical): 6,
			string(QueueDefault):  3,
			string(QueueLow):      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.Error("AsynqProcessor: task processing error",
				zap.String("task_type", task.Type()),
				zap.ByteString("payload", task.Payload()),
				zap.Error(err),
			)
		}),
	})
	embedder := rag.NewChromaEmbedder()
	return &AsynqProcessor{
		server:       server,
		logger:       logger,
		documentRepo: documentRepo,
		ragEmbedder:  embedder,
	}
}

func (processor *AsynqProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmbedDocument, processor.ProcessTaskEmbedDocument)

	processor.logger.Info("AsynqProcessor: starting task processing server")
	if err := processor.server.Start(mux); err != nil {
		processor.logger.Fatal("AsynqProcessor: failed to start task processing server", zap.Error(err))
		return err
	}
	return nil
}

func (processor *AsynqProcessor) Shutdown() {
	processor.logger.Info("AsynqProcessor: shutting down task processor")
	processor.server.Shutdown()
}
