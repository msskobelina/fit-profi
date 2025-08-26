package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/msskobelina/fit-profi/internal/api/emails"
	"github.com/msskobelina/fit-profi/internal/controller"
	"github.com/msskobelina/fit-profi/internal/domain"
	"github.com/msskobelina/fit-profi/internal/repository"
	"github.com/msskobelina/fit-profi/internal/service"
	"github.com/msskobelina/fit-profi/pkg/httpserver"
	"github.com/msskobelina/fit-profi/pkg/logger"
	"github.com/msskobelina/fit-profi/pkg/mysql"
)

func Run() {
	l := logger.New(os.Getenv("LOG_LEVEL"))

	sql, err := mysql.New(mysql.MySQLConfig{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Database: os.Getenv("MYSQL_DATABASE"),
	})
	if err != nil {
		l.Fatal("failed to connect to mysql", "err", err)
	}

	if err = sql.DB.AutoMigrate(
		&domain.User{},
		&domain.UserToken{},
		&domain.RevokedToken{},
	); err != nil {
		l.Fatal("automigration failed", "err", err)
	}

	apis := service.APIs{Emails: emails.New()}
	repos := service.Repositories{Users: repository.NewUsersRepo(sql)}
	services := service.Services{Users: service.NewUserService(repos, apis)}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())

	controller.NewRouter(e, services, l, repos)

	httpServer := httpserver.New(
		e,
		httpserver.Port(os.Getenv("HTTP_PORT")),
		httpserver.ReadTimeout(60*time.Second),
		httpserver.WriteTimeout(60*time.Second),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify", "err", err)
	}

	if err = httpServer.Shutdown(); err != nil {
		l.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}
