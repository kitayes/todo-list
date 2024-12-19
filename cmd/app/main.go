package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo/internal/application"
	delivery "todo/internal/delivery/http"
	"todo/internal/repository"
	"todo/pkg/config"
	"todo/pkg/logger"
	service "todo/pkg/services"
)

// TODO: дописать сваггер на все эндпоинты
// TODO: почитать про миграции, добавить в pkg реализацию обновления миграций
type Config struct {
	Repo   repository.Config `envPrefix:"REPO_"`
	logger logger.Config     `envPrefix:"LOGGER_"`
	Http   delivery.Config   `envPrefix:"HTTP_"`
}

// @title           TODO list
// @version         1.0
// @description     Пет проект, заметки.
// @termsOfService  http://swagger.io/terms/

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// TODO: посмотреть что делает метод  godotenv.Load() на примере
	if err := godotenv.Load(); err != nil {
		slog.Error("error loading env variables: %s", err.Error())
		return
	}

	cfg := Config{}
	if err := config.ReadEnvConfig(&cfg); err != nil {
		slog.Error("error initializing configs: %s", err.Error())
		return
	}

	log := logger.NewLogger(&cfg.logger)

	repos := repository.NewRepository(&cfg.Repo, log)
	services := application.NewService(repos, log)
	handlers := delivery.NewHandler(services, &cfg.Http, log)

	srv := service.NewServiceManager()
	srv.AddService(
		repos,
		services,
		handlers,
	)
	// TODO: прочитать про контекст, где используют, виды контекста
	ctx := context.Background()
	go func() {
		if err := srv.Run(ctx); err != nil {
			log.Error("error occured while running http server: %s", err.Error())
			return
		}
	}()

	log.Info("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("TodoApp Shutting Down")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Error("error occured on server shutting down: %s", err.Error())
		return
	}

}
