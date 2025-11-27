package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LoggerSetter interface {
	SetLogger(*zap.Logger)
}

func ZapLoggerMiddleware(base *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// generate or use existing request ID
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set("req_id", reqID)

		// create request-scoped logger
		reqLogger := base.With(
			zap.String("req_id", reqID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		// store in Gin context
		c.Set("logger", reqLogger)
		c.Writer.Header().Set("X-Request-ID", reqID)

		c.Next() // process the request

		// log after request
		end := time.Now()
		latency := end.Sub(start)

		status := c.Writer.Status()
		entry := reqLogger.With(
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		if len(c.Errors) > 0 {
			entry = entry.With(zap.String("error", c.Errors.ByType(gin.ErrorTypeAny).String()))
		}

		switch {
		case status >= 500:
			entry.Error("Request failed")
		case status >= 400:
			entry.Warn("Client error")
		default:
			entry.Info("Request completed")
		}
	}
}
