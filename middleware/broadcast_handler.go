package middleware

import (
	"camforchat/internal"
	"context"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

var (
	ctxKeyBroadcastHandler = internal.ContextKey("BroadcastHandler")
)

// GetBroadcastHandler returns BroadcastHandler from context
func GetBroadcastHandler(ctx context.Context) (*internal.BroadcastHandler, bool) {
	u, ok := ctx.Value(ctxKeyBroadcastHandler).(*internal.BroadcastHandler)

	return u, ok
}

// BroadcastHandler is middleware for passing BroadcastHandler between requests
func BroadcastHandler(bh *internal.BroadcastHandler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyBroadcastHandler, bh))
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}
