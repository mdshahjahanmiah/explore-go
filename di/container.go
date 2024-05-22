package di

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	bsHttp "github.com/mdshahjahanmiah/explore-go/http"
	"go.uber.org/dig"
)

type StartCloser interface {
	Start() error
	Close()
}

type Container struct {
	*dig.Container
}

func New() Container {
	c := Container{dig.New()}

	c.Provide(func() (context.Context, <-chan struct{}) {
		interrupt := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			// interrupt signal sent from terminal
			signal.Notify(sigint, os.Interrupt)
			// sigterm signal sent from kubernetes
			signal.Notify(sigint, syscall.SIGTERM)
			<-sigint

			close(interrupt)
		}()

		return context.TODO(), interrupt
	})

	// register a global abort handler
	c.Provide(func() chan error {
		return make(chan error, 1)
	})

	// register a waitgroup
	c.Provide(func() *sync.WaitGroup {
		return &sync.WaitGroup{}
	})

	return c
}

func (c Container) Provide(constructor interface{}, opts ...dig.ProvideOption) {
	err := c.Container.Provide(constructor, opts...)
	if err != nil {
		slog.Error("error in dig container provide", "err", err)
		c.Close()
		os.Exit(1)
	}
}

func (c Container) Invoke(function interface{}, opts ...dig.InvokeOption) {
	err := c.Container.Invoke(function, opts...)
	if err != nil {
		slog.Error("error in dig container invoke", "err", err)
		c.Close()
		os.Exit(1)
	}
}

func (c Container) ProvideMonitoringEndpoints(group string) {
	c.Container.Provide(createEndpoint("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		slog.Debug("service is healthy")
	})), dig.Group(group))
}

func createEndpoint(pattern string, handler http.Handler) func() bsHttp.Endpoint {
	return func() bsHttp.Endpoint {
		r := chi.NewRouter()
		r.Handle(pattern, handler)
		return bsHttp.Endpoint{Pattern: pattern, Handler: r}
	}
}

// The Start method is responsible for initiating the functionality and returns an error if any error occurs during the start process.
func (c Container) Start() error {
	return c.Container.Invoke(func(in struct {
		dig.In
		List      []StartCloser `group:"startclose"`
		Interrupt <-chan struct{}
		ErrChan   chan error
	}) {
		slog.Info("di container is starting up")
		for _, sc := range in.List {
			sc.Start()
		}

		select {
		case <-in.Interrupt:
			slog.Info("di container is interrupted and/or unexpected shutdown")
		case err := <-in.ErrChan:
			slog.Error("di container chan", "err", err)
		}

		c.Close()
	})
}

// The Close method is responsible for closing and cleaning up any resources used by the type.
func (c Container) Close() error {
	return c.Container.Invoke(func(in struct {
		dig.In
		List []StartCloser `group:"startclose"`
	}) {
		slog.Info("di container cleanup is done")
		for _, sc := range in.List {
			sc.Close()
		}
	})
}
