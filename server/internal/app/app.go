package app

import (
	models "Server_part_finance_control/server/internal/domains"
	grpcApp "Server_part_finance_control/server/internal/app/grpc"
	"Server_part_finance_control/server/internal/repository"
	"log/slog"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(
	log *slog.Logger, 
	port int, 
	appConfig models.App, 
	userRepo *repository.UserRepository,
	) *App{
		grpcapp := grpcApp.New(log, port, appConfig, userRepo)

		return &App{
			GRPCServer: grpcapp,
		}
}
