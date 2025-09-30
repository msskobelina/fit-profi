package integrations

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Service interface {
	ConnectURL(ctx context.Context, userID int) (string, error)
	ExchangeCallback(ctx context.Context, state, code string) error
	GetGoogleClient(ctx context.Context, userID int) (*oauth2.Config, *oauth2.Token, *UserIntegration, error)
}

type service struct {
	repo  Repository
	oauth *oauth2.Config
	nowFn func() time.Time
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
		oauth: &oauth2.Config{
			ClientID:     must("GOOGLE_CLIENT_ID"),
			ClientSecret: must("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  must("GOOGLE_REDIRECT_URL"),
			Scopes:       []string{calendar.CalendarScope},
			Endpoint:     google.Endpoint,
		},
		nowFn: time.Now,
	}
}

func (s *service) ConnectURL(ctx context.Context, userID int) (string, error) {
	state := s.signState(userID)
	return s.oauth.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (s *service) ExchangeCallback(ctx context.Context, state, code string) error {
	uid, ok := s.verifyState(state, 10*time.Minute)
	if !ok || code == "" {
		return fmt.Errorf("invalid state/code")
	}
	tok, err := s.oauth.Exchange(ctx, code)
	if err != nil {
		return err
	}
	_, err = s.repo.Upsert(ctx, &UserIntegration{
		UserID:       uid,
		Provider:     ProviderGoogle,
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiryUnix:   tok.Expiry.Unix(),
		Scope:        calendar.CalendarScope,
		CalendarID:   "primary",
		Timezone:     "Europe/Kyiv", // TODO: взяти з профілю в майбутньому
	})
	return err
}

func (s *service) GetGoogleClient(ctx context.Context, userID int) (*oauth2.Config, *oauth2.Token, *UserIntegration, error) {
	integ, err := s.repo.GetByUser(ctx, userID, ProviderGoogle)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("google not connected")
	}
	tok := &oauth2.Token{
		AccessToken:  integ.AccessToken,
		RefreshToken: integ.RefreshToken,
		Expiry:       time.Unix(integ.ExpiryUnix, 0),
	}
	return s.oauth, tok, integ, nil
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing env " + k)
	}
	return v
}

func (s *service) signState(uid int) string {
	payload := fmt.Sprintf("%d:%d", uid, s.nowFn().Unix())
	mac := hmac.New(sha256.New, []byte(os.Getenv("HMAC_SECRET")))
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig
}

func (s *service) verifyState(state string, maxAge time.Duration) (int, bool) {
	parts := strings.Split(state, ".")
	if len(parts) != 2 {
		return 0, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, false
	}
	mac := hmac.New(sha256.New, []byte(os.Getenv("HMAC_SECRET")))
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
	if s.nowFn().Sub(time.Unix(ts, 0)) > maxAge {
		return 0, false
	}
	return uid, true
}
