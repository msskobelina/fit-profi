package authorize

import (
	"context"
	"net/http"
	"strings"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type LogoutHandler interface {
	Logout(ctx context.Context, cmd cmdAuthorize.LogoutUserCommand) error
}

func LogoutController(io controller.IO, h LogoutHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("Authorization")
		token := ""
		if len(raw) > 7 && strings.HasPrefix(strings.ToLower(raw), "bearer ") {
			token = raw[7:]
		}
		_ = h.Logout(r.Context(), cmdAuthorize.LogoutUserCommand{Token: token})
		w.WriteHeader(http.StatusNoContent)
	})
}
