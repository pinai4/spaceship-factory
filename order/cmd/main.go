package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/pinai4/spaceship-factory/order/internal/app"
	"github.com/pinai4/spaceship-factory/order/internal/config"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
)

const configPath = "./deploy/compose/order/.env"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appLogger := logger.New(cfg.Logger.Level(), cfg.Logger.AsJSON())
	appCloser := closer.New(appLogger, syscall.SIGINT, syscall.SIGTERM)
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown(appLogger, appCloser)

	a, err := app.New(appCtx, cfg, appCloser, appLogger)
	if err != nil {
		appLogger.Error(appCtx, "❌ Failed to create application", logger.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		appLogger.Error(appCtx, "❌ Error while running application", logger.Error(err))
		return
	}
}

func gracefulShutdown(log *logger.Logger, closer *closer.Closer) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		log.Error(ctx, "❌ Error during shutdown", logger.Error(err))
	}
}
