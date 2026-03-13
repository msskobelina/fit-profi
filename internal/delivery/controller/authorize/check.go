package authorize

import (
	"net/http"

	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

func CheckController(io controller.IO) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Result(map[string]any{
			"userId": r.Context().Value("userID"),
			"role":   r.Context().Value("userRole"),
		}, w)
	})
}
