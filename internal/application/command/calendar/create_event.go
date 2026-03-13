package calendar

import "time"

type CreateEventCommand struct {
	UserID     int
	CalendarID string
	Title      string
	Notes      string
	Start      time.Time
	End        time.Time
	Reminders  []int64
	Attendees  []string
}

type CreateEventResult struct {
	EventID string
}
