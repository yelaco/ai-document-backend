package dtos

import (
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type AnswerQuestionParams struct {
	Question string `form:"question" binding:"required"`
}

// SSE Event types
type SSEEventType string

const (
	SSEEventTypeMessage   SSEEventType = "message"
	SSEEventTypeError     SSEEventType = "error"
	SSEEventTypeProgress  SSEEventType = "progress"
	SSEEventTypeComplete  SSEEventType = "complete"
	SSEEventTypeHeartbeat SSEEventType = "heartbeat"
)

// SSE Event structure
type SSEEvent struct {
	Type      SSEEventType `json:"type"`
	Data      interface{}  `json:"data"`
	Timestamp time.Time    `json:"timestamp"`
}

// Chat streaming response
type ChatStreamResponse struct {
	ChatID    types.ChatID           `json:"chat_id"`
	Content   string                 `json:"content"`
	IsPartial bool                   `json:"is_partial"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Progress event data
type ProgressEventData struct {
	Step        string `json:"step"`
	Progress    int    `json:"progress"` // 0-100
	Message     string `json:"message"`
	TotalSteps  int    `json:"total_steps,omitempty"`
	CurrentStep int    `json:"current_step,omitempty"`
}

// Error event data
type ErrorEventData struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Details      string `json:"details,omitempty"`
}
