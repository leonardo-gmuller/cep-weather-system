package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/config"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/gateway/api"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/telemetry"
	"golang.org/x/sync/errgroup"
)

// Injected on build via ldflags.
var (
	BuildTime   = "undefined"
	BuildCommit = "undefined"
	BuildTag    = "undefined"
)

func main() {
	mainCtx := context.Background()

	// Config
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load configurations: %v", err)
	}

	// Open Telemetry
	otel, err := telemetry.NewOtel(mainCtx, cfg.Otel, string(cfg.Environment), BuildTag)
	if err != nil {
		log.Fatalf("failed to start otel: %v", err)
	}

	ctx := telemetry.ContextWithTracer(mainCtx, otel.Tracer)

	// Application
	appl, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	// Server
	server := &http.Server{
		Addr:         cfg.Server.Address,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		Handler:      api.New(cfg, appl.UseCase).Handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful Shutdown
	stopCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	group, groupCtx := errgroup.WithContext(stopCtx)

	//nolint:wrapcheck
	group.Go(func() error {
		log.Printf("starting api server")

		return server.ListenAndServe()
	})
	//nolint:contextcheck
	group.Go(func() error {
		<-groupCtx.Done()

		log.Printf("stopping api; interrupt signal received")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.App.GracefulShutdownTimeout)
		defer cancel()

		var errs error

		if err := server.Shutdown(timeoutCtx); err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to stop server: %w", err))
		}

		return errs
	})

	if err := group.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("api exit reason: %v", err)
	}

	stop()
}
