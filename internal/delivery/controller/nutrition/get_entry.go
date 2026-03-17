package nutrition

import (
	"context"
	"net/http"
	"strconv"

	qryNutrition "github.com/msskobelina/fit-profi/internal/application/query/nutrition"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type GetEntryHandler interface {
	GetEntry(context.Context, qryNutrition.GetEntryQuery) (*model.DiaryEntry, error)
}

func GetEntryController(io controller.IO, h GetEntryHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		idStr := controller.PathParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.GetEntry(r.Context(), qryNutrition.GetEntryQuery{EntryID: id, UserID: userID})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
