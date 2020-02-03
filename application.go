package camforchat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/confirm"
	"github.com/volatiletech/authboss/defaults"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/valve"

	"gitlab.com/isqad/camforchat/utils"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/joho/godotenv"
)

// Application structure
type Application struct {
	env    AppEnv
	logger *zap.Logger
	db     *sqlx.DB
	ab     *authboss.Authboss
}

// Getenv returns current application environment
func (app *Application) getenv() AppEnv {
	if app.env == AppEnv("") {
		osEnv := os.Getenv("APP_ENV")

		if osEnv == "" ||
			AppEnv(osEnv) != AppEnvTest && AppEnv(osEnv) != AppEnvProduction && AppEnv(osEnv) != AppEnvDevelopment {

			app.env = AppEnvDevelopment
		} else {
			app.env = AppEnv(osEnv)
		}
	}

	return app.env
}

func (app *Application) initConfig() {
	err := godotenv.Load(fmt.Sprintf("%s.%s", ".env", app.getenv()))
	if err != nil {
		log.Fatal("Application initialization failed: ", err)
	}
}

func (app *Application) initLog() {
	logger, err := NewLogger(app.getenv())
	if err != nil {
		log.Fatal("Application initialization failed: ", err)
	}

	app.logger = logger
}

func (app *Application) initDB() {
	spec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("CAMFORCHAT_DB_USER"),
		os.Getenv("CAMFORCHAT_DB_PASSWORD"),
		os.Getenv("CAMFORCHAT_DB_HOST"),
		os.Getenv("CAMFORCHAT_DB_NAME"),
	)

	db, err := NewDb(spec)
	if err != nil {
		log.Fatal("Application initialization failed: ", err)
	}

	app.db = db
}

func (app *Application) initAuthboss() {
	ab := authboss.New()
	ab.Config.Paths.RootURL = os.Getenv("CAMFORCHAT_ROOT_URL")

	ab.Config.Storage.Server = NewUserStorer(app.db)

	userSessionStore := NewUserSessionStore()
	ab.Config.Storage.SessionState = userSessionStore.SessionStorer
	ab.Config.Storage.CookieState = userSessionStore.CookieStorer

	ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "templates")
	ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "templates")
	// SetCore(config *authboss.Config, readJSON, useUsername bool) {
	defaults.SetCore(&ab.Config, false, false)

	ab.Config.Core.BodyReader = abBodyReader()
	ab.Config.Core.Mailer = NewSMTPMailer(
		os.Getenv("CAMFORCHAT_MAIL_SERVER"),
		smtp.PlainAuth(
			"",
			os.Getenv("CAMFORCHAT_MAIL_USERNAME"),
			os.Getenv("CAMFORCHAT_MAIL_PASSWORD"),
			os.Getenv("CAMFORCHAT_MAIL_HOSTNAME"),
		),
	)

	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}
	ab.Config.Modules.RoutesRedirectOnUnauthed = true

	ab.Config.Mail.From = os.Getenv("CAMFORCHAT_MAIL_FROM")
	ab.Config.Mail.RootURL = os.Getenv("CAMFORCHAT_ROOT_URL")
	ab.Config.Mail.FromName = os.Getenv("CAMFORCHAT_MAIL_FROM_NAME")

	if err := ab.Init(); err != nil {
		log.Fatal("Application initialization failed: ", err)
	}

	app.ab = ab
}

// Run initializes and runs application
func (app *Application) Run() {
	app.initConfig()
	// Logging
	app.initLog()
	// Database
	app.initDB()
	// Authentication, registration
	app.initAuthboss()

	// Broadcasts
	broadcastHandlerContext, broadcastHandlerCancel := context.WithCancel(context.Background())
	defer broadcastHandlerCancel()

	// Run broadcast Handler
	broadcastHandler := NewBroadcastHandler()
	broadcastHandler.Run(broadcastHandlerContext)

	broadcastHanderMiddleware := BroadcastHandlerMiddleware(broadcastHandler)

	wrtc := NewWebrtc()
	webrtcMiddleware := WebrtcAPIMiddleware(wrtc)

	// Web server, routing
	r := chi.NewRouter()
	r.Use(
		middleware.RealIP,
		LoggerMiddleware(app.logger),
		DbMiddleware(app.db),
		app.ab.LoadClientStateMiddleware,
		CurrentUserMiddleware(app.ab),
		middleware.Recoverer,
		webrtcMiddleware,
	)

	r.Route("/broadcasts", func(r chi.Router) {
		// Require auth
		r.Use(
			app.ab.LoadClientStateMiddleware,
			authboss.Middleware2(app.ab, authboss.RequireNone, authboss.RespondRedirect),
			confirm.Middleware(app.ab),
		)

		r.Get("/new", BroadcastsNew)
		r.With(broadcastHanderMiddleware).
			Post("/", BroadcastsCreate)
		r.Get("/{broadcastID}", BroadcastsShow)

		r.With(broadcastHanderMiddleware).
			Post("/{broadcastID}/viewers", ViewersCreate)
	})

	r.Group(func(r chi.Router) {
		r.Use(app.ab.LoadClientStateMiddleware, authboss.ModuleListMiddleware(app.ab))
		r.Mount("/auth", http.StripPrefix("/auth", app.ab.Config.Core.Router))
	})

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/", HomePage)

	valv := valve.New()
	baseCtx := valv.Context()
	server := &http.Server{
		Addr:         os.Getenv("CAMFORCHAT_LISTEN_ADDR"),
		Handler:      chi.ServerBaseContext(baseCtx, r),
		ErrorLog:     log.New(&utils.FwdToZapWriter{Logger: app.logger}, "", 0),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down..")

			// first valv
			valv.Shutdown(20 * time.Second)

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			// start http shutdown
			server.Shutdown(ctx)
			app.Shutdown()

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(21 * time.Second):
				log.Println("not all connections done")
			case <-ctx.Done():

			}
		}
	}()

	err := server.ListenAndServe()
	if err != nil {
		app.logger.Fatal("Start server is failed", zap.Error(err))
	}
}

// Shutdown stops application gracefully
func (app *Application) Shutdown() {
	app.logger.Sync()
	app.db.Close()
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
