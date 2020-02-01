package middleware

import (
	"camforchat/internal"
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
	"log"
	"net/http"
)

var (
	ctxKeyCurrentUser = internal.ContextKey("CurrentUser")
)

// GetCurrentUser extract current user from context
func GetCurrentUser(ctx context.Context) (*internal.User, bool) {
	user, ok := ctx.Value(ctxKeyCurrentUser).(*internal.User)
	return user, ok
}

// CurrentUserDataInject is middleware for injecting currentUser data
func CurrentUserDataInject(ab *authboss.Authboss) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			var (
				currentUserID   int64
				currentUserName string
				user            *internal.User
			)
			logger, ok := GetLog(r.Context())
			if !ok {
				log.Println("Can't get logger from context")
				http.Error(w, http.StatusText(500), 500)
				return
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			currentUser, err := ab.LoadCurrentUser(&r)
			if err != nil && err != authboss.ErrUserNotFound {
				logger.Error("loading current user failed", zap.Error(err))
				http.Error(w, http.StatusText(500), 500)
				return
			}

			if currentUser != nil {
				user = currentUser.(*internal.User)
				currentUserID = user.ID
				currentUserName = user.Name
			}

			data := authboss.HTMLData{
				"loggedin":          currentUser != nil,
				"current_user_id":   currentUserID,
				"current_user_name": currentUserName,
				//"csrf_token":        nosurf.Token(*r),
				"flash_success": authboss.FlashSuccess(w, r),
				"flash_error":   authboss.FlashError(w, r),
			}

			newCtx := context.WithValue(r.Context(), ctxKeyCurrentUser, user)
			newCtx = context.WithValue(newCtx, authboss.CTXKeyData, data)

			r = r.WithContext(newCtx)
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}
