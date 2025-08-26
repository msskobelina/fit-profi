package controller

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/msskobelina/fit-profi/internal/service"
)

func authMiddleware(s service.Services) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := c.Request().Header.Get("Authorization")
			parts := strings.SplitN(raw, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
				return errorResponse(c, http.StatusUnauthorized, "missing or malformed Authorization header")
			}

			ok, uid, role := s.Users.VerifyAccessToken(c.Request().Context(), parts[1])
			if !ok {
				return errorResponse(c, http.StatusUnauthorized, "invalid token")
			}

			c.Set("userID", uid)
			c.Set("userRole", role)
			return next(c)
		}
	}
}
