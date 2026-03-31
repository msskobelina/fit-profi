package authorize

import (
	"context"
	"net/http"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

// SendEmailRequest is the body for POST /users/send-email.
type SendEmailRequest struct {
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

type SendEmailHandler interface {
	SendResetEmail(ctx context.Context, cmd cmdAuthorize.SendResetEmailCommand) error
}

// SendEmailController godoc
//
//	@Summary		Send password-reset email
//	@Description	Sends a one-time password-reset link to the given email address.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SendEmailRequest		true	"Email address"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	controller.ErrorResponse
//	@Router			/users/send-email [post]
func SendEmailController(io controller.IO, h SendEmailHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SendEmailRequest
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
