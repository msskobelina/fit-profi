package authorize

import (
	"context"
	"net/http"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

// RegisterRequest is the body for POST /users/register.
type RegisterRequest struct {
	FullName string `json:"fullName" validate:"required"                 example:"John Doe"`
	Email    string `json:"email"    validate:"required,email"            example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6"            example:"secret123"`
}

type RegisterHandler interface {
	Register(ctx context.Context, cmd cmdAuthorize.RegisterUserCommand) (*cmdAuthorize.RegisterUserResult, error)
}

// RegisterController godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new account and returns a JWT token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RegisterRequest					true	"Registration payload"
//	@Success		200		{object}	cmdAuthorize.RegisterUserResult
//	@Failure		400		{object}	controller.ErrorResponse
//	@Router			/users/register [post]
func RegisterController(io controller.IO, h RegisterHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.Register(r.Context(), cmdAuthorize.RegisterUserCommand{
			FullName: req.FullName,
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
