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
