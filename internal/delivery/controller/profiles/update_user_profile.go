package profiles

import (
	"context"
	"net/http"

	cmdProfiles "github.com/msskobelina/fit-profi/internal/application/command/profiles"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type UpdateUserProfileHandler interface {
	UpdateUserProfile(context.Context, cmdProfiles.UpdateUserProfileCommand) (*model.UserProfile, error)
}

func UpdateUserProfileController(io controller.IO, h UpdateUserProfileHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req struct {
			FullName    string  `json:"fullName"`
			Age         int     `json:"age"`
			WeightKg    float32 `json:"weightKg"`
			Goal        string  `json:"goal"`
			Description string  `json:"description"`
		}
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.UpdateUserProfile(r.Context(), cmdProfiles.UpdateUserProfileCommand{
			UserID:      userID,
			FullName:    req.FullName,
			Age:         req.Age,
			WeightKg:    req.WeightKg,
			Goal:        req.Goal,
			Description: req.Description,
		})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
