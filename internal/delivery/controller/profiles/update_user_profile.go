package profiles

import (
	"context"
	"net/http"

	cmdProfiles "github.com/msskobelina/fit-profi/internal/application/command/profiles"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

// UpdateUserProfileRequest is the body for PUT /profiles/user.
type UpdateUserProfileRequest struct {
	FullName    string  `json:"fullName"    validate:"required"                                                   example:"John Doe"`
	Age         int     `json:"age"         validate:"required,gt=0,lt=130"                                       example:"26"`
	WeightKg    float32 `json:"weightKg"    validate:"required,gt=0"                                              example:"74.0"`
	Goal        string  `json:"goal"        validate:"required,oneof=lose_weight gain_weight rehab keep_fit competition" example:"lose_weight"`
	Description string  `json:"description"                                                                        example:"Updated description"`
}

type UpdateUserProfileHandler interface {
	UpdateUserProfile(context.Context, cmdProfiles.UpdateUserProfileCommand) (*model.UserProfile, error)
}

// UpdateUserProfileController godoc
//
//	@Summary		Update user profile
//	@Description	Replaces the fitness profile for the authenticated user.
//	@Tags			Profiles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateUserProfileRequest	true	"Updated profile data"
//	@Success		200		{object}	model.UserProfile
//	@Failure		400		{object}	controller.ErrorResponse
//	@Failure		401		{object}	controller.ErrorResponse
//	@Router			/profiles/user [put]
func UpdateUserProfileController(io controller.IO, h UpdateUserProfileHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req UpdateUserProfileRequest
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
