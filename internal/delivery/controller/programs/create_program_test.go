package programs_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cmdPrograms "github.com/msskobelina/fit-profi/internal/application/command/programs"
	"github.com/msskobelina/fit-profi/internal/delivery/boundary"
	"github.com/msskobelina/fit-profi/internal/delivery/controller/programs"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type mockCreateProgramHandler struct {
	result *model.TrainingProgram
	err    error
	gotCmd cmdPrograms.CreateProgramCommand
}

func (m *mockCreateProgramHandler) CreateProgram(_ context.Context, cmd cmdPrograms.CreateProgramCommand) (*model.TrainingProgram, error) {
	m.gotCmd = cmd
	return m.result, m.err
}

func requestWithUserID(method, path, body string, userID int) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req.WithContext(context.WithValue(req.Context(), "userID", userID))
}

func TestCreateProgramController(t *testing.T) {
	successResult := &model.TrainingProgram{
		ID:    1,
		Title: "My Program",
	}

	tests := []struct {
		name       string
		body       string
		userID     int
		handler    *mockCreateProgramHandler
		wantStatus int
		wantErrKey string
	}{
		{
			name:       "valid request with title only",
			body:       `{"title":"My Program"}`,
			userID:     5,
			handler:    &mockCreateProgramHandler{result: successResult},
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid request with days",
			body:       `{"title":"Full Body","description":"3 days/week","days":[{"dayNumber":1,"title":"Day 1"}]}`,
			userID:     5,
			handler:    &mockCreateProgramHandler{result: successResult},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing title",
			body:       `{"description":"No title here"}`,
			userID:     5,
			handler:    &mockCreateProgramHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "empty title",
			body:       `{"title":""}`,
			userID:     5,
			handler:    &mockCreateProgramHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "invalid JSON",
			body:       `{bad json`,
			userID:     5,
			handler:    &mockCreateProgramHandler{},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name:       "handler error",
			body:       `{"title":"My Program"}`,
			userID:     5,
			handler:    &mockCreateProgramHandler{err: &testError{"db error"}},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := boundary.New()
			h := programs.CreateProgramController(io, tt.handler)

			req := requestWithUserID(http.MethodPost, "/api/v1/programs", tt.body, tt.userID)
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

func TestCreateProgramController_UserIDFromContext(t *testing.T) {
	handler := &mockCreateProgramHandler{result: &model.TrainingProgram{ID: 1, Title: "Test"}}
	io := boundary.New()
	h := programs.CreateProgramController(io, handler)

	req := requestWithUserID(http.MethodPost, "/api/v1/programs", `{"title":"Test"}`, 42)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d; body: %s", w.Code, w.Body.String())
	}
	if handler.gotCmd.UserID != 42 {
		t.Errorf("UserID = %d, want 42", handler.gotCmd.UserID)
	}
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
