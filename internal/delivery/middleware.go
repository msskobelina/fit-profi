package delivery

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	qryAuthorize "github.com/msskobelina/fit-profi/internal/application/query/authorize"
	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
)

type TokenVerifier interface {
	VerifyToken(ctx context.Context, q qryAuthorize.VerifyTokenQuery) (*qryAuthorize.VerifyTokenResult, error)
}

func AuthMiddleware(verifier TokenVerifier) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := c.Request().Header.Get("Authorization")
			parts := strings.SplitN(raw, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
				return c.JSON(http.StatusUnauthorized, utilsErrors.Error{Message: "missing or malformed Authorization header"})
			}

			res, err := verifier.VerifyToken(c.Request().Context(), qryAuthorize.VerifyTokenQuery{Token: parts[1]})
			if err != nil || res == nil {
				return c.JSON(http.StatusUnauthorized, utilsErrors.Error{Message: "invalid token"})
			}

			ctx := context.WithValue(c.Request().Context(), "userID", res.UserID)
			ctx = context.WithValue(ctx, "userRole", res.Role)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if role, _ := c.Request().Context().Value("userRole").(string); role != "admin" {
			return c.JSON(http.StatusForbidden, utilsErrors.Error{Message: "admin only"})
		}
		return next(c)
	}
}
