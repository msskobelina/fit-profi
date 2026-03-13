package programs

import (
	"context"
	"net/http"

	cmdPrograms "github.com/msskobelina/fit-profi/internal/application/command/programs"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type TrackProgressHandler interface {
	TrackProgress(context.Context, cmdPrograms.TrackProgressCommand) (*model.ExerciseProgress, error)
}

func TrackProgressController(io controller.IO, h TrackProgressHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req struct {
			ExerciseID int    `json:"exerciseId"`
			Sets       int    `json:"sets"`
			Reps       int    `json:"reps"`
			WeightKg   int    `json:"weightKg"`
			Notes      string `json:"notes"`
		}
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.TrackProgress(r.Context(), cmdPrograms.TrackProgressCommand{
			UserID:     userID,
			ExerciseID: req.ExerciseID,
			Sets:       req.Sets,
			Reps:       req.Reps,
			WeightKg:   req.WeightKg,
			Notes:      req.Notes,
		})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
