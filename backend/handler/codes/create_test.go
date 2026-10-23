package codes_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler/codes"
	"github.com/randspace0/ketikode/backend/handler/handlertest"
	repocodes "github.com/randspace0/ketikode/backend/repository/codes"
)

// Posts the same form the home page submits, against a stub executor.
func TestCreateFromHomePageForm(t *testing.T) {
	var got url.Values
	stub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		got = r.PostForm
		w.Write([]byte(`{"stdout":"45\n","stderr":"","success":true}`))
	}))
	defer stub.Close()
	t.Setenv("EXECUTOR_URL", stub.URL)

	e := handlertest.New(t)
	e.POST("/code/create", codes.Create)

	form := url.Values{"language": {"python3"}, "code": {"print(sum(range(10)))"}, "stdin": {"in"}}
	req := httptest.NewRequest(http.MethodPost, "/code/create", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/1" {
		t.Fatalf("got %d %q: %s", rec.Code, rec.Header().Get("Location"), rec.Body)
	}
	if got.Get("language") != "python3" || got.Get("code") != "print(sum(range(10)))" || got.Get("stdin") != "in" {
		t.Fatalf("executor received %v", got)
	}
	saved, _ := (&repocodes.Codes{ID: 1}).FindByID()
	if saved == nil || saved.Stdout != "45\n" || saved.Owner != 0 {
		t.Fatalf("saved submission: %+v", saved)
	}
}

func TestCreateRequiresLanguage(t *testing.T) {
	e := handlertest.New(t)
	e.POST("/code/create", codes.Create)

	req := httptest.NewRequest(http.MethodPost, "/code/create", strings.NewReader("code=x"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusFound {
		t.Fatal("a submission without a language must not be saved")
	}
}
