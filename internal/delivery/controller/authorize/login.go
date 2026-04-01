package authorize

import (
	"context"
	"net/http"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

// LoginRequest is the body for POST /users/login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required"       example:"secret123"`
}

type LoginHandler interface {
	Login(ctx context.Context, cmd cmdAuthorize.LoginUserCommand) (*cmdAuthorize.LoginUserResult, error)
}

// LoginController godoc
//
//	@Summary		Login
//	@Description	Authenticates the user and returns a JWT token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest					true	"Login credentials"
//	@Success		200		{object}	cmdAuthorize.LoginUserResult
//	@Failure		400		{object}	controller.ErrorResponse
//	@Router			/users/login [post]
func LoginController(io controller.IO, h LoginHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.Login(r.Context(), cmdAuthorize.LoginUserCommand{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
