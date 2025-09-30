package calendar

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) Register(g *echo.Group, authMW echo.MiddlewareFunc) {
	r := g.Group("/calendar", authMW)

	r.GET("/list", h.listCalendars)
	r.POST("/me/events", h.createSelfEvent)
}

func (h *Handler) listCalendars(c echo.Context) error {
	uid := c.Get("userID").(int)
	out, err := h.svc.ListUserCalendars(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, out)
}

type createSelfEventReq struct {
	CalendarID string   `json:"calendarId,omitempty"`
	Title      string   `json:"title"`
	Notes      string   `json:"notes"`
	Start      string   `json:"start"`
	End        string   `json:"end"`
	Reminders  []int64  `json:"reminders,omitempty"`
	Attendees  []string `json:"attendees,omitempty"`
}

func (h *Handler) createSelfEvent(c echo.Context) error {
	uid := c.Get("userID").(int)
	in := new(createSelfEventReq)
	if err := c.Bind(in); err != nil {
		return c.JSON(400, map[string]string{"error": "bad body"})
	}
	st, err := time.Parse(time.RFC3339, in.Start)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "bad start"})
	}
	en, err := time.Parse(time.RFC3339, in.End)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "bad end"})
	}
	id, err := h.svc.CreateSelfEvent(c.Request().Context(), uid, CreateSelfEventInput{
		CalendarID: in.CalendarID,
		Title:      in.Title,
		Notes:      in.Notes,
		Start:      st,
		End:        en,
		Reminders:  in.Reminders,
		Attendees:  in.Attendees,
	})
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]any{"eventId": id})
}
