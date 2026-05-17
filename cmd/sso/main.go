package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lfssxxx/auth_service/internal/app"
	"github.com/lfssxxx/auth_service/internal/config"
	"github.com/lfssxxx/auth_service/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	log := setupLogger(cfg.Env)
	log.Info("Application start")

	pool, err := storage.NewPool(ctx, storage.NewConfigMust())
	if err != nil {
		log.Error("pgx pool was not created", "err", err)
		os.Exit(1)
	}
	// TODO: App Init
	application := app.New(log, cfg.Port, pool, cfg.TokenTTL)
	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()

	log.Info("application stopped")
	// TODO: Start gRPC Server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
