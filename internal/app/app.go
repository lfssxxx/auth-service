package app

import (
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	grpcapp "github.com/lfssxxx/auth_service/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, pgxpool *pgxpool.Pool, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
