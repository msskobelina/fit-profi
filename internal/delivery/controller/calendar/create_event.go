package calendar

import (
	"context"
	"net/http"
	"time"

	cmdCalendar "github.com/msskobelina/fit-profi/internal/application/command/calendar"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

// CreateEventRequest is the body for POST /calendar/me/events.
type CreateEventRequest struct {
	CalendarID string   `json:"calendarId,omitempty"           example:"primary"`
	Title      string   `json:"title"       validate:"required" example:"Training session"`
	Notes      string   `json:"notes"                           example:"Leg day"`
	Start      string   `json:"start"       validate:"required" example:"2024-03-15T10:00:00Z"`
	End        string   `json:"end"         validate:"required" example:"2024-03-15T11:00:00Z"`
	Reminders  []int64  `json:"reminders,omitempty"             example:"10,30"`
	Attendees  []string `json:"attendees,omitempty"             example:"coach@example.com"`
}

// CreateEventResult documents the response for POST /calendar/me/events.
type CreateEventResult struct {
	EventID string `json:"eventId" example:"abc123xyz"`
}

type CreateEventHandler interface {
	CreateEvent(ctx context.Context, cmd cmdCalendar.CreateEventCommand) (*cmdCalendar.CreateEventResult, error)
}

// CreateEventController godoc
//
//	@Summary		Create calendar event
//	@Description	Creates a new event in the user's connected Google Calendar.
//	@Tags			Calendar
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateEventRequest	true	"Event data"
//	@Success		200		{object}	CreateEventResult
//	@Failure		400		{object}	controller.ErrorResponse
//	@Failure		401		{object}	controller.ErrorResponse
//	@Router			/calendar/me/events [post]
func CreateEventController(io controller.IO, h CreateEventHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userID").(int)
		var req CreateEventRequest
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
