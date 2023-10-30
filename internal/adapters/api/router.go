package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/zeusito/go-blueprint/internal/adapters/config"
)

// HTTPServer http server
type HTTPServer struct {
	sc     config.ServerConfigurations
	Router *chi.Mux
	server *http.Server
}

func NewHTTPServer(serverConf config.ServerConfigurations) *HTTPServer {
	router := chi.NewRouter()

	// APM middleware
	//router.Use(apmchiv5.Middleware())

	// A good base middleware stack
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request models (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	return &HTTPServer{
		sc:     serverConf,
		Router: router,
	}
}

func (r *HTTPServer) Start(log *zap.SugaredLogger) {
	// Listening address
	listeningAddr := ":" + r.sc.Port
	log.Infof("Server listening on port %s", listeningAddr)

	// Customizing the server
	r.server = &http.Server{
		Addr:         listeningAddr,
		Handler:      r.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server
	if err := r.server.ListenAndServe(); err != nil {
		log.Errorf("Server error: %v", err)
	}
}

func (r *HTTPServer) Shutdown(ctx context.Context) {
	_ = r.server.Shutdown(ctx)
}
