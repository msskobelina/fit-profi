package profiles

import (
	"context"
	"net/http"

	qryProfiles "github.com/msskobelina/fit-profi/internal/application/query/profiles"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type GetCoachProfileHandler interface {
	GetCoachProfile(context.Context, qryProfiles.GetCoachProfileQuery) (*model.CoachProfile, error)
}

func GetCoachProfileController(io controller.IO, h GetCoachProfileHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		res, err := h.GetCoachProfile(r.Context(), qryProfiles.GetCoachProfileQuery{UserID: userID})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
