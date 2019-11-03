package middleware

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/volatiletech/authboss"
	"net/http"
)

var (
	ctxKeyAuthBoss = contextKey("Authboss")
)

// GetAuthBoss returns Authboss object from context
func GetAuthBoss(ctx context.Context) (*authboss.Authboss, bool) {
	u, ok := ctx.Value(ctxKeyAuthBoss).(*authboss.Authboss)

	return u, ok
}

// AuthBoss is middleware for store Authboss between requests
func AuthBoss(ab *authboss.Authboss) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyAuthBoss, ab))
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}
