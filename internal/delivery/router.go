package delivery

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	ctrlAuthorize "github.com/msskobelina/fit-profi/internal/delivery/controller/authorize"
	ctrlCalendar "github.com/msskobelina/fit-profi/internal/delivery/controller/calendar"
	ctrlIntegrations "github.com/msskobelina/fit-profi/internal/delivery/controller/integrations"
	ctrlNutrition "github.com/msskobelina/fit-profi/internal/delivery/controller/nutrition"
	ctrlProfiles "github.com/msskobelina/fit-profi/internal/delivery/controller/profiles"
	ctrlPrograms "github.com/msskobelina/fit-profi/internal/delivery/controller/programs"
)

type App interface {
	// authorize
	ctrlAuthorize.RegisterHandler
	ctrlAuthorize.LoginHandler
	ctrlAuthorize.LogoutHandler
	ctrlAuthorize.SendEmailHandler
	ctrlAuthorize.ResetPasswordHandler
	// profiles
	ctrlProfiles.CreateUserProfileHandler
	ctrlProfiles.UpdateUserProfileHandler
	ctrlProfiles.GetUserProfileHandler
	ctrlProfiles.CreateCoachProfileHandler
	ctrlProfiles.UpdateCoachProfileHandler
	ctrlProfiles.GetCoachProfileHandler
	// programs
	ctrlPrograms.CreateProgramHandler
	ctrlPrograms.GetProgramHandler
	ctrlPrograms.ListProgramsHandler
	ctrlPrograms.DeleteProgramHandler
	ctrlPrograms.TrackProgressHandler
	// nutrition
	ctrlNutrition.CreateEntryHandler
	ctrlNutrition.ListEntriesHandler
	ctrlNutrition.GetEntryHandler
	ctrlNutrition.UpdateEntryHandler
	ctrlNutrition.DeleteEntryHandler
	// integrations
	ctrlIntegrations.ConnectGoogleHandler
	ctrlIntegrations.ExchangeCallbackHandler
	// calendar
	ctrlCalendar.ListCalendarsHandler
	ctrlCalendar.CreateEventHandler
}

func wrap(h http.Handler, params ...string) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		for _, p := range params {
			ctx = context.WithValue(ctx, controller.PathParamKey(p), c.Param(p))
		}
		h.ServeHTTP(c.Response().Writer, c.Request().WithContext(ctx))
		return nil
	}
}

func Register(e *echo.Echo, app App, io controller.IO, authMW echo.MiddlewareFunc) {
	v1 := e.Group("/api/v1")

	// users (public)
	users := v1.Group("/users")
	users.POST("/register", wrap(ctrlAuthorize.RegisterController(io, app)))
	users.POST("/login", wrap(ctrlAuthorize.LoginController(io, app)))
	users.POST("/send-email", wrap(ctrlAuthorize.SendEmailController(io, app)))
	users.PATCH("/reset-password", wrap(ctrlAuthorize.ResetPasswordController(io, app)))

	// users (private)
	usersPriv := v1.Group("/users", authMW)
	usersPriv.POST("/logout", wrap(ctrlAuthorize.LogoutController(io, app)))
	usersPriv.GET("/check", wrap(ctrlAuthorize.CheckController(io)))

	// profiles
	prof := v1.Group("/profiles", authMW)
	prof.POST("/user", wrap(ctrlProfiles.CreateUserProfileController(io, app)))
	prof.GET("/user", wrap(ctrlProfiles.GetUserProfileController(io, app)))
	prof.PUT("/user", wrap(ctrlProfiles.UpdateUserProfileController(io, app)))
	prof.POST("/coach", wrap(ctrlProfiles.CreateCoachProfileController(io, app)))
	prof.GET("/coach", wrap(ctrlProfiles.GetCoachProfileController(io, app)))
	prof.PUT("/coach", wrap(ctrlProfiles.UpdateCoachProfileController(io, app)))

	// programs
	prog := v1.Group("/programs", authMW)
	prog.POST("", wrap(ctrlPrograms.CreateProgramController(io, app)))
	prog.GET("/:id", wrap(ctrlPrograms.GetProgramController(io, app), "id"))
	prog.GET("", wrap(ctrlPrograms.ListProgramsController(io, app)))
	prog.DELETE("/:id", wrap(ctrlPrograms.DeleteProgramController(io, app), "id"))
	prog.POST("/progress", wrap(ctrlPrograms.TrackProgressController(io, app)))

	// nutrition
	nutr := v1.Group("/nutrition", authMW)
	nutr.POST("/entries", wrap(ctrlNutrition.CreateEntryController(io, app)))
	nutr.GET("/entries", wrap(ctrlNutrition.ListEntriesController(io, app)))
	nutr.GET("/entries/:id", wrap(ctrlNutrition.GetEntryController(io, app), "id"))
	nutr.PUT("/entries/:id", wrap(ctrlNutrition.UpdateEntryController(io, app), "id"))
	nutr.DELETE("/entries/:id", wrap(ctrlNutrition.DeleteEntryController(io, app), "id"))

	// integrations
	v1.GET("/integrations/google/callback", wrap(ctrlIntegrations.ExchangeCallbackController(io, app)))
	integ := v1.Group("/integrations/google", authMW)
	integ.GET("/connect", wrap(ctrlIntegrations.ConnectGoogleController(io, app)))

	// calendar
	cal := v1.Group("/calendar", authMW)
	cal.GET("/list", wrap(ctrlCalendar.ListCalendarsController(io, app)))
	cal.POST("/me/events", wrap(ctrlCalendar.CreateEventController(io, app)))
}
