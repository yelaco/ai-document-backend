package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/dtos"
	"go.uber.org/zap"
)

type SSEHandler struct {
	logger      *zap.Logger
	chatService interfaces.ChatService
}

func NewSSEHandler(logger *zap.Logger, chatService interfaces.ChatService) *SSEHandler {
	return &SSEHandler{
		logger:      logger,
		chatService: chatService,
	}
}

// StreamChat handles SSE streaming for chat responses
func (h *SSEHandler) StreamChat(c *gin.Context) {
	chatID, err := types.NewChatIDFromString(c.Param("id"))
	if err != nil {
		_ = c.Error(fmt.Errorf("SSEHandler.StreamChat: invalid chat ID: %w", err))
		c.JSON(http.StatusBadRequest, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "invalid chat ID",
			},
		})
		return
	}

	// Bind query parameters
	var req dtos.AnswerQuestionParams
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(fmt.Errorf("SSEHandler.StreamChat: failed to bind query: %w", err))
		c.JSON(http.StatusBadRequest, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "invalid request payload",
			},
		})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// Get client context for cancellation
	clientGone := c.Request.Context().Done()

	// Create a channel for streaming events
	eventChan := make(chan dtos.SSEEvent, 10)
	defer close(eventChan)

	// Start the chat processing in a goroutine
	go h.processChatStream(c.Request.Context(), chatID, req, eventChan)

	// Stream events to client
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			h.logger.Info("Client disconnected")
			return false
		case event, ok := <-eventChan:
			if !ok {
				return false
			}

			// Send the event
			if err := sse.Encode(w, sse.Event{
				Event: string(event.Type),
				Data:  event,
			}); err != nil {
				h.logger.Error("Failed to encode SSE event", zap.Error(err))
				return false
			}

			// Check if this is a completion event
			if event.Type == dtos.SSEEventTypeComplete || event.Type == dtos.SSEEventTypeError {
				return false
			}

			return true
		case <-time.After(30 * time.Second):
			// Send heartbeat every 30 seconds
			heartbeat := dtos.SSEEvent{
				Type:      dtos.SSEEventTypeHeartbeat,
				Data:      gin.H{"message": "heartbeat"},
				Timestamp: time.Now(),
			}

			if err := sse.Encode(w, sse.Event{
				Event: string(heartbeat.Type),
				Data:  heartbeat,
			}); err != nil {
				h.logger.Error("Failed to send heartbeat", zap.Error(err))
				return false
			}
			return true
		}
	})
}

// processChatStream simulates chat processing and sends events
func (h *SSEHandler) processChatStream(ctx context.Context, chatID types.ChatID, req dtos.AnswerQuestionParams, eventChan chan<- dtos.SSEEvent) {
	defer func() {
		// Send completion event
		eventChan <- dtos.SSEEvent{
			Type:      dtos.SSEEventTypeComplete,
			Data:      gin.H{"message": "Processing complete"},
			Timestamp: time.Now(),
		}
	}()

	// Send progress events
	steps := []string{"Processing question", "Searching documents", "Generating response", "Finalizing answer"}

	for i, step := range steps {
		select {
		case <-ctx.Done():
			return
		default:
			// Send progress update
			progressEvent := dtos.SSEEvent{
				Type: dtos.SSEEventTypeProgress,
				Data: dtos.ProgressEventData{
					Step:        step,
					Progress:    (i + 1) * 25,
					Message:     fmt.Sprintf("Step %d of %d: %s", i+1, len(steps), step),
					TotalSteps:  len(steps),
					CurrentStep: i + 1,
				},
				Timestamp: time.Now(),
			}

			eventChan <- progressEvent
			time.Sleep(1 * time.Second) // Simulate processing time
		}
	}

	// Simulate streaming chat response
	fullResponse := "This is a simulated AI response that demonstrates how SSE streaming works with your chat application. The response can be sent in chunks to provide a real-time typing effect."

	words := []string{}
	for _, word := range fullResponse {
		if string(word) == " " && len(words) > 0 {
			// Send partial response
			content := ""
			for _, w := range words {
				content += w + " "
			}

			chatEvent := dtos.SSEEvent{
				Type: dtos.SSEEventTypeMessage,
				Data: dtos.ChatStreamResponse{
					ChatID:    chatID,
					Content:   content,
					IsPartial: true,
					Metadata: map[string]interface{}{
						"word_count": len(words),
					},
				},
				Timestamp: time.Now(),
			}

			select {
			case <-ctx.Done():
				return
			case eventChan <- chatEvent:
				time.Sleep(100 * time.Millisecond) // Simulate typing delay
			}
			words = []string{}
		} else {
			words = append(words, string(word))
		}
	}

	// Send final complete message
	finalEvent := dtos.SSEEvent{
		Type: dtos.SSEEventTypeMessage,
		Data: dtos.ChatStreamResponse{
			ChatID:    chatID,
			Content:   fullResponse,
			IsPartial: false,
			Metadata: map[string]interface{}{
				"final":      true,
				"word_count": len(fullResponse),
			},
		},
		Timestamp: time.Now(),
	}

	eventChan <- finalEvent
}

// StreamProgress provides a generic SSE endpoint for progress updates
func (h *SSEHandler) StreamProgress(c *gin.Context) {
	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "task_id parameter is required",
			},
		})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	clientGone := c.Request.Context().Done()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case <-time.After(2 * time.Second):
			// Simulate progress updates
			event := dtos.SSEEvent{
				Type: dtos.SSEEventTypeProgress,
				Data: dtos.ProgressEventData{
					Step:     "Processing",
					Progress: 50,
					Message:  fmt.Sprintf("Processing task %s...", taskID),
				},
				Timestamp: time.Now(),
			}

			if err := sse.Encode(w, sse.Event{
				Event: string(event.Type),
				Data:  event,
			}); err != nil {
				return false
			}
			return true
		}
	})
}
