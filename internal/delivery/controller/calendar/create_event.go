package calendar

import (
	"context"
	"net/http"
	"time"

	cmdCalendar "github.com/msskobelina/fit-profi/internal/application/command/calendar"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type CreateEventHandler interface {
	CreateEvent(ctx context.Context, cmd cmdCalendar.CreateEventCommand) (*cmdCalendar.CreateEventResult, error)
}

func CreateEventController(io controller.IO, h CreateEventHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req struct {
			CalendarID string   `json:"calendarId,omitempty"`
			Title      string   `json:"title"`
			Notes      string   `json:"notes"`
			Start      string   `json:"start"`
			End        string   `json:"end"`
			Reminders  []int64  `json:"reminders,omitempty"`
			Attendees  []string `json:"attendees,omitempty"`
		}
		if err := io.Read(&req, r.Body); err != nil {
			io.Error(err, r, w)
			return
		}
		start, err := time.Parse(time.RFC3339, req.Start)
		if err != nil {
			io.Error(err, r, w)
			return
		}
		end, err := time.Parse(time.RFC3339, req.End)
		if err != nil {
			io.Error(err, r, w)
			return
		}
		res, err := h.CreateEvent(r.Context(), cmdCalendar.CreateEventCommand{
			UserID:     userID,
			CalendarID: req.CalendarID,
			Title:      req.Title,
			Notes:      req.Notes,
			Start:      start,
			End:        end,
			Reminders:  req.Reminders,
			Attendees:  req.Attendees,
		})
		if err != nil {
			io.Error(err, r, w)
			return
		}
		io.Result(res, w)
	})
}
