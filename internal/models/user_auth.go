package models

import (
	"regexp"

	"github.com/jmoiron/sqlx"

	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/defaults"
)

// NewUserAuth builds authboss object
func NewUserAuth(rootURL string, dbConn *sqlx.DB) (*authboss.Authboss, error) {
	ab := authboss.New()
	ab.Config.Paths.RootURL = rootURL

	ab.Config.Storage.Server = NewUserStorer(dbConn)

	userSessionStore := NewUserSessionStore()
	ab.Config.Storage.SessionState = userSessionStore.SessionStorer
	ab.Config.Storage.CookieState = userSessionStore.CookieStorer

	ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "templates")
	ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "templates")

	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}

	ab.Config.Modules.RoutesRedirectOnUnauthed = true

	defaults.SetCore(&ab.Config, false, false)
	ab.Config.Core.BodyReader = abBodyReader()

	if err := ab.Init(); err != nil {
		return nil, err
	}

	return ab, nil
}

func abBodyReader() defaults.HTTPBodyReader {
	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]{1,}`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength: 2,
	}

	return defaults.HTTPBodyReader{
		ReadJSON: false,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, nameRule},
			"recover_end": {passwordRule},
		},
		Confirms: map[string][]string{
			"register":    {"password", authboss.ConfirmPrefix + "password"},
			"recover_end": {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": []string{"email", "name", "password"},
		},
	}
}
