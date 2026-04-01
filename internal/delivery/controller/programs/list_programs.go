package programs

import (
	"context"
	"net/http"

	qryPrograms "github.com/msskobelina/fit-profi/internal/application/query/programs"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type ListProgramsHandler interface {
	ListPrograms(context.Context, qryPrograms.ListProgramsQuery) ([]model.TrainingProgram, error)
}

// ListProgramsController godoc
//
//	@Summary		List training programs
//	@Description	Returns all training programs belonging to the authenticated user.
//	@Tags			Programs
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		model.TrainingProgram
//	@Failure		400	{object}	controller.ErrorResponse
//	@Failure		401	{object}	controller.ErrorResponse
//	@Router			/programs [get]
func ListProgramsController(io controller.IO, h ListProgramsHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		res, err := h.ListPrograms(r.Context(), qryPrograms.ListProgramsQuery{UserID: userID})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
