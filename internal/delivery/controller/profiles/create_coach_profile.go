package profiles

import (
	"context"
	"net/http"

	cmdProfiles "github.com/msskobelina/fit-profi/internal/application/command/profiles"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

// CreateCoachProfileRequest is the body for POST /profiles/coach.
type CreateCoachProfileRequest struct {
	FullName     string                   `json:"fullName"     validate:"required"                          example:"Jane Smith"`
	Category     string                   `json:"category"     validate:"required,oneof=standard master professional" example:"master"`
	Info         string                   `json:"info"         validate:"required"                          example:"10 years of coaching experience"`
	Achievements []model.CoachAchievement `json:"achievements"`
	Education    []model.CoachEducation   `json:"education"`
}

type CreateCoachProfileHandler interface {
	CreateCoachProfile(context.Context, cmdProfiles.CreateCoachProfileCommand) (*model.CoachProfile, error)
}

// CreateCoachProfileController godoc
//
//	@Summary		Create coach profile
//	@Description	Creates a coach profile for the authenticated user.
//	@Tags			Profiles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateCoachProfileRequest	true	"Coach profile data"
//	@Success		200		{object}	model.CoachProfile
//	@Failure		400		{object}	controller.ErrorResponse
//	@Failure		401		{object}	controller.ErrorResponse
//	@Router			/profiles/coach [post]
func CreateCoachProfileController(io controller.IO, h CreateCoachProfileHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req CreateCoachProfileRequest
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
