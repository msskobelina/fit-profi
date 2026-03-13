package nutrition

import (
	"context"
	"net/http"
	"time"

	cmdNutrition "github.com/msskobelina/fit-profi/internal/application/command/nutrition"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type CreateEntryHandler interface {
	CreateEntry(context.Context, cmdNutrition.CreateEntryCommand) (*model.DiaryEntry, error)
}

func CreateEntryController(io controller.IO, h CreateEntryHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req struct {
			Date     string            `json:"date"`
			MealType string            `json:"mealType"`
			Items    []model.DiaryItem `json:"items"`
		}
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.CreateEntry(r.Context(), cmdNutrition.CreateEntryCommand{
			UserID:   userID,
			Date:     date,
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
