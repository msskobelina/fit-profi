package nutrition

import (
	"context"
	"net/http"
	"time"

	qryNutrition "github.com/msskobelina/fit-profi/internal/application/query/nutrition"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type ListEntriesHandler interface {
	ListEntries(context.Context, qryNutrition.ListEntriesQuery) ([]model.DiaryEntry, error)
}

func ListEntriesController(io controller.IO, h ListEntriesHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		dateStr := r.URL.Query().Get("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			date = time.Now()
		}
		res, err := h.ListEntries(r.Context(), qryNutrition.ListEntriesQuery{UserID: userID, Date: date})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
