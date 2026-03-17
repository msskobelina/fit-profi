package authorize

import (
	"context"
	"net/http"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type ResetPasswordHandler interface {
	ResetPassword(ctx context.Context, cmd cmdAuthorize.ResetPasswordCommand) error
}

func ResetPasswordController(io controller.IO, h ResetPasswordHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Token    string `json:"token"    validate:"required"`
			Password string `json:"password" validate:"required,min=6"`
		}
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		if err := h.ResetPassword(r.Context(), cmdAuthorize.ResetPasswordCommand{
			Token:    req.Token,
			Password: req.Password,
		}); err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(map[string]any{}, w)
	})
}
