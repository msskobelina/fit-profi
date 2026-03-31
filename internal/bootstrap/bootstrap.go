package bootstrap

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googlecalendar "google.golang.org/api/calendar/v3"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	cmdCalendar "github.com/msskobelina/fit-profi/internal/application/command/calendar"
	cmdIntegrations "github.com/msskobelina/fit-profi/internal/application/command/integrations"
	cmdNutrition "github.com/msskobelina/fit-profi/internal/application/command/nutrition"
	cmdProfiles "github.com/msskobelina/fit-profi/internal/application/command/profiles"
	cmdPrograms "github.com/msskobelina/fit-profi/internal/application/command/programs"
	qryAuthorize "github.com/msskobelina/fit-profi/internal/application/query/authorize"
	qryCalendar "github.com/msskobelina/fit-profi/internal/application/query/calendar"
	qryNutrition "github.com/msskobelina/fit-profi/internal/application/query/nutrition"
	qryProfiles "github.com/msskobelina/fit-profi/internal/application/query/profiles"
	qryPrograms "github.com/msskobelina/fit-profi/internal/application/query/programs"
	"github.com/msskobelina/fit-profi/internal/delivery"
	"github.com/msskobelina/fit-profi/internal/delivery/boundary"
	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/infrastructure/email"
	repoAuthorize "github.com/msskobelina/fit-profi/internal/infrastructure/repository/authorize"
	repoCalendar "github.com/msskobelina/fit-profi/internal/infrastructure/repository/calendar"
	repoIntegrations "github.com/msskobelina/fit-profi/internal/infrastructure/repository/integrations"
	repoNutrition "github.com/msskobelina/fit-profi/internal/infrastructure/repository/nutrition"
	repoProfiles "github.com/msskobelina/fit-profi/internal/infrastructure/repository/profiles"
	repoPrograms "github.com/msskobelina/fit-profi/internal/infrastructure/repository/programs"
	"github.com/msskobelina/fit-profi/pkg/analytics"
	"github.com/msskobelina/fit-profi/pkg/httpserver"
	"github.com/msskobelina/fit-profi/pkg/logger"
	metricPkg "github.com/msskobelina/fit-profi/pkg/metric"
	metricEntity "github.com/msskobelina/fit-profi/pkg/metric/entity"
	"github.com/msskobelina/fit-profi/pkg/mysql"
	promMetric "github.com/msskobelina/fit-profi/pkg/observability/prometheus"
)

func Run() {
	l := logger.New(os.Getenv("LOG_LEVEL"))

	// metrics
	promAdapter := promMetric.NewAdapter()
	metricAdapter := promMetric.NewMetricAdapter(promAdapter)
	metricService := metricPkg.NewService(
		metricAdapter,
		metricEntity.UserCreated{},
		metricEntity.LoginFailed{},
	)

	// mysql
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
		&model.User{},
		&model.UserToken{},
		&model.RevokedToken{},
		&model.UserProfile{},
		&model.CoachProfile{},
		&model.CoachAchievement{},
		&model.CoachEducation{},
		&model.TrainingProgram{},
		&model.ProgramDay{},
		&model.ProgramExercise{},
		&model.ExerciseProgress{},
		&model.DiaryEntry{},
		&model.DiaryItem{},
		&model.UserIntegration{},
	); err != nil {
		l.Fatal("automigration failed", "err", err)
	}

	// infrastructure
	emailSender := email.NewSender()
	mixpanel := analytics.NewMixpanel(
		os.Getenv("MIXPANEL_TOKEN"),
		os.Getenv("MIXPANEL_API_HOST"),
	)

	hmacSecret := os.Getenv("HMAC_SECRET")
	adminName := os.Getenv("ADMIN_USER_FULLNAME")
	adminEmail := os.Getenv("ADMIN_USER_EMAIL")

	oauthCfg := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{googlecalendar.CalendarScope},
		Endpoint:     google.Endpoint,
	}

	// repositories
	usersRepo := repoAuthorize.NewRepository(sql)
	profilesRepo := repoProfiles.NewRepository(sql)
	programsRepo := repoPrograms.NewRepository(sql)
	nutritionRepo := repoNutrition.NewRepository(sql)
	integrationsRepo := repoIntegrations.NewRepository(sql)
	_ = repoCalendar.NewRepository(sql)

	// application
	app := &application{
		// authorize
		registerUser:   cmdAuthorize.NewRegisterUserService(usersRepo, mixpanel, metricService, hmacSecret, adminName, adminEmail),
		loginUser:      cmdAuthorize.NewLoginUserService(usersRepo, mixpanel, metricService, hmacSecret, adminName, adminEmail),
		logoutUser:     cmdAuthorize.NewLogoutUserService(usersRepo, hmacSecret),
		sendResetEmail: cmdAuthorize.NewSendResetEmailService(usersRepo, emailSender, hmacSecret, adminName, adminEmail),
		resetPassword:  cmdAuthorize.NewResetPasswordService(usersRepo, hmacSecret),
		verifyToken:    qryAuthorize.NewVerifyTokenService(usersRepo, hmacSecret),
		// profiles
		createUserProfile:  cmdProfiles.NewCreateUserProfileService(profilesRepo),
		updateUserProfile:  cmdProfiles.NewUpdateUserProfileService(profilesRepo),
		getUserProfile:     qryProfiles.NewGetUserProfileService(profilesRepo),
		createCoachProfile: cmdProfiles.NewCreateCoachProfileService(profilesRepo),
		updateCoachProfile: cmdProfiles.NewUpdateCoachProfileService(profilesRepo),
		getCoachProfile:    qryProfiles.NewGetCoachProfileService(profilesRepo),
		// programs
		createProgram: cmdPrograms.NewCreateProgramService(programsRepo),
		deleteProgram: cmdPrograms.NewDeleteProgramService(programsRepo),
		trackProgress: cmdPrograms.NewTrackProgressService(programsRepo),
		getProgram:    qryPrograms.NewGetProgramService(programsRepo),
		listPrograms:  qryPrograms.NewListProgramsService(programsRepo),
		// nutrition
		createEntry: cmdNutrition.NewCreateEntryService(nutritionRepo),
		updateEntry: cmdNutrition.NewUpdateEntryService(nutritionRepo),
		deleteEntry: cmdNutrition.NewDeleteEntryService(nutritionRepo),
		listEntries: qryNutrition.NewListEntriesService(nutritionRepo),
		getEntry:    qryNutrition.NewGetEntryService(nutritionRepo),
		// integrations
		connectGoogle:    cmdIntegrations.NewConnectGoogleService(oauthCfg, hmacSecret),
		exchangeCallback: cmdIntegrations.NewExchangeCallbackService(integrationsRepo, oauthCfg, hmacSecret),
		// calendar
		listCalendars: qryCalendar.NewListCalendarsService(integrationsRepo, oauthCfg),
		createEvent:   cmdCalendar.NewCreateEventService(integrationsRepo, oauthCfg),
	}

	// delivery
	io := boundary.New()
	authMW := delivery.AuthMiddleware(app)

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

	e.GET("/ping", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	e.GET("/metrics", echo.WrapHandler(promAdapter.Handler()))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	delivery.Register(e, app, io, authMW)

	// server
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
		l.Info("bootstrap - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		if err != nil {
			l.Error("bootstrap - Run - httpServer.Notify", "err", err)
		}
	}
	if err = httpServer.Shutdown(); err != nil {
		l.Error("bootstrap - Run - httpServer.Shutdown", "err", err)
	}
}
