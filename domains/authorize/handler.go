package authorize

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Register(g *echo.Group, authMW echo.MiddlewareFunc) {
	ug := g.Group("/users")
	ug.POST("/register", h.register)
	ug.POST("/login", h.login)
	ug.POST("/send-email", h.sendEmail)
	ug.PATCH("/reset-password", h.reset)

	priv := ug.Group("", authMW)
	priv.POST("/logout", h.logout)
	priv.GET("/check", h.check)
}

// swagger:route POST /users/register Users usersRegister
// Register new user
//
// responses:
//
//	200: body:authResponse
//	400: body:errorResponse
func (h *Handler) register(c echo.Context) error {
	req := new(RegisterUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	out, err := h.service.Register(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, outOrErr(out, err))
	}
	return c.JSON(http.StatusOK, out)
}

// swagger:route POST /users/login Users usersLogin
// Login user
//
// responses:
//
//	200: body:authResponse
//	400: body:errorResponse
func (h *Handler) login(c echo.Context) error {
	req := new(LoginUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	out, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, outOrErr(out, err))
	}
	return c.JSON(http.StatusOK, out)
}

// swagger:route POST /users/logout Users usersLogout
// Logout (revoke token)
//
// security:
//   - Bearer: []
//
// responses:
//
//	204: description: no content
//	401: body:errorResponse
func (h *Handler) logout(c echo.Context) error {
	raw := c.Request().Header.Get("Authorization")
	if len(raw) > 7 && strings.HasPrefix(strings.ToLower(raw), "bearer ") {
		_ = h.service.Logout(c.Request().Context(), raw[7:])
	}
	return c.NoContent(http.StatusNoContent)
}

// swagger:route POST /users/send-email Users usersSendEmail
// Send password reset email
//
// responses:
//
//	200: description: ok
//	400: body:errorResponse
func (h *Handler) sendEmail(c echo.Context) error {
	req := new(SendEmailRequest)
	if err := c.Bind(req); err != nil || req.Email == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	if err := h.service.SendEmail(c.Request().Context(), req.Email); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{})
}

// swagger:route PATCH /users/reset-password Users usersResetPassword
// Reset password by token
//
// responses:
//
//	200: description: ok
//	400: body:errorResponse
func (h *Handler) reset(c echo.Context) error {
	req := new(ResetPasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	if err := h.service.ResetPassword(c.Request().Context(), req.Token, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{})
}

// swagger:route GET /users/check Users usersCheck
// Get current user context (from token)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: description: ok
//	401: body:errorResponse
func (h *Handler) check(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"userId":   c.Get("userID"),
		"role":     c.Get("userRole"),
		"fullName": c.Get("userFullName"),
	})
}

func outOrErr(out *AuthResponse, err error) any {
	if out != nil {
		return out
	}
	return ErrorResponse{Error: err.Error()}
}
