package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain_Ordering(t *testing.T) {
	var order []string
	makeMw := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}

	h := Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		makeMw("first"),
		makeMw("second"),
		makeMw("third"),
	)

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	expected := []string{"first", "second", "third"}
	for i, name := range expected {
		if order[i] != name {
			t.Errorf("expected middleware %s at position %d, got %s", name, i, order[i])
		}
	}
}

func TestRecoverPanic_500(t *testing.T) {
	panic := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	h := RecoverPanic(panic)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
	if rec.Body.String() != "internal server error\n" {
		t.Errorf("expected body 'internal server error', got %q", rec.Body.String())
	}
}

func TestRecoveryPanic_NoPanic(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := RecoverPanic(ok)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestLogging(t *testing.T) {
	// This test just ensures that the Logging middleware doesn't interfere with normal operation.
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := Logging(ok)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestLogging_StatusCode(t *testing.T) {
	// This test ensures that the Logging middleware correctly captures the status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot) // 418
	})
	h := Logging(handler)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != http.StatusTeapot {
		t.Errorf("expected status 418, got %d", rec.Code)
	}
}
