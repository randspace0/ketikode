package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler"
)

func TestRateLimitIsPerIP(t *testing.T) {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPDirect()
	e.POST("/run", func(c echo.Context) error { return c.NoContent(http.StatusOK) }, handler.RateLimit(0.001, 2))

	call := func(remote, forwarded string) int {
		req := httptest.NewRequest(http.MethodPost, "/run", nil)
		req.RemoteAddr = remote
		if forwarded != "" {
			req.Header.Set("X-Forwarded-For", forwarded)
		}
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		return rec.Code
	}

	for i := 0; i < 2; i++ {
		if code := call("10.0.0.1:1000", ""); code != http.StatusOK {
			t.Fatalf("request %d within the burst got %d", i+1, code)
		}
	}
	if code := call("10.0.0.1:1000", ""); code != http.StatusTooManyRequests {
		t.Fatalf("third request got %d, want 429", code)
	}
	if code := call("10.0.0.1:1000", "203.0.113.9"); code != http.StatusTooManyRequests {
		t.Fatal("a spoofed X-Forwarded-For must not reset the limit")
	}
	if code := call("10.0.0.2:1000", ""); code != http.StatusOK {
		t.Fatalf("another client got %d, want 200", code)
	}
}
