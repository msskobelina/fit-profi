package authorize

import (
	"context"
	"net/http"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type SendEmailHandler interface {
	SendResetEmail(ctx context.Context, cmd cmdAuthorize.SendResetEmailCommand) error
}

func SendEmailController(io controller.IO, h SendEmailHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email string `json:"email" validate:"required,email"`
		}
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		if err := h.SendResetEmail(r.Context(), cmdAuthorize.SendResetEmailCommand{Email: req.Email}); err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(map[string]any{}, w)
	})
}
