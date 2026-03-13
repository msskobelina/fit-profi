package profiles

import (
	"context"
	"net/http"

	qryProfiles "github.com/msskobelina/fit-profi/internal/application/query/profiles"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type GetUserProfileHandler interface {
	GetUserProfile(context.Context, qryProfiles.GetUserProfileQuery) (*model.UserProfile, error)
}

func GetUserProfileController(io controller.IO, h GetUserProfileHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		res, err := h.GetUserProfile(r.Context(), qryProfiles.GetUserProfileQuery{UserID: userID})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
