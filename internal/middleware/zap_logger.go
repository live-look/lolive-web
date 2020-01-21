package middleware

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	ctxKeyLogger = contextKey("Logger")
)

// GetLog returns logger from context
func GetLog(ctx context.Context) (*zap.Logger, bool) {
	l, ok := ctx.Value(ctxKeyLogger).(*zap.Logger)
	return l, ok
}

// ZapLogger is middleware for passing logger between requests
func ZapLogger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				l.Info(r.Method+" "+r.URL.Path,
					zap.String("ip", r.RemoteAddr),
					zap.String("ua", r.UserAgent()),
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Duration("lat", time.Since(start)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
				)
				//zap.String("reqId", middleware.GetReqID(r.Context())))
			}()
			r = r.WithContext(context.WithValue(r.Context(), ctxKeyLogger, l))
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(h)
	}
}
