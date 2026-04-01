package nutrition_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cmdNutrition "github.com/msskobelina/fit-profi/internal/application/command/nutrition"
	"github.com/msskobelina/fit-profi/internal/delivery/boundary"
	"github.com/msskobelina/fit-profi/internal/delivery/controller/nutrition"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type mockCreateEntryHandler struct {
	result *model.DiaryEntry
	err    error
	gotCmd cmdNutrition.CreateEntryCommand
}

func (m *mockCreateEntryHandler) CreateEntry(_ context.Context, cmd cmdNutrition.CreateEntryCommand) (*model.DiaryEntry, error) {
	m.gotCmd = cmd
	return m.result, m.err
}

func requestWithUserID(method, path, body string, userID int) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req.WithContext(context.WithValue(req.Context(), "userID", userID))
}

func TestCreateEntryController(t *testing.T) {
	successResult := &model.DiaryEntry{ID: 1, UserID: 5}

	validBody := `{"date":"2024-01-15","mealType":"breakfast","items":[{"name":"Oats","grams":100,"calories":350}]}`

	tests := []struct {
		name       string
		body       string
		userID     int
		handler    *mockCreateEntryHandler
		wantStatus int
		wantErrKey string
	}{
		{
			name:       "valid request",
			body:       validBody,
			userID:     5,
			handler:    &mockCreateEntryHandler{result: successResult},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing date",
			body:       `{"mealType":"breakfast","items":[{"name":"Oats","grams":100,"calories":350}]}`,
			userID:     5,
			handler:    &mockCreateEntryHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "invalid mealType",
			body:       `{"date":"2024-01-15","mealType":"brunch","items":[{"name":"Oats","grams":100}]}`,
			userID:     5,
			handler:    &mockCreateEntryHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "missing mealType",
			body:       `{"date":"2024-01-15","items":[{"name":"Oats","grams":100}]}`,
			userID:     5,
			handler:    &mockCreateEntryHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "empty items list",
			body:       `{"date":"2024-01-15","mealType":"lunch","items":[]}`,
			userID:     5,
			handler:    &mockCreateEntryHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "missing items",
			body:       `{"date":"2024-01-15","mealType":"dinner"}`,
			userID:     5,
			handler:    &mockCreateEntryHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "invalid date format",
			body:       `{"date":"15-01-2024","mealType":"breakfast","items":[{"name":"Oats"}]}`,
			userID:     5,
			handler:    &mockCreateEntryHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "all valid mealTypes",
			body:       `{"date":"2024-01-15","mealType":"snack","items":[{"name":"Apple","grams":150}]}`,
			userID:     5,
			handler:    &mockCreateEntryHandler{result: successResult},
			wantStatus: http.StatusOK,
		},
		{
			name:       "handler error",
			body:       validBody,
			userID:     5,
			handler:    &mockCreateEntryHandler{err: &testError{"db error"}},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := boundary.New()
			h := nutrition.CreateEntryController(io, tt.handler)

			req := requestWithUserID(http.MethodPost, "/api/v1/nutrition/entries", tt.body, tt.userID)
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

func TestCreateEntryController_PassesUserIDFromContext(t *testing.T) {
	handler := &mockCreateEntryHandler{result: &model.DiaryEntry{ID: 1, UserID: 15}}
	io := boundary.New()
	h := nutrition.CreateEntryController(io, handler)

	req := requestWithUserID(http.MethodPost, "/api/v1/nutrition/entries",
		`{"date":"2024-01-15","mealType":"lunch","items":[{"name":"Rice","grams":200}]}`, 15)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d; body: %s", w.Code, w.Body.String())
	}
	if handler.gotCmd.UserID != 15 {
		t.Errorf("UserID passed to handler = %d, want 15", handler.gotCmd.UserID)
	}
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
