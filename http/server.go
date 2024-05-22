package http

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

type Endpoint struct {
	Pattern string
	Handler http.Handler
}

type ServerConfig struct {
	HttpAddress    string
	MaxHeaderBytes *int
	ReadTimeout    *time.Duration
	WriteTimeout   *time.Duration
	IdleTimeout    *time.Duration
}

type Server struct {
	httpServer       *http.Server
	httpServerConfig *ServerConfig
}

func NewServer(serverConfig *ServerConfig, ee []Endpoint, mm []func(http.Handler) http.Handler) *Server {
	router := chi.NewRouter()

	// middlewares
	for _, m := range mm {
		router.Use(m)
	}

	for _, e := range ee {
		if len(e.Pattern) == 0 {
			continue
		}

		if e.Pattern[0] != '/' {
			slog.Error("invalid routing pattern, must begin with '/'", "pattern", e.Pattern)
			continue
		}
		router.Handle(e.Pattern, e.Handler)
	}

	httpServer := &http.Server{
		Addr:           serverConfig.HttpAddress,
		Handler:        router,
		MaxHeaderBytes: DefaultInt(serverConfig.MaxHeaderBytes, http.DefaultMaxHeaderBytes),
		ReadTimeout:    DefaultDuration(serverConfig.ReadTimeout, 10*time.Second),
		WriteTimeout:   DefaultDuration(serverConfig.WriteTimeout, 20*time.Second),
		IdleTimeout:    DefaultDuration(serverConfig.IdleTimeout, 360*time.Second),
	}

	return &Server{
		httpServer:       httpServer,
		httpServerConfig: serverConfig,
	}
}

func (s *Server) Start() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server listen and serve", "err", err)
		}
	}()
	slog.Info("http server is started successfully", "addr", s.httpServerConfig.HttpAddress)

	return nil
}

func (s *Server) Close() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
		slog.Error("shutdown error during server stopped", "addr", s.httpServerConfig.HttpAddress, "err", err)
	}
	slog.Info("http server is stopped successfully", "addr", s.httpServerConfig.HttpAddress)
}

func DefaultInt(value *int, defaultValue int) int {
	if value != nil {
		return *value
	}
	return defaultValue
}

func DefaultDuration(value *time.Duration, defaultValue time.Duration) time.Duration {
	if value != nil {
		return *value
	}
	return defaultValue
}
