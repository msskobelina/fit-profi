package authorize

import (
	"net/http"

	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

// CheckResponse is the response body for GET /users/check.
type CheckResponse struct {
	UserID int    `json:"userId" example:"42"`
	Role   string `json:"role"   example:"user"`
}

// CheckController godoc
//
//	@Summary		Check auth token
//	@Description	Returns the user ID and role extracted from the JWT token.
//	@Tags			Auth
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	CheckResponse
//	@Failure		401	{object}	controller.ErrorResponse
//	@Router			/users/check [get]
func CheckController(io controller.IO) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Result(map[string]any{
			"userId": r.Context().Value("userID"),
			"role":   r.Context().Value("userRole"),
		}, w)
	})
}
