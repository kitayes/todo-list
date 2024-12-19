package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo/internal/application"
	delivery "todo/internal/delivery/http"
	"todo/internal/repository"
	"todo/pkg/config"
	service "todo/pkg/services"
)

type Config struct {
	Repo repository.Config `envPrefix:"REPO_"`
	Http delivery.Config   `envPrefix:"HTTP_"`
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
	if err := godotenv.Load(); err != nil {
		logrus.Errorf("error loading env variables: %s", err.Error())
		return
	}
	logrus.SetFormatter(new(logrus.JSONFormatter))
	cfg := Config{}
	if err := config.ReadEnvConfig(&cfg); err != nil {
		logrus.Errorf("error initializing configs: %s", err.Error())
		return
	}

	repos := repository.NewRepository(&cfg.Repo)
	services := application.NewService(repos)
	handlers := delivery.NewHandler(services, &cfg.Http)

	srv := service.NewServiceManager()
	srv.AddService(
		repos,
		services,
		handlers,
	)
	ctx := context.Background()
	go func() {
		if err := srv.Run(ctx); err != nil {
			logrus.Errorf("error occured while running http server: %s", err.Error())
			return
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
		return
	}

}
