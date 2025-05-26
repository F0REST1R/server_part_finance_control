package main

import (
	"Server_part_finance_control/server/internal/app"
	domains "Server_part_finance_control/server/internal/domains"
	PSQL "Server_part_finance_control/server/internal/PostgresSQL" // Исправлен импорт
	"Server_part_finance_control/server/internal/repository"
	"Server_part_finance_control/server/pkg/config"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. Инициализация логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("service", "auth-service")

	// 2. Загрузка конфигурации
	cfg, err := config.New()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	
	// 3. Подключение к PostgreSQL
	pgConfig := PSQL.Config{
		Host:     cfg.POSTGRES.Host,
		Port:     cfg.POSTGRES.Port,
		Username: cfg.POSTGRES.Username,
		Password: cfg.POSTGRES.Password,
		Database: cfg.POSTGRES.Database,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Исправленный вызов конструктора
	pgClient, err := PSQL.NewPost(pgConfig, ctx)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pgClient.Close()

	// 4. Инициализация репозитория
	userRepo := repository.NewUserRepository(pgClient.Conn)

	// 5. Настройка приложения (исправлено согласно новой структуры App)
	appConfig := domains.App{
		ID:     1, // Должно быть числом согласно структуре
		Name:   "auth-service",
		Secret: "your_jwt_secret_key_here",
	}

	// 6. Создание gRPC приложения
	application := app.New(
		logger,
		cfg.GRPC.PORT,
		appConfig,
		userRepo,
	)

	// 7. Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("starting gRPC server", "port", cfg.GRPC.PORT)
		application.GRPCServer.MustRun()
	}()

	<-shutdown
	logger.Info("shutting down server")

	// Остановка gRPC сервера
	application.GRPCServer.Stop()

	logger.Info("server stopped gracefully")
}