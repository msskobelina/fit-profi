package nutrition

import (
	"context"
	"net/http"
	"strconv"

	cmdNutrition "github.com/msskobelina/fit-profi/internal/application/command/nutrition"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type UpdateEntryHandler interface {
	UpdateEntry(context.Context, cmdNutrition.UpdateEntryCommand) (*model.DiaryEntry, error)
}

func UpdateEntryController(io controller.IO, h UpdateEntryHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		idStr := controller.PathParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			io.Error(err, r, w)
			return
		}
		var req struct {
			MealType string            `json:"mealType"`
			Items    []model.DiaryItem `json:"items"`
		}
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
