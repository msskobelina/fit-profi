package authorize

import (
	"context"
	"net/http"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type RegisterHandler interface {
	Register(ctx context.Context, cmd cmdAuthorize.RegisterUserCommand) (*cmdAuthorize.RegisterUserResult, error)
}

func RegisterController(io controller.IO, h RegisterHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			FullName string `json:"fullName" validate:"required"`
			Email    string `json:"email"    validate:"required,email"`
			Password string `json:"password" validate:"required,min=6"`
		}
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
