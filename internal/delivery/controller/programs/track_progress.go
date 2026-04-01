package programs

import (
	"context"
	"net/http"

	cmdPrograms "github.com/msskobelina/fit-profi/internal/application/command/programs"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

// TrackProgressRequest is the body for POST /programs/progress.
type TrackProgressRequest struct {
	ExerciseID int    `json:"exerciseId" validate:"required,gt=0" example:"5"`
	Sets       int    `json:"sets"       validate:"required,gt=0" example:"3"`
	Reps       int    `json:"reps"       validate:"required,gt=0" example:"12"`
	WeightKg   int    `json:"weightKg"   validate:"gte=0"         example:"60"`
	Notes      string `json:"notes"                               example:"Felt strong today"`
}

type TrackProgressHandler interface {
	TrackProgress(context.Context, cmdPrograms.TrackProgressCommand) (*model.ExerciseProgress, error)
}

// TrackProgressController godoc
//
//	@Summary		Track exercise progress
//	@Description	Records a completed set of an exercise for the authenticated user.
//	@Tags			Programs
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		TrackProgressRequest	true	"Progress data"
//	@Success		200		{object}	model.ExerciseProgress
//	@Failure		400		{object}	controller.ErrorResponse
//	@Failure		401		{object}	controller.ErrorResponse
//	@Router			/programs/progress [post]
func TrackProgressController(io controller.IO, h TrackProgressHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req TrackProgressRequest
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
