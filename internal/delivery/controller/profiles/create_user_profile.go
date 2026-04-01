package profiles

import (
	"context"
	"net/http"

	"github.com/msskobelina/fit-profi/internal/application/command/profiles"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type (
	CreateUserProfileHandler interface {
		CreateUserProfile(context.Context, profiles.CreateUserProfileCommand) (*model.UserProfile, error)
	}

	// CreateUserProfileRequest is the body for POST /profiles/user.
	CreateUserProfileRequest struct {
		FullName    string  `json:"fullName"    validate:"required"                                                   example:"John Doe"`
		Age         int     `json:"age"         validate:"required,gt=0,lt=130"                                       example:"25"`
		WeightKg    float32 `json:"weightKg"    validate:"required,gt=0"                                              example:"75.5"`
		Goal        string  `json:"goal"        validate:"required,oneof=lose_weight gain_weight rehab keep_fit competition" example:"keep_fit"`
		Description string  `json:"description"                                                                        example:"I want to stay healthy"`
	}
)

// CreateUserProfileController godoc
//
//	@Summary		Create user profile
//	@Description	Creates a fitness profile for the authenticated user.
//	@Tags			Profiles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateUserProfileRequest	true	"User profile data"
//	@Success		200		{object}	model.UserProfile
//	@Failure		400		{object}	controller.ErrorResponse
//	@Failure		401		{object}	controller.ErrorResponse
//	@Router			/profiles/user [post]
func CreateUserProfileController(
	io controller.IO,
	handler CreateUserProfileHandler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)

		var req CreateUserProfileRequest

		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}

		cmd := profiles.CreateUserProfileCommand{
			UserID:      userID,
			FullName:    req.FullName,
			Age:         req.Age,
			WeightKg:    req.WeightKg,
			Goal:        req.Goal,
			Description: req.Description,
		}

		resp, err := handler.CreateUserProfile(r.Context(), cmd)
		if err != nil {
			io.Error(err, r, w)
			return
		}

		io.Result(resp, w)
	})
}
