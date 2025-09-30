package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/msskobelina/fit-profi/domains/authorize"
	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
)

func AuthMiddleware(s authorize.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := c.Request().Header.Get("Authorization")
			parts := strings.SplitN(raw, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
				return c.JSON(http.StatusUnauthorized, utilsErrors.Error{Message: "missing or malformed Authorization header"})
			}
			ok, uid, role := s.VerifyAccessToken(c.Request().Context(), parts[1])
			if !ok {
				return c.JSON(http.StatusUnauthorized, utilsErrors.Error{Message: "invalid token"})
			}
			c.Set("userID", uid)
			c.Set("userRole", role)
			return next(c)
		}
	}
}

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if role, _ := c.Get("userRole").(string); role != "admin" {
			return c.JSON(http.StatusForbidden, utilsErrors.Error{Message: "admin only"})
		}
		return next(c)
	}
}
