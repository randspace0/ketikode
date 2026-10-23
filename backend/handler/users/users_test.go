package users_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/handler/handlertest"
	"github.com/randspace0/ketikode/backend/handler/users"
	"github.com/randspace0/ketikode/backend/repository"
	repo "github.com/randspace0/ketikode/backend/repository/users"
)

func app(t *testing.T) *echo.Echo {
	e := handlertest.New(t)
	e.POST("/signup", users.SignupProcessHandler)
	e.POST("/signin", users.SigninProcessHandler)
	e.POST("/settings", handler.RequiresSignin(users.SettingsProcessHandler))
	return e
}

const signupBody = `{"username":"alice","email":"a@x.io","password":"correct-horse","passwordConfirmation":"correct-horse"}`

func signinCookie(t *testing.T, e *echo.Echo, password string) *http.Cookie {
	t.Helper()
	rec := handlertest.JSON(e, "POST", "/signin", `{"username":"alice","password":"`+password+`"}`)
	if rec.Code != http.StatusOK || len(rec.Result().Cookies()) == 0 {
		t.Fatalf("signin failed: %d %s", rec.Code, rec.Body)
	}
	return rec.Result().Cookies()[0]
}

func TestSignupStoresHashedPassword(t *testing.T) {
	e := app(t)
	if rec := handlertest.JSON(e, "POST", "/signup", signupBody); rec.Code != http.StatusOK {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
	u, _ := repo.FindByUsername("alice")
	if u == nil || u.Password == "correct-horse" || !u.CheckPassword("correct-horse") {
		t.Fatalf("password must be stored hashed: %+v", u)
	}
}

func TestSignupRejectsBadInput(t *testing.T) {
	e := app(t)
	for name, body := range map[string]string{
		"short password": `{"username":"a","email":"a@x.io","password":"short","passwordConfirmation":"short"}`,
		"too long":       `{"username":"a","email":"a@x.io","password":"` + strings.Repeat("a", 73) + `","passwordConfirmation":"` + strings.Repeat("a", 73) + `"}`,
		"mismatch":       `{"username":"a","email":"a@x.io","password":"correct-horse","passwordConfirmation":"different-one"}`,
		"invalid email":  `{"username":"a","email":"nope","password":"correct-horse","passwordConfirmation":"correct-horse"}`,
		"malformed json": `{`,
	} {
		if rec := handlertest.JSON(e, "POST", "/signup", body); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: got %d", name, rec.Code)
		}
	}
}

func TestSignupDuplicateIsConflictWithoutLeakingSQL(t *testing.T) {
	e := app(t)
	handlertest.JSON(e, "POST", "/signup", signupBody)
	rec := handlertest.JSON(e, "POST", "/signup", signupBody)
	if rec.Code != http.StatusConflict {
		t.Fatalf("got %d %s", rec.Code, rec.Body)
	}
	if body := strings.ToLower(rec.Body.String()); strings.Contains(body, "constraint") || strings.Contains(body, "sql") {
		t.Fatalf("database error leaked: %s", body)
	}
}

func TestSignin(t *testing.T) {
	e := app(t)
	handlertest.JSON(e, "POST", "/signup", signupBody)

	wrong := handlertest.JSON(e, "POST", "/signin", `{"username":"alice","password":"wrong"}`)
	unknown := handlertest.JSON(e, "POST", "/signin", `{"username":"nobody","password":"wrong"}`)
	if wrong.Code != http.StatusUnauthorized || unknown.Code != http.StatusUnauthorized || wrong.Body.String() != unknown.Body.String() {
		t.Fatalf("wrong and unknown user must look the same: %d %s / %d %s", wrong.Code, wrong.Body, unknown.Code, unknown.Body)
	}
	signinCookie(t, e, "correct-horse")
}

func TestChangePassword(t *testing.T) {
	e := app(t)
	handlertest.JSON(e, "POST", "/signup", signupBody)
	cookie := signinCookie(t, e, "correct-horse")

	if rec := handlertest.JSON(e, "POST", "/settings", `{"passwordOld":"correct-horse","password":"new-password"}`); rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("signed-out request must redirect, got %d", rec.Code)
	}
	if rec := handlertest.JSON(e, "POST", "/settings", `{"passwordOld":"wrong","password":"new-password"}`, cookie); rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong old password: %d", rec.Code)
	}
	if rec := handlertest.JSON(e, "POST", "/settings", `{"passwordOld":"correct-horse","password":"short"}`, cookie); rec.Code != http.StatusBadRequest {
		t.Fatalf("short new password: %d", rec.Code)
	}
	if rec := handlertest.JSON(e, "POST", "/settings", `{"passwordOld":"correct-horse","password":"new-password"}`, cookie); rec.Code != http.StatusOK {
		t.Fatalf("change failed: %d %s", rec.Code, rec.Body)
	}

	if rec := handlertest.JSON(e, "POST", "/signin", `{"username":"alice","password":"correct-horse"}`); rec.Code != http.StatusUnauthorized {
		t.Fatal("old password must stop working")
	}
	signinCookie(t, e, "new-password")
}

func TestCookieOfDeletedUserIsSignedOut(t *testing.T) {
	e := app(t)
	handlertest.JSON(e, "POST", "/signup", signupBody)
	cookie := signinCookie(t, e, "correct-horse")

	repository.DB.Where("username = ?", "alice").Delete(&repo.Users{})

	rec := handlertest.JSON(e, "POST", "/settings", `{"passwordOld":"x","password":"new-password"}`, cookie)
	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("want redirect to signin, got %d", rec.Code)
	}
}
