package programs

import (
	"context"
	"net/http"
	"strconv"

	cmdPrograms "github.com/msskobelina/fit-profi/internal/application/command/programs"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type DeleteProgramHandler interface {
	DeleteProgram(context.Context, cmdPrograms.DeleteProgramCommand) error
}

// DeleteProgramController godoc
//
//	@Summary		Delete training program
//	@Description	Deletes a training program owned by the authenticated user.
//	@Tags			Programs
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Program ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	controller.ErrorResponse
//	@Failure		401	{object}	controller.ErrorResponse
//	@Router			/programs/{id} [delete]
func DeleteProgramController(io controller.IO, h DeleteProgramHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		idStr := controller.PathParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			io.Error(err, r, w)
			return
		}
		if err := h.DeleteProgram(r.Context(), cmdPrograms.DeleteProgramCommand{
			ProgramID: id,
			UserID:    userID,
		}); err != nil {
			io.Error(err, r, w)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
