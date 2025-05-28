package main

import (
	"context"
	_ "github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log/slog"
	"todo/internal/application"
	delivery "todo/internal/delivery/http"
	"todo/internal/repository"
	"todo/pkg/config"
	"todo/pkg/logger"
	service "todo/pkg/services"
)

// TODO: прочитать про unit-test'ы, mock'и, покрыть все тестами
// TODO: прочитать про линтер и добавить линтер в проект
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

	srv := service.NewManager(log)
	srv.AddService(
		repos,
		services,
		handlers,
	)
	// TODO: прочитать про контекст, где используют, виды контекста
	ctx := context.Background()
	if err := srv.Run(ctx); err != nil {
		err := errors.Wrap(err, "s.listRepo.GetById(...) err:")
		log.Error(err.Error())
		return
	}

	log.Info("TodoApp Started")
}
