package middleware

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"net/http"
)

var (
	ctxKeyDb = contextKey("Db")
)

// GetDb return database connection link
func GetDb(ctx context.Context) (*sqlx.DB, bool) {
	l, ok := ctx.Value(ctxKeyDb).(*sqlx.DB)
	return l, ok
}

// NewDb creates new db connection link
func NewDb(spec string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", spec)
}

// Db is middleware for passing db link between requests
func Db(c *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyDb, c))
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}