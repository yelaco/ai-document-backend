package logger

import (
	"os"

	"github.com/yelaco/ai-document-backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(env config.AppEnv) *zap.Logger {
	if env == config.DevEnv {
		cfg := zap.NewDevelopmentEncoderConfig()
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(cfg),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		))
	}
	l, _ := zap.NewProduction()
	return l
}
