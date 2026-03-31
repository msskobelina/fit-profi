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

type mockRegisterHandler struct {
	result *cmdAuthorize.RegisterUserResult
	err    error
}

func (m *mockRegisterHandler) Register(_ context.Context, _ cmdAuthorize.RegisterUserCommand) (*cmdAuthorize.RegisterUserResult, error) {
	return m.result, m.err
}

func TestRegisterController(t *testing.T) {
	successResult := &cmdAuthorize.RegisterUserResult{
		Token:    "jwt-token",
		UserID:   1,
		FullName: "John Doe",
		Email:    "john@example.com",
	}

	tests := []struct {
		name       string
		body       string
		handler    *mockRegisterHandler
		wantStatus int
		wantErrKey string
	}{
		{
			name:       "valid registration",
			body:       `{"fullName":"John Doe","email":"john@example.com","password":"secret123"}`,
			handler:    &mockRegisterHandler{result: successResult},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing fullName",
			body:       `{"email":"john@example.com","password":"secret123"}`,
			handler:    &mockRegisterHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "invalid email",
			body:       `{"fullName":"John","email":"not-an-email","password":"secret123"}`,
			handler:    &mockRegisterHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "password too short",
			body:       `{"fullName":"John","email":"john@example.com","password":"abc"}`,
			handler:    &mockRegisterHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "missing password",
			body:       `{"fullName":"John","email":"john@example.com"}`,
			handler:    &mockRegisterHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "invalid JSON",
			body:       `{invalid}`,
			handler:    &mockRegisterHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "empty body",
			body:       ``,
			handler:    &mockRegisterHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "handler returns error",
			body:       `{"fullName":"John Doe","email":"john@example.com","password":"secret123"}`,
			handler:    &mockRegisterHandler{err: &testError{"email already exists"}},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := boundary.New()
			h := authorize.RegisterController(io, tt.handler)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/register", strings.NewReader(tt.body))
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

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
