package camforchat

import (
	"encoding/base64"
	"time"

	"github.com/gorilla/sessions"
	"github.com/volatiletech/authboss-clientstate"
)

const (
	sessionCookieName = "camforchat"
)

// UserSessionStore represents session store
type UserSessionStore struct {
	CookieStorer  abclientstate.CookieStorer
	SessionStorer abclientstate.SessionStorer
}

// NewUserSessionStore builds UserSessionStore
func NewUserSessionStore() *UserSessionStore {
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)

	cookieStore := abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false

	sessionStore := abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((30 * 24 * time.Hour) / time.Second))

	return &UserSessionStore{CookieStorer: cookieStore, SessionStorer: sessionStore}
}
