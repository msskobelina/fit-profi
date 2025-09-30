package calendar

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	integ "github.com/msskobelina/fit-profi/domains/integrations"
)

type Service interface {
	CreateSelfEvent(ctx context.Context, uid int, in CreateSelfEventInput) (string, error)
	ListUserCalendars(ctx context.Context, uid int) ([]CalendarInfo, error)
}

type CreateSelfEventInput struct {
	CalendarID string
	Title      string
	Notes      string
	Start      time.Time
	End        time.Time
	Reminders  []int64
	Attendees  []string
}

type TimeSlot struct{ Start, End time.Time }
type CalendarInfo struct{ ID, Summary string }

type service struct {
	repo  Repository
	integ integ.Service
}

func NewService(repo Repository, integSvc integ.Service) Service {
	return &service{repo: repo, integ: integSvc}
}

func (s *service) calSvc(ctx context.Context, userID int) (*calendar.Service, *integ.UserIntegration, *oauth2.Config, error) {
	cfg, tok, i, err := s.integ.GetGoogleClient(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	svc, err := calendar.NewService(ctx, option.WithTokenSource(cfg.TokenSource(ctx, tok)))
	if err != nil {
		return nil, nil, nil, err
	}
	return svc, i, cfg, nil
}

func (s *service) ListUserCalendars(ctx context.Context, uid int) ([]CalendarInfo, error) {
	svc, _, _, err := s.calSvc(ctx, uid)
	if err != nil {
		return nil, err
	}
	resp, err := svc.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	out := make([]CalendarInfo, 0, len(resp.Items))
	for _, it := range resp.Items {
		out = append(out, CalendarInfo{ID: it.Id, Summary: it.Summary})
	}
	return out, nil
}

func (s *service) CreateSelfEvent(ctx context.Context, uid int, in CreateSelfEventInput) (string, error) {
	svc, integ, _, err := s.calSvc(ctx, uid)
	if err != nil {
		return "", err
	}
	calID := in.CalendarID
	if calID == "" {
		calID = "primary"
	}

	var ro []*calendar.EventReminder
	for _, m := range in.Reminders {
		ro = append(ro, &calendar.EventReminder{Method: "popup", Minutes: m})
	}

	var atts []*calendar.EventAttendee
	for _, email := range in.Attendees {
		atts = append(atts, &calendar.EventAttendee{Email: email})
	}

	rem := &calendar.EventReminders{Overrides: ro}
	if len(ro) == 0 {
		rem.UseDefault = true
		rem.ForceSendFields = []string{"UseDefault"}
	} else {
		rem.UseDefault = false
		rem.ForceSendFields = []string{"UseDefault"}
	}

	ev := &calendar.Event{
		Summary:     in.Title,
		Description: in.Notes,
		Start:       &calendar.EventDateTime{DateTime: in.Start.Format(time.RFC3339), TimeZone: integ.Timezone},
		End:         &calendar.EventDateTime{DateTime: in.End.Format(time.RFC3339), TimeZone: integ.Timezone},
		Reminders:   rem,
		Attendees:   atts,
	}

	out, err := svc.Events.Insert(calID, ev).Do()
	if err != nil {
		return "", err
	}
	return out.Id, nil
}
