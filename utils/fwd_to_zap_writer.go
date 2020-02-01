package utils

import (
	"go.uber.org/zap"
)

// FwdToZapWriter is wrapper for pass zap logger to net/http server
type FwdToZapWriter struct {
	Logger *zap.Logger
}

func (fw *FwdToZapWriter) Write(p []byte) (n int, err error) {
	fw.Logger.Error(string(p[:]))
	return len(p), nil
}
