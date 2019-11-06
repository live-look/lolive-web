package middleware

import (
	"camforchat/internal/models"
	"context"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

var (
	ctxKeyWebrtcAPI = contextKey("WebrtcAPI")
)

// GetWebrtcAPI returns database connection link
func GetWebrtcAPI(ctx context.Context) (*models.Webrtc, bool) {
	w, ok := ctx.Value(ctxKeyWebrtcAPI).(*models.Webrtc)
	return w, ok
}

// WebrtcAPI is middleware for passing Webrtc between requests
func WebrtcAPI(wrtc *models.Webrtc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyWebrtcAPI, wrtc))
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}
