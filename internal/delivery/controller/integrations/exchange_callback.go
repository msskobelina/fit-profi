package integrations

import (
	"context"
	"net/http"

	cmdIntegrations "github.com/msskobelina/fit-profi/internal/application/command/integrations"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type ExchangeCallbackHandler interface {
	ExchangeCallback(ctx context.Context, cmd cmdIntegrations.ExchangeCallbackCommand) error
}

// ExchangeCallbackController godoc
//
//	@Summary		Google OAuth callback
//	@Description	Receives the OAuth authorization code from Google and stores the access token.
//	@Tags			Integrations
//	@Produce		plain
//	@Param			state	query		string	true	"OAuth state parameter"
//	@Param			code	query		string	true	"OAuth authorization code"
//	@Success		200		{string}	string	"Google connected ✓"
//	@Failure		400		{object}	controller.ErrorResponse
//	@Router			/integrations/google/callback [get]
func ExchangeCallbackController(io controller.IO, h ExchangeCallbackHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")
		if err := h.ExchangeCallback(r.Context(), cmdIntegrations.ExchangeCallbackCommand{
			State: state,
			Code:  code,
		}); err != nil {
			io.Error(err, r, w)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Google connected ✓"))
	})
}
