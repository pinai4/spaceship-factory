package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pinai4/spaceship-factory/order/internal/config"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
	"github.com/pinai4/spaceship-factory/platform/pkg/sqldb/migrator"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

type App struct {
	config      *config.Config
	closer      *closer.Closer
	logger      *logger.Logger
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context, config *config.Config, closer *closer.Closer, logger *logger.Logger) (*App, error) {
	a := &App{config: config, closer: closer, logger: logger}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initDBMigrations,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer(a.config, a.closer)
	return nil
}

func (a *App) initDBMigrations(ctx context.Context) error {
	migratorRunner := migrator.New(a.diContainer.PostgresDB(ctx).DB, a.config.DBMigrations.Path())
	if err := migratorRunner.Up(); err != nil {
		return err
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return err
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	// Manual health check handler
	r.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			a.logger.Error(ctx, "failed to write health response", logger.Error(err))
		}
	})

	a.httpServer = &http.Server{
		Addr:              a.config.HTTPServer.Address(),
		Handler:           r,
		ReadHeaderTimeout: a.config.HTTPServer.ReadTimeout(),
	}
	a.closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	a.logger.Info(ctx, "🚀 HTTP-server listening", logger.String("address", a.config.HTTPServer.Address()))

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
