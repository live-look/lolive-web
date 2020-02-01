package camforchat

import (
	"context"
	"go.uber.org/zap"
)

var (
	CtxKeyLogger = ContextKey("Logger")
)

// NewLogger build new application logger
func NewLogger(env AppEnv) (*zap.Logger, error) {
	if env == AppEnvProduction {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

// GetLog returns logger from context
func GetLogger(ctx context.Context) (*zap.Logger, bool) {
	l, ok := ctx.Value(CtxKeyLogger).(*zap.Logger)
	return l, ok
}
