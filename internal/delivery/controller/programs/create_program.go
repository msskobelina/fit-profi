package programs

import (
	"context"
	"net/http"

	cmdPrograms "github.com/msskobelina/fit-profi/internal/application/command/programs"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

// CreateProgramRequest is the body for POST /programs.
type CreateProgramRequest struct {
	Title       string             `json:"title"       validate:"required" example:"Beginner Full Body"`
	Description string             `json:"description"                     example:"3-day full-body routine"`
	Days        []model.ProgramDay `json:"days"`
}

type CreateProgramHandler interface {
	CreateProgram(context.Context, cmdPrograms.CreateProgramCommand) (*model.TrainingProgram, error)
}

// CreateProgramController godoc
//
//	@Summary		Create training program
//	@Description	Creates a new training program with optional days and exercises.
//	@Tags			Programs
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateProgramRequest	true	"Program data"
//	@Success		200		{object}	model.TrainingProgram
//	@Failure		400		{object}	controller.ErrorResponse
//	@Failure		401		{object}	controller.ErrorResponse
//	@Router			/programs [post]
func CreateProgramController(io controller.IO, h CreateProgramHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req CreateProgramRequest
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.CreateProgram(r.Context(), cmdPrograms.CreateProgramCommand{
			UserID:      userID,
			Title:       req.Title,
			Description: req.Description,
			Days:        req.Days,
		})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
