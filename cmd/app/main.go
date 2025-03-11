package main

import (
	"context"
	"io"
	stdLog "log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/AleksandrVishniakov/jwt-auth/internal/configs"
	"github.com/AleksandrVishniakov/jwt-auth/internal/handlers"
	"github.com/AleksandrVishniakov/jwt-auth/internal/repository"
	"github.com/AleksandrVishniakov/jwt-auth/internal/repository/db"
	"github.com/AleksandrVishniakov/jwt-auth/internal/roles"
	"github.com/AleksandrVishniakov/jwt-auth/internal/servers/httpserver"
	"github.com/AleksandrVishniakov/jwt-auth/internal/tokenizer"
	"github.com/AleksandrVishniakov/jwt-auth/internal/usecases"
)

const (
	configPath = "./config.yaml"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := configs.MustConfig()
	stdLog.Printf("running on %s environment\n", cfg.Env)

	log := logger(os.Stdout, cfg.Env)

	if err := run(ctx, log, &cfg); err != nil {
		stdLog.Fatalf("Server stopping: %s\n", err.Error())
	}
}

func run(
	ctx context.Context,
	log *slog.Logger,
	cfg *configs.Config,
) error {
	rolesList := configs.MustParseRoles(configPath)

	database, err := repository.NewPostgresDB(&repository.DBConfigs{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.User,
		DBName:   cfg.DB.DBName,
		Password: cfg.DB.Password,
		SSLMode:  "disable",
	})
	if err != nil {
		return err
	}

	queries := db.New(database)
	repo := repository.New(log, database, queries)
	tokenGenerator := tokenizer.New([]byte(cfg.JWTSignature), time.Hour)
	roleManager := roles.NewManager(log, repo)

	for alias, role := range rolesList {
		if err := roleManager.CreateRole(ctx, alias, role.Permissions, role.Default, role.Super); err != nil {
			return err
		}
	}

	usecase := usecases.New(log, repo, tokenGenerator)

	err = usecase.CreateSuperUser(ctx, cfg.Admin.Login, cfg.Admin.Password)
	if err != nil {
		return err
	}

	handler := handlers.New(log, usecase, tokenGenerator)

	server := httpserver.NewHTTPServer(ctx, cfg.HTTP.Port, handler.InitRoutes())
	defer server.Shutdown(ctx)

	log.Info("running http server", slog.Int("port", cfg.HTTP.Port))
	if err := server.Run(); err != nil {
		return err
	}

	return nil
}

func logger(w io.Writer, env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "production":
		log = slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	slog.SetDefault(log)
	return log
}


