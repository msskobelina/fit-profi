package api

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/msskobelina/fit-profi/api/emails"
	"github.com/msskobelina/fit-profi/domains/authorize"
	"github.com/msskobelina/fit-profi/domains/calendar"
	"github.com/msskobelina/fit-profi/domains/integrations"
	"github.com/msskobelina/fit-profi/domains/nutrition"
	"github.com/msskobelina/fit-profi/domains/profiles"
	"github.com/msskobelina/fit-profi/domains/programs"

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
		// users
		&authorize.User{},
		&authorize.UserToken{},
		&authorize.RevokedToken{},

		// profiles
		&profiles.UserProfile{},
		&profiles.CoachProfile{},
		&profiles.CoachAchievement{},
		&profiles.CoachEducation{},

		// programs
		&programs.TrainingProgram{},
		&programs.ProgramDay{},
		&programs.ProgramExercise{},
		&programs.ExerciseProgress{},

		// nutrition
		&nutrition.DiaryEntry{},
		&nutrition.DiaryItem{},

		// calendar & integrations
		&integrations.UserIntegration{},
	); err != nil {
		l.Fatal("automigration failed", "err", err)
	}

	emailApi := emails.New()
	authRepo := authorize.NewRepository(sql)
	authService := authorize.NewService(
		authRepo,
		emailApi,
		os.Getenv("HMAC_SECRET"),
		os.Getenv("ADMIN_USER_FULLNAME"),
		os.Getenv("ADMIN_USER_EMAIL"),
	)
	authMW := AuthMiddleware(authService)

	profilesRepo := profiles.NewRepository(sql)
	profilesService := profiles.NewService(profilesRepo)

	programsRepo := programs.NewRepository(sql)
	programsService := programs.NewService(programsRepo)

	nutritionRepo := nutrition.NewRepository(sql)
	nutritionService := nutrition.NewService(nutritionRepo)

	integRepo := integrations.NewRepository(sql)
	integSvc := integrations.NewService(integRepo)

	calRepo := calendar.NewRepository(sql)
	calSvc := calendar.NewService(calRepo, integSvc)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet, http.MethodPost, http.MethodPut,
			http.MethodPatch, http.MethodDelete, http.MethodOptions,
		},
		AllowHeaders: []string{"*"},
	}))

	// health
	e.GET("/ping", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	v1 := e.Group("/api/v1")

	authorize.NewHandler(authService).Register(v1, authMW)
	profiles.NewHandler(profilesService).Register(v1, authMW)
	programs.NewHandler(programsService).Register(v1, authMW)
	nutrition.NewHandler(nutritionService).Register(v1, authMW)
	integrations.NewHandler(integSvc).Register(v1, authMW)
	calendar.NewHandler(calSvc).Register(v1, authMW)

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
		if err != nil {
			l.Error("app - Run - httpServer.Notify", "err", err)
		}
	}
	// print
	if err = httpServer.Shutdown(); err != nil {
		l.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}
