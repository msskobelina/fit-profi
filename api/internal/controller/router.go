package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/msskobelina/fit-profi/internal/service"
	"github.com/msskobelina/fit-profi/pkg/logger"
)

func NewRouter(e *echo.Echo, s service.Services, l logger.Interface, _ service.Repositories) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodOptions},
		AllowHeaders: []string{"*"},
	}))

	e.GET("/ping", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	h := e.Group("/api/v1")
	newUserRoutes(h, s, l)
}
