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

type ListCalendarsHandler interface {
	ListCalendars(ctx context.Context, q ListCalendarsQuery) ([]CalendarInfo, error)
}

type listCalendarsService struct {
	integRepo repository.IntegrationsRepository
	oauthCfg  *oauth2.Config
}

func NewListCalendarsService(integRepo repository.IntegrationsRepository, oauthCfg *oauth2.Config) ListCalendarsHandler {
	return &listCalendarsService{integRepo: integRepo, oauthCfg: oauthCfg}
}

func (s *listCalendarsService) ListCalendars(ctx context.Context, q ListCalendarsQuery) ([]CalendarInfo, error) {
	integ, err := s.integRepo.GetByUserAndProvider(ctx, q.UserID, model.ProviderGoogle)
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
