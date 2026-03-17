package integrations

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type ExchangeCallbackHandler interface {
	ExchangeCallback(ctx context.Context, cmd ExchangeCallbackCommand) error
}

type exchangeCallbackService struct {
	repo       repository.IntegrationsRepository
	oauth      *oauth2.Config
	hmacSecret string
}

func NewExchangeCallbackService(
	repo repository.IntegrationsRepository,
	oauthCfg *oauth2.Config,
	hmacSecret string,
) ExchangeCallbackHandler {
	return &exchangeCallbackService{repo: repo, oauth: oauthCfg, hmacSecret: hmacSecret}
}

func (s *exchangeCallbackService) ExchangeCallback(ctx context.Context, cmd ExchangeCallbackCommand) error {
	uid, ok := s.verifyState(cmd.State, 10*time.Minute)
	if !ok || cmd.Code == "" {
		return fmt.Errorf("invalid state/code")
	}

	tok, err := s.oauth.Exchange(ctx, cmd.Code)
	if err != nil {
		return err
	}

	_, err = s.repo.Upsert(ctx, model.UserIntegration{
		UserID:       uid,
		Provider:     model.ProviderGoogle,
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiryUnix:   tok.Expiry.Unix(),
		Scope:        calendar.CalendarScope,
		CalendarID:   "primary",
		Timezone:     "Europe/Kyiv",
	})

	return err
}

func (s *exchangeCallbackService) verifyState(state string, maxAge time.Duration) (int, bool) {
	parts := strings.Split(state, ".")
	if len(parts) != 2 {
		return 0, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, false
	}
	mac := hmac.New(sha256.New, []byte(s.hmacSecret))
	mac.Write(raw)
	if base64.RawURLEncoding.EncodeToString(mac.Sum(nil)) != parts[1] {
		return 0, false
	}
	sp := strings.Split(string(raw), ":")
	if len(sp) != 2 {
		return 0, false
	}
	uid, err := strconv.Atoi(sp[0])
	if err != nil {
		return 0, false
	}
	ts, err := strconv.ParseInt(sp[1], 10, 64)
	if err != nil {
		return 0, false
	}
	if time.Now().Sub(time.Unix(ts, 0)) > maxAge {
		return 0, false
	}

	return uid, true
}
