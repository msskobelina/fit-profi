package nutrition

import (
	"context"
	"net/http"
	"strconv"

	cmdNutrition "github.com/msskobelina/fit-profi/internal/application/command/nutrition"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

// UpdateEntryRequest is the body for PUT /nutrition/entries/:id.
type UpdateEntryRequest struct {
	MealType string            `json:"mealType" validate:"required,oneof=breakfast lunch dinner snack" example:"dinner"`
	Items    []model.DiaryItem `json:"items"    validate:"required,min=1"`
}

type UpdateEntryHandler interface {
	UpdateEntry(context.Context, cmdNutrition.UpdateEntryCommand) (*model.DiaryEntry, error)
}

// UpdateEntryController godoc
//
//	@Summary		Update nutrition diary entry
//	@Description	Replaces the meal type and food items of an existing diary entry.
//	@Tags			Nutrition
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Entry ID"
//	@Param			body	body		UpdateEntryRequest	true	"Updated entry data"
//	@Success		200		{object}	model.DiaryEntry
//	@Failure		400		{object}	controller.ErrorResponse
//	@Failure		401		{object}	controller.ErrorResponse
//	@Router			/nutrition/entries/{id} [put]
func UpdateEntryController(io controller.IO, h UpdateEntryHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		idStr := controller.PathParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			io.Error(err, r, w)
			return
		}
		var req UpdateEntryRequest
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.UpdateEntry(r.Context(), cmdNutrition.UpdateEntryCommand{
			EntryID:  id,
			UserID:   userID,
			MealType: req.MealType,
			Items:    req.Items,
		})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
