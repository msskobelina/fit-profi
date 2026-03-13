package programs

import (
	"context"
	"net/http"

	cmdPrograms "github.com/msskobelina/fit-profi/internal/application/command/programs"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type CreateProgramHandler interface {
	CreateProgram(context.Context, cmdPrograms.CreateProgramCommand) (*model.TrainingProgram, error)
}

func CreateProgramController(io controller.IO, h CreateProgramHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req struct {
			Title       string             `json:"title"`
			Description string             `json:"description"`
			Days        []model.ProgramDay `json:"days"`
		}
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
