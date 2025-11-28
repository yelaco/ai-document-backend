package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

const (
	TypeEmbedDocument = "document:embed"
)

type PayloadEmbedDocument struct {
	DocumentID uuid.UUID `json:"document_id"`
	UserID     uuid.UUID `json:"user_id"`
}

func NewTaskEmbedDocument(documentID uuid.UUID, userID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(PayloadEmbedDocument{
		DocumentID: documentID,
		UserID:     userID,
	})
	if err != nil {
		return nil, fmt.Errorf("NewTaskProcessDocument: failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeEmbedDocument, payload), nil
}

func (distributor *AsynqTaskDistributor) DistributeTaskEmbedDocument(ctx context.Context, payload *PayloadEmbedDocument, opt TaskProcessingOption) error {
	task, err := NewTaskEmbedDocument(payload.DocumentID, payload.UserID)
	if err != nil {
		return fmt.Errorf("AsynqTaskDistributor.DistributeTaskEmbedDocument: failed to create task: %w", err)
	}
	asynqOpts := []asynq.Option{
		asynq.MaxRetry(opt.MaxRetry),
		asynq.ProcessIn(time.Duration(opt.ProcessIn) * time.Second),
		asynq.Queue(string(opt.Queue)),
	}
	info, err := distributor.client.EnqueueContext(ctx, task, asynqOpts...)
	if err != nil {
		return fmt.Errorf("AsynqTaskDistributor.DistributeTaskEmbedDocument: failed to enqueue task: %w", err)
	}
	distributor.logger.Info("AsynqTaskDistributor.DistributeTaskEmbedDocument: task enqueued",
		zap.String("task_id", info.ID),
		zap.String("document_id", payload.DocumentID.String()),
		zap.String("user_id", payload.UserID.String()),
		zap.String("queue", info.Queue),
		zap.Int("max_retry", info.MaxRetry),
	)
	return nil
}

func (processor *AsynqProcessor) ProcessTaskEmbedDocument(ctx context.Context, task *asynq.Task) error {
	var payload PayloadEmbedDocument
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("DocumentProcessor.ProcessTask: failed to unmarshal payload: %w", err)
	}
	document, err := processor.documentRepo.GetDocument(ctx, payload.DocumentID, payload.UserID)
	if err != nil {
		return fmt.Errorf("DocumentProcessor.ProcessTask: failed to get document: %w", err)
	}
	processor.logger.Info("Processing document", zap.String("name", document.OriginalName))
	return nil
}
