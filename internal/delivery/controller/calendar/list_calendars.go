package calendar

import (
	"context"
	"net/http"

	qryCalendar "github.com/msskobelina/fit-profi/internal/application/query/calendar"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type ListCalendarsHandler interface {
	ListCalendars(ctx context.Context, q qryCalendar.ListCalendarsQuery) ([]qryCalendar.CalendarInfo, error)
}

func ListCalendarsController(io controller.IO, h ListCalendarsHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		res, err := h.ListCalendars(r.Context(), qryCalendar.ListCalendarsQuery{UserID: userID})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
