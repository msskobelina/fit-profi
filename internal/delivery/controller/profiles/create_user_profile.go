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

	createUserProfileRequest struct {
		UserID      int     `json:"userId"`
		FullName    string  `json:"fullName"`
		Age         int     `json:"age"`
		WeightKg    float32 `json:"weightKg"`
		Goal        string  `json:"goal"`
		Description string  `json:"description"`
	}
)

func CreateUserProfileController(
	io controller.IO,
	handler CreateUserProfileHandler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var req createUserProfileRequest

		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}

		cmd := profiles.CreateUserProfileCommand{
			UserID:      req.UserID,
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
