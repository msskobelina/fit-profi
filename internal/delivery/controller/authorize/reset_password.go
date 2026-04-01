package authorize

import (
	"context"
	"net/http"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

// ResetPasswordRequest is the body for PATCH /users/reset-password.
type ResetPasswordRequest struct {
	Token    string `json:"token"    validate:"required"       example:"eyJhbGciOiJIUzI1NiJ9..."`
	Password string `json:"password" validate:"required,min=6" example:"newSecret123"`
}

type ResetPasswordHandler interface {
	ResetPassword(ctx context.Context, cmd cmdAuthorize.ResetPasswordCommand) error
}

// ResetPasswordController godoc
//
//	@Summary		Reset password
//	@Description	Resets the user password using a one-time reset token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ResetPasswordRequest	true	"Reset token and new password"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	controller.ErrorResponse
//	@Router			/users/reset-password [patch]
func ResetPasswordController(io controller.IO, h ResetPasswordHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ResetPasswordRequest
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
