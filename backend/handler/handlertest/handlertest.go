// Package handlertest builds an echo app with sessions, validation and an
// in-memory database, for testing handlers end to end without a browser.
package handlertest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/repository/repotest"
)

func New(t *testing.T) *echo.Echo {
	t.Helper()
	repotest.Open(t)
	e := echo.New()
	e.Validator = &handler.CustomValidator{Validator: validator.New()}
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(strings.Repeat("k", 32)))))
	return e
}

// JSON sends a JSON request and returns the response.
func JSON(e *echo.Echo, method, path, body string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}
