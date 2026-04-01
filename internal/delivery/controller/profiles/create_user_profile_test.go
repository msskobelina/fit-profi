package profiles_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cmdProfiles "github.com/msskobelina/fit-profi/internal/application/command/profiles"
	"github.com/msskobelina/fit-profi/internal/delivery/boundary"
	"github.com/msskobelina/fit-profi/internal/delivery/controller/profiles"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type mockCreateUserProfileHandler struct {
	result *model.UserProfile
	err    error
	gotCmd cmdProfiles.CreateUserProfileCommand
}

func (m *mockCreateUserProfileHandler) CreateUserProfile(_ context.Context, cmd cmdProfiles.CreateUserProfileCommand) (*model.UserProfile, error) {
	m.gotCmd = cmd
	return m.result, m.err
}

func requestWithUserID(method, path, body string, userID int) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req.WithContext(context.WithValue(req.Context(), "userID", userID))
}

func TestCreateUserProfileController(t *testing.T) {
	successResult := &model.UserProfile{
		ID:       1,
		UserID:   10,
		FullName: "John Doe",
		Age:      25,
		WeightKg: 75.0,
		Goal:     model.Goal("lose_weight"),
	}

	tests := []struct {
		name       string
		body       string
		userID     int
		handler    *mockCreateUserProfileHandler
		wantStatus int
		wantErrKey string
	}{
		{
			name:       "valid request",
			body:       `{"fullName":"John Doe","age":25,"weightKg":75.0,"goal":"lose_weight"}`,
			userID:     10,
			handler:    &mockCreateUserProfileHandler{result: successResult},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing fullName",
			body:       `{"age":25,"weightKg":75.0,"goal":"lose_weight"}`,
			userID:     10,
			handler:    &mockCreateUserProfileHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "invalid goal",
			body:       `{"fullName":"John","age":25,"weightKg":75.0,"goal":"invalid_goal"}`,
			userID:     10,
			handler:    &mockCreateUserProfileHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "negative age",
			body:       `{"fullName":"John","age":-1,"weightKg":75.0,"goal":"lose_weight"}`,
			userID:     10,
			handler:    &mockCreateUserProfileHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "zero weightKg",
			body:       `{"fullName":"John","age":25,"weightKg":0,"goal":"lose_weight"}`,
			userID:     10,
			handler:    &mockCreateUserProfileHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "userID from context not body",
			body:       `{"fullName":"John","age":25,"weightKg":75.0,"goal":"keep_fit"}`,
			userID:     99,
			handler:    &mockCreateUserProfileHandler{result: successResult},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := boundary.New()
			h := profiles.CreateUserProfileController(io, tt.handler)

			req := requestWithUserID(http.MethodPost, "/api/v1/profiles/user", tt.body, tt.userID)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", w.Code, tt.wantStatus, w.Body.String())
			}

			if tt.wantErrKey != "" {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("decode: %v", err)
				}
				if _, ok := resp[tt.wantErrKey]; !ok {
					t.Errorf("expected key %q in response, got: %v", tt.wantErrKey, resp)
				}
			}
		})
	}
}

func TestCreateUserProfileController_UserIDFromContext(t *testing.T) {
	handler := &mockCreateUserProfileHandler{result: &model.UserProfile{UserID: 7}}
	io := boundary.New()
	h := profiles.CreateUserProfileController(io, handler)

	req := requestWithUserID(http.MethodPost, "/api/v1/profiles/user",
		`{"fullName":"Alice","age":30,"weightKg":60.0,"goal":"keep_fit"}`, 7)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if handler.gotCmd.UserID != 7 {
		t.Errorf("UserID passed to handler = %d, want 7", handler.gotCmd.UserID)
	}
}
