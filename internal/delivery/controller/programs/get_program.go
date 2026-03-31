package programs

import (
	"context"
	"net/http"
	"strconv"

	qryPrograms "github.com/msskobelina/fit-profi/internal/application/query/programs"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type GetProgramHandler interface {
	GetProgram(context.Context, qryPrograms.GetProgramQuery) (*model.TrainingProgram, error)
}

// GetProgramController godoc
//
//	@Summary		Get training program
//	@Description	Returns a training program by ID.
//	@Tags			Programs
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Program ID"
//	@Success		200	{object}	model.TrainingProgram
//	@Failure		400	{object}	controller.ErrorResponse
//	@Failure		401	{object}	controller.ErrorResponse
//	@Router			/programs/{id} [get]
func GetProgramController(io controller.IO, h GetProgramHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := controller.PathParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.GetProgram(r.Context(), qryPrograms.GetProgramQuery{ProgramID: id})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
