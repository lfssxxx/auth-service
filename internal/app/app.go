package app

import (
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	grpcapp "github.com/lfssxxx/auth_service/internal/app/grpc"
	"github.com/lfssxxx/auth_service/internal/services/auth"
	"github.com/lfssxxx/auth_service/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, pgxpool *pgxpool.Pool, tokenTTL time.Duration) *App {

	storage := postgres.New(pgxpool)

	authService := auth.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
