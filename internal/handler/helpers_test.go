package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rearurides/eagle-bank/internal/handler/middleware"
)

func assertResponse[T any](
	t *testing.T,
	w *httptest.ResponseRecorder,
	wantStatus int,
	wantBody *T, wantErr *errorResponse,
) {
	t.Helper()

	if w.Code != wantStatus {
		t.Errorf("expected status: %d, got: %d", wantStatus, w.Code)
	}

	if wantBody != nil {
		var got T
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}

		if diff := cmp.Diff(*wantBody, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	} else if wantErr != nil {
		var got errorResponse
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to unmarshal error body: %v", err)
		}
		if diff := cmp.Diff(*wantErr, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func newAuthRequest(t *testing.T, method, path, body, tokenUserID string) *http.Request {
	t.Helper()
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, tokenUserID)
	return req.WithContext(ctx)
}
