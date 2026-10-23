package adminutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAuth(t *testing.T) {
	e := echo.New()
	e.GET("/admin", func(c echo.Context) error { return c.String(http.StatusOK, "ok") }, Auth("root", "s3cret"))

	for _, tc := range []struct {
		name, user, pass string
		auth             bool
		want             int
	}{
		{"no credentials", "", "", false, http.StatusUnauthorized},
		{"wrong password", "root", "nope", true, http.StatusUnauthorized},
		{"wrong user", "other", "s3cret", true, http.StatusUnauthorized},
		{"correct", "root", "s3cret", true, http.StatusOK},
	} {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		if tc.auth {
			req.SetBasicAuth(tc.user, tc.pass)
		}
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, rec.Code, tc.want)
		}
	}
}

func TestCredentialsRequireBoth(t *testing.T) {
	t.Setenv("ADMIN_USER", "root")
	t.Setenv("ADMIN_PASSWORD", "")
	if _, _, ok := Credentials(); ok {
		t.Fatal("password missing: backoffice must stay disabled")
	}
	t.Setenv("ADMIN_PASSWORD", "x")
	if _, _, ok := Credentials(); !ok {
		t.Fatal("both set: backoffice must be enabled")
	}
}

func TestNewPager(t *testing.T) {
	for _, tc := range []struct {
		page       int
		total      int64
		wantPage   int
		prev, next int
		wantPages  int
	}{
		{1, 0, 1, 0, 0, 1},
		{1, 50, 1, 0, 0, 1},
		{1, 51, 1, 0, 2, 2},
		{2, 120, 2, 1, 3, 3},
		{99, 120, 3, 2, 0, 3}, // out of range clamps to last page
		{-4, 120, 1, 0, 2, 3},
	} {
		p := NewPager(tc.page, tc.total)
		if p.Page != tc.wantPage || p.Prev != tc.prev || p.Next != tc.next || p.Pages != tc.wantPages {
			t.Errorf("page=%d total=%d: got %+v", tc.page, tc.total, p)
		}
	}
}
