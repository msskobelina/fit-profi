package profiles

import (
	"context"
	"net/http"

	cmdProfiles "github.com/msskobelina/fit-profi/internal/application/command/profiles"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

// UpdateCoachProfileRequest is the body for PUT /profiles/coach.
type UpdateCoachProfileRequest struct {
	FullName     string                   `json:"fullName"     validate:"required"                          example:"Jane Smith"`
	Category     string                   `json:"category"     validate:"required,oneof=standard master professional" example:"professional"`
	Info         string                   `json:"info"         validate:"required"                          example:"Updated bio"`
	Achievements []model.CoachAchievement `json:"achievements"`
	Education    []model.CoachEducation   `json:"education"`
}

type UpdateCoachProfileHandler interface {
	UpdateCoachProfile(context.Context, cmdProfiles.UpdateCoachProfileCommand) (*model.CoachProfile, error)
}

// UpdateCoachProfileController godoc
//
//	@Summary		Update coach profile
//	@Description	Replaces the coach profile for the authenticated user.
//	@Tags			Profiles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateCoachProfileRequest	true	"Updated coach profile data"
//	@Success		200		{object}	model.CoachProfile
//	@Failure		400		{object}	controller.ErrorResponse
//	@Failure		401		{object}	controller.ErrorResponse
//	@Router			/profiles/coach [put]
func UpdateCoachProfileController(io controller.IO, h UpdateCoachProfileHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req UpdateCoachProfileRequest
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.UpdateCoachProfile(r.Context(), cmdProfiles.UpdateCoachProfileCommand{
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
