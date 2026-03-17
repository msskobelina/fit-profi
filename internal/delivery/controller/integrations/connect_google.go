package integrations

import (
	"context"
	"net/http"

	cmdIntegrations "github.com/msskobelina/fit-profi/internal/application/command/integrations"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type ConnectGoogleHandler interface {
	ConnectGoogle(ctx context.Context, cmd cmdIntegrations.ConnectGoogleCommand) (*cmdIntegrations.ConnectGoogleResult, error)
}

func ConnectGoogleController(io controller.IO, h ConnectGoogleHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		res, err := h.ConnectGoogle(r.Context(), cmdIntegrations.ConnectGoogleCommand{UserID: userID})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		http.Redirect(w, r, res.RedirectURL, http.StatusFound)
	})
}
