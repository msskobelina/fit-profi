package integrations

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type ConnectGoogleHandler interface {
	ConnectGoogle(ctx context.Context, cmd ConnectGoogleCommand) (*ConnectGoogleResult, error)
}

type connectGoogleService struct {
	oauth      *oauth2.Config
	hmacSecret string
}

func NewConnectGoogleService(oauthCfg *oauth2.Config, hmacSecret string) ConnectGoogleHandler {
	return &connectGoogleService{oauth: oauthCfg, hmacSecret: hmacSecret}
}

func (s *connectGoogleService) ConnectGoogle(_ context.Context, cmd ConnectGoogleCommand) (*ConnectGoogleResult, error) {
	state := s.signState(cmd.UserID)
	url := s.oauth.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return &ConnectGoogleResult{RedirectURL: url}, nil
}

func (s *connectGoogleService) signState(uid int) string {
	payload := fmt.Sprintf("%d:%d", uid, time.Now().Unix())
	mac := hmac.New(sha256.New, []byte(s.hmacSecret))
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig
}
