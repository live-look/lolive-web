package middleware

import (
	"camforchat/internal/models"
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/volatiletech/authboss"
	"net/http"
)

// CurrentUserDataInject is middleware for injecting currentUser data
func CurrentUserDataInject(ab *authboss.Authboss) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			data := layoutData(w, &r, ab)

			r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, data))
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}

func layoutData(w http.ResponseWriter, r **http.Request, ab *authboss.Authboss) authboss.HTMLData {
	var (
		currentUserID   int64
		currentUserName string
	)

	userInter, err := ab.LoadCurrentUser(r)

	if userInter != nil && err == nil {
		user := userInter.(*models.User)
		currentUserName = user.Name
		currentUserID = user.ID
	}

	return authboss.HTMLData{
		"loggedin":          userInter != nil,
		"current_user_id":   currentUserID,
		"current_user_name": currentUserName,
		//"csrf_token":        nosurf.Token(*r),
		"flash_success": authboss.FlashSuccess(w, *r),
		"flash_error":   authboss.FlashError(w, *r),
	}
}
