package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"context"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server App structure
type App struct{}

func AppNew() *App {
	return &App{}
}

func (app *App) Start() error {
	router := chi.NewRouter()
	router.Use(
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	router.Handle("/metrics", promhttp.Handler())
	router.Get("/", HomePage)
	router.Get("/broadcasts/new", BroadcastNew)

	// TODO: extract to internal
	// Serve static assets
	// serves files from web/static dir
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	staticPrefix := "/assets/"
	staticDir := path.Join(cwd, "web", staticPrefix)
	router.Method(
		"GET",
		staticPrefix+"*",
		http.StripPrefix(staticPrefix, http.FileServer(http.Dir(staticDir))),
	)

	// Favicon
	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		if err := serveStaticFile(staticDir+"/favicon.ico", w); err != nil {
			log.Println(err)
		}
	})

	// Handle 404
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		if err := serveStaticFile(staticDir+"/404.html", w); err != nil {
			log.Println(err)
		}
	})

	server := &http.Server{
		Addr:         "0.0.0.0:3333",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err = server.ListenAndServeTLS("configs/cert/localhost.crt", "configs/cert/localhost.key")
	if err != nil {
		return err
	}

	<-serverCtx.Done()

	return nil
}

func serveStaticFile(filePath string, w io.Writer) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	buf := make([]byte, 4*1024) // 4Kb
	if _, err = io.CopyBuffer(w, f, buf); err != nil {
		return err
	}

	return nil
}
