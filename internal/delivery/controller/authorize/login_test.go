package authorize_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	"github.com/msskobelina/fit-profi/internal/delivery/boundary"
	"github.com/msskobelina/fit-profi/internal/delivery/controller/authorize"
)

type mockLoginHandler struct {
	result *cmdAuthorize.LoginUserResult
	err    error
}

func (m *mockLoginHandler) Login(_ context.Context, _ cmdAuthorize.LoginUserCommand) (*cmdAuthorize.LoginUserResult, error) {
	return m.result, m.err
}

func TestLoginController(t *testing.T) {
	successResult := &cmdAuthorize.LoginUserResult{
		Token:    "jwt-token",
		UserID:   1,
		FullName: "John Doe",
		Email:    "john@example.com",
	}

	tests := []struct {
		name       string
		body       string
		handler    *mockLoginHandler
		wantStatus int
		wantErrKey string
	}{
		{
			name:       "valid login",
			body:       `{"email":"john@example.com","password":"secret123"}`,
			handler:    &mockLoginHandler{result: successResult},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing email",
			body:       `{"password":"secret123"}`,
			handler:    &mockLoginHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "invalid email format",
			body:       `{"email":"not-an-email","password":"secret123"}`,
			handler:    &mockLoginHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "missing password",
			body:       `{"email":"john@example.com"}`,
			handler:    &mockLoginHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "invalid JSON",
			body:       `{not-json`,
			handler:    &mockLoginHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "wrong password from handler",
			body:       `{"email":"john@example.com","password":"wrongpass"}`,
			handler:    &mockLoginHandler{err: &testError{"Wrong password"}},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "user not found from handler",
			body:       `{"email":"nobody@example.com","password":"somepass"}`,
			handler:    &mockLoginHandler{err: &testError{"user not found"}},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := boundary.New()
			h := authorize.LoginController(io, tt.handler)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", w.Code, tt.wantStatus, w.Body.String())
			}

			if tt.wantErrKey != "" {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if _, ok := resp[tt.wantErrKey]; !ok {
					t.Errorf("expected key %q in response, got: %v", tt.wantErrKey, resp)
				}
			}
		})
	}
}

func TestLoginController_ResponseBody(t *testing.T) {
	expected := &cmdAuthorize.LoginUserResult{
		Token:    "my-token",
		UserID:   42,
		FullName: "Jane Doe",
		Email:    "jane@example.com",
	}

	io := boundary.New()
	h := authorize.LoginController(io, &mockLoginHandler{result: expected})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login",
		strings.NewReader(`{"email":"jane@example.com","password":"password"}`))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var got cmdAuthorize.LoginUserResult
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got.Token != expected.Token {
		t.Errorf("Token = %q, want %q", got.Token, expected.Token)
	}
	if got.UserID != expected.UserID {
		t.Errorf("UserID = %d, want %d", got.UserID, expected.UserID)
	}
	if got.Email != expected.Email {
		t.Errorf("Email = %q, want %q", got.Email, expected.Email)
	}
}
