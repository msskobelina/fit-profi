package authorize

import (
	"context"
	"net/http"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type LoginHandler interface {
	Login(ctx context.Context, cmd cmdAuthorize.LoginUserCommand) (*cmdAuthorize.LoginUserResult, error)
}

func LoginController(io controller.IO, h LoginHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"    validate:"required,email"`
			Password string `json:"password" validate:"required"`
		}
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
