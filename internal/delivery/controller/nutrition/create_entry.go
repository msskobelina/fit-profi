package nutrition

import (
	"context"
	"net/http"
	"time"

	cmdNutrition "github.com/msskobelina/fit-profi/internal/application/command/nutrition"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

// CreateEntryRequest is the body for POST /nutrition/entries.
type CreateEntryRequest struct {
	Date     string            `json:"date"     validate:"required"                               example:"2024-03-15"`
	MealType string            `json:"mealType" validate:"required,oneof=breakfast lunch dinner snack" example:"lunch"`
	Items    []model.DiaryItem `json:"items"    validate:"required,min=1"`
}

type CreateEntryHandler interface {
	CreateEntry(context.Context, cmdNutrition.CreateEntryCommand) (*model.DiaryEntry, error)
}

// CreateEntryController godoc
//
//	@Summary		Create nutrition diary entry
//	@Description	Records a meal with food items for the authenticated user.
//	@Tags			Nutrition
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateEntryRequest	true	"Diary entry"
//	@Success		200		{object}	model.DiaryEntry
//	@Failure		400		{object}	controller.ErrorResponse
//	@Failure		401		{object}	controller.ErrorResponse
//	@Router			/nutrition/entries [post]
func CreateEntryController(io controller.IO, h CreateEntryHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req CreateEntryRequest
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
