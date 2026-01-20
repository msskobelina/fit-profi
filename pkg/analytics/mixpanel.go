package analytics

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client interface {
	Track(ctx context.Context, event string, distinctID string, props map[string]any) error
}

type Mixpanel struct {
	token   string
	apiHost string
	http    *http.Client
}

func NewMixpanel(token, apiHost string) *Mixpanel {
	if apiHost == "" {
		apiHost = os.Getenv("MIXPANEL_API_HOST")
	}
	return &Mixpanel{
		token:   token,
		apiHost: strings.TrimRight(apiHost, "/"),
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (m *Mixpanel) Track(ctx context.Context, event, distinctID string, props map[string]any) error {
	if m == nil || m.token == "" || event == "" || distinctID == "" {
		return nil
	}

	if props == nil {
		props = map[string]any{}
	}

	props["token"] = m.token
	props["distinct_id"] = distinctID

	payload := []map[string]any{
		{"event": event, "properties": props},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(b)

	form := url.Values{}
	form.Set("data", encoded)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		m.apiHost+"/track",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 || strings.TrimSpace(string(body)) != "1" {
		return fmt.Errorf("mixpanel track failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}
