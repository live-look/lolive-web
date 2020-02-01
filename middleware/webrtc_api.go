package middleware

import (
	"camforchat/internal"
	"context"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

var (
	ctxKeyWebrtcAPI = internal.ContextKey("WebrtcAPI")
)

// GetWebrtcAPI returns database connection link
func GetWebrtcAPI(ctx context.Context) (*internal.Webrtc, bool) {
	w, ok := ctx.Value(ctxKeyWebrtcAPI).(*internal.Webrtc)
	return w, ok
}

// WebrtcAPI is middleware for passing Webrtc between requests
func WebrtcAPI(wrtc *internal.Webrtc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyWebrtcAPI, wrtc))
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}
