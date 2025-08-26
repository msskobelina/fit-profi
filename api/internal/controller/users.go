package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msskobelina/fit-profi/internal/service"
	"github.com/msskobelina/fit-profi/pkg/logger"
)

type userRoutes struct {
	service service.Services
	logger  logger.Interface
}

func newUserRoutes(g *echo.Group, s service.Services, l logger.Interface) {
	r := &userRoutes{s, l}
	ug := g.Group("/users")

	// public
	ug.POST("/register", r.registerUser)
	ug.POST("/login", r.loginUser)
	ug.POST("/sendemail", r.sendEmail)
	ug.PATCH("/resetpassword", r.resetPassword)

	// private
	auth := ug.Group("", authMiddleware(s))
	auth.POST("/logout", r.logoutUser)
	auth.GET("/check", r.check)
}

type registerUserRequestBody struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type registerUserResponse struct {
	Token    string         `json:"token"`
	UserID   int            `json:"userId"`
	FullName string         `json:"fullName"`
	Email    string         `json:"email"`
	Error    *service.Error `json:"error,omitempty"`
}

func (r *userRoutes) registerUser(c echo.Context) error {
	var body registerUserRequestBody
	if err := c.Bind(&body); err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := r.service.Users.RegisterUser(c.Request().Context(), &service.RegisterUserInput{
		FullName: body.FullName, Email: body.Email, Password: body.Password,
	})
	if err != nil {
		if se, ok := err.(*service.Error); ok {
			return c.JSON(http.StatusBadRequest, registerUserResponse{Error: se})
		}
		return errorResponse(c, http.StatusInternalServerError, "failed to register user")
	}
	return c.JSON(http.StatusOK, registerUserResponse{
		Token:    out.Token,
		UserID:   out.UserID,
		FullName: out.FullName,
		Email:    out.Email,
	})
}

type loginUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginUserResponse struct {
	Token    string         `json:"token"`
	UserID   int            `json:"userId"`
	FullName string         `json:"fullName"`
	Email    string         `json:"email"`
	Error    *service.Error `json:"error,omitempty"`
}

func (r *userRoutes) loginUser(c echo.Context) error {
	var body loginUserRequestBody
	if err := c.Bind(&body); err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := r.service.Users.LoginUser(c.Request().Context(), &service.LoginUserInput{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		if se, ok := err.(*service.Error); ok {
			return c.JSON(http.StatusBadRequest, loginUserResponse{Error: se})
		}
		return errorResponse(c, http.StatusInternalServerError, "failed to login user")
	}
	return c.JSON(http.StatusOK, loginUserResponse{
		Token:    out.Token,
		UserID:   out.UserID,
		FullName: out.FullName,
		Email:    out.Email,
	})
}

func (r *userRoutes) logoutUser(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		_ = r.service.Users.Logout(c.Request().Context(), auth[7:])
	}

	return c.NoContent(http.StatusNoContent)
}

type sendEmailRequestBody struct {
	Email string `json:"email"`
}
type sendEmailResponse struct {
	Error *service.Error `json:"error,omitempty"`
}

func (r *userRoutes) sendEmail(c echo.Context) error {
	var body sendEmailRequestBody
	if err := c.Bind(&body); err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	if err := r.service.Users.SendEmail(c.Request().Context(), &service.SendUserEmailInput{Email: body.Email}); err != nil {
		if se, ok := err.(*service.Error); ok {
			return c.JSON(http.StatusBadRequest, sendEmailResponse{Error: se})
		}
		return errorResponse(c, http.StatusInternalServerError, "failed to send email")
	}
	return c.JSON(http.StatusOK, sendEmailResponse{})
}

type resetPasswordRequestBody struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
type resetPasswordResponse struct {
	Error *service.Error `json:"error,omitempty"`
}

func (r *userRoutes) resetPassword(c echo.Context) error {
	var body resetPasswordRequestBody
	if err := c.Bind(&body); err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	if err := r.service.Users.ResetPassword(c.Request().Context(), &service.ResetPasswordInput{
		Token:    body.Token,
		Password: body.Password,
	}); err != nil {
		if se, ok := err.(*service.Error); ok {
			return c.JSON(http.StatusBadRequest, resetPasswordResponse{Error: se})
		}
		return errorResponse(c, http.StatusInternalServerError, "failed to reset password")
	}
	return c.JSON(http.StatusOK, resetPasswordResponse{})
}

func (r *userRoutes) check(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"userId": c.Get("userID"),
		"role":   c.Get("userRole"),
	})
}
