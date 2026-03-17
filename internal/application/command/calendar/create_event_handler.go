package calendar

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	googlecalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type CreateEventHandler interface {
	CreateEvent(ctx context.Context, cmd CreateEventCommand) (*CreateEventResult, error)
}

type createEventService struct {
	integRepo repository.IntegrationsRepository
	oauthCfg  *oauth2.Config
}

func NewCreateEventService(integRepo repository.IntegrationsRepository, oauthCfg *oauth2.Config) CreateEventHandler {
	return &createEventService{integRepo: integRepo, oauthCfg: oauthCfg}
}

func (s *createEventService) CreateEvent(ctx context.Context, cmd CreateEventCommand) (*CreateEventResult, error) {
	integ, err := s.integRepo.GetByUserAndProvider(ctx, cmd.UserID, model.ProviderGoogle)
	if err != nil {
		return nil, err
	}

	tok := &oauth2.Token{
		AccessToken:  integ.AccessToken,
		RefreshToken: integ.RefreshToken,
		Expiry:       time.Unix(integ.ExpiryUnix, 0),
	}
	svc, err := googlecalendar.NewService(ctx, option.WithTokenSource(s.oauthCfg.TokenSource(ctx, tok)))
	if err != nil {
		return nil, err
	}

	calID := cmd.CalendarID
	if calID == "" {
		calID = "primary"
	}

	var reminders []*googlecalendar.EventReminder
	for _, m := range cmd.Reminders {
		reminders = append(reminders, &googlecalendar.EventReminder{Method: "popup", Minutes: m})
	}

	var attendees []*googlecalendar.EventAttendee
	for _, email := range cmd.Attendees {
		attendees = append(attendees, &googlecalendar.EventAttendee{Email: email})
	}

	rem := &googlecalendar.EventReminders{Overrides: reminders}
	if len(reminders) == 0 {
		rem.UseDefault = true
		rem.ForceSendFields = []string{"UseDefault"}
	} else {
		rem.UseDefault = false
		rem.ForceSendFields = []string{"UseDefault"}
	}

	ev := &googlecalendar.Event{
		Summary:     cmd.Title,
		Description: cmd.Notes,
		Start:       &googlecalendar.EventDateTime{DateTime: cmd.Start.Format(time.RFC3339), TimeZone: integ.Timezone},
		End:         &googlecalendar.EventDateTime{DateTime: cmd.End.Format(time.RFC3339), TimeZone: integ.Timezone},
		Reminders:   rem,
		Attendees:   attendees,
	}

	out, err := svc.Events.Insert(calID, ev).Do()
	if err != nil {
		return nil, err
	}

	return &CreateEventResult{EventID: out.Id}, nil
}
