package grpcApp

import (
	models "Server_part_finance_control/server/internal/domains"
	authgrpc "Server_part_finance_control/server/internal/grpc/auth"
	"Server_part_finance_control/server/internal/repository"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log *slog.Logger
	GRPCServer *grpc.Server
	port int
	userRepo *repository.UserRepository
	appConfig models.App
}

func (a *App) MustRun(){
	if err := a.Run(); err != nil{
		panic(err)
	}
}

func New(log *slog.Logger, port int, appConfig models.App, userRepo *repository.UserRepository) *App{
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, appConfig, userRepo)
	reflection.Register(gRPCServer)

	return &App{
		log: log,
		GRPCServer: gRPCServer,
		port: port,
		userRepo: userRepo,
		appConfig: appConfig,
	}
}  

func (a *App) Run() error {
	const op = "grpc.run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil{
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.GRPCServer.Serve(l); err != nil{
		return fmt.Errorf("%s: %w", op, err)
	}


	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	a.GRPCServer.GracefulStop()
}