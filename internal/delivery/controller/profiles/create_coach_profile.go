package profiles

import (
	"context"
	"net/http"

	cmdProfiles "github.com/msskobelina/fit-profi/internal/application/command/profiles"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type CreateCoachProfileHandler interface {
	CreateCoachProfile(context.Context, cmdProfiles.CreateCoachProfileCommand) (*model.CoachProfile, error)
}

func CreateCoachProfileController(io controller.IO, h CreateCoachProfileHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req struct {
			FullName     string                   `json:"fullName"`
			Category     string                   `json:"category"`
			Info         string                   `json:"info"`
			Achievements []model.CoachAchievement `json:"achievements"`
			Education    []model.CoachEducation   `json:"education"`
		}
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.CreateCoachProfile(r.Context(), cmdProfiles.CreateCoachProfileCommand{
			UserID:       userID,
			FullName:     req.FullName,
			Category:     req.Category,
			Info:         req.Info,
			Achievements: req.Achievements,
			Education:    req.Education,
		})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
