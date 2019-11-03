package main

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/volatiletech/authboss/auth"
	_ "github.com/volatiletech/authboss/confirm"
	_ "github.com/volatiletech/authboss/logout"
	_ "github.com/volatiletech/authboss/register"

	"context"
	"fmt"
	"path/filepath"

	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss/confirm"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/valve"

	"camforchat/internal/handlers"
	appMiddleware "camforchat/internal/middleware"
	"camforchat/internal/models"
	"camforchat/internal/usecases"

	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Logging
	logger, err := appMiddleware.NewZapLogger()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	// Databse
	dbSpec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	dbConn, err := appMiddleware.NewDb(dbSpec)
	if err != nil {
		log.Fatalf("Opening db connectoin failed: %v", err)
	}
	defer dbConn.Close()

	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("Db is not available: %v %s", err, dbSpec)
	}

	// Authentication, registration
	ab, err := models.NewUserAuth(os.Getenv("APP_ROOT_URL"), dbConn)
	if err != nil {
		log.Fatalf("Error initialization of authboss: %v", err)
	}

	//m := webrtc.MediaEngine{}
	// m.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))
	//m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	//webrtcApi := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	// Broadcasts
	broadcastHandlerContext, broadcastHandlerCancel := context.WithCancel(context.Background())
	defer broadcastHandlerCancel()

	// Run broadcast Handler
	broadcastHandler := models.NewBroadcastHandler()
	broadcastHandler.Run(broadcastHandlerContext)

	broadcastHanderMiddleware := appMiddleware.BroadcastHandler(broadcastHandler)

	// Web server, routing
	r := chi.NewRouter()
	r.Use(
		middleware.RealIP,
		appMiddleware.ZapLogger(logger),
		appMiddleware.Db(dbConn),
		ab.LoadClientStateMiddleware,
		appMiddleware.CurrentUserDataInject(ab),
		middleware.Recoverer,
	)

	r.Route("/broadcasts", func(r chi.Router) {
		// Require auth
		r.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondRedirect), confirm.Middleware(ab))

		r.Get("/new", handlers.BroadcastsNew)
		r.With(broadcastHanderMiddleware).
			Post("/", handlers.BroadcastsCreate)
		r.Get("/{broadcastId}", handlers.BroadcastsShow)

		r.With(broadcastHanderMiddleware).
			Post("/{broadcastID}/viewers", handlers.ViewersCreate)
	})

	r.Group(func(r chi.Router) {
		r.Use(ab.LoadClientStateMiddleware, authboss.ModuleListMiddleware(ab))
		r.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))
	})

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/", handlers.HomePage)

	// Server static files
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	path := "/static"
	fs := http.StripPrefix(path, http.FileServer(http.Dir(filesDir)))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

	valv := valve.New()
	baseCtx := valv.Context()
	server := &http.Server{
		Addr:         os.Getenv("APP_LISTEN_ADDR"),
		Handler:      chi.ServerBaseContext(baseCtx, r),
		ErrorLog:     log.New(&usecases.FwdToZapWriter{Logger: logger}, "", 0),
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

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(21 * time.Second):
				log.Println("not all connections done")
			case <-ctx.Done():

			}
		}
	}()

	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal("Start server is failed", zap.Error(err))
	}
}
