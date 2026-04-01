package nutrition

import (
	"context"
	"net/http"
	"strconv"

	cmdNutrition "github.com/msskobelina/fit-profi/internal/application/command/nutrition"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type DeleteEntryHandler interface {
	DeleteEntry(context.Context, cmdNutrition.DeleteEntryCommand) error
}

// DeleteEntryController godoc
//
//	@Summary		Delete nutrition diary entry
//	@Description	Deletes a diary entry belonging to the authenticated user.
//	@Tags			Nutrition
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Entry ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	controller.ErrorResponse
//	@Failure		401	{object}	controller.ErrorResponse
//	@Router			/nutrition/entries/{id} [delete]
func DeleteEntryController(io controller.IO, h DeleteEntryHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		idStr := controller.PathParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			io.Error(err, r, w)
			return
		}
		if err := h.DeleteEntry(r.Context(), cmdNutrition.DeleteEntryCommand{
			EntryID: id,
			UserID:  userID,
		}); err != nil {
			io.Error(err, r, w)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
