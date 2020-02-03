package camforchat

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"net/http"
)

var (
	ctxKeyDb = ContextKey("Db")
)

// NewDb initialize db connection
func NewDb(spec string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", spec)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetDb return database connection link
func GetDb(ctx context.Context) (*sqlx.DB, bool) {
	db, ok := ctx.Value(ctxKeyDb).(*sqlx.DB)
	return db, ok
}

// DbMiddleware is middleware for passing db link between requests
func DbMiddleware(c *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyDb, c))
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}
