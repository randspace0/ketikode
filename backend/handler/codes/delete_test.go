package codes_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/handler/codes"
	"github.com/randspace0/ketikode/backend/handler/handlertest"
	"github.com/randspace0/ketikode/backend/handler/users"
	repocodes "github.com/randspace0/ketikode/backend/repository/codes"
)

func app(t *testing.T) *echo.Echo {
	e := handlertest.New(t)
	e.POST("/signup", users.SignupProcessHandler)
	e.POST("/signin", users.SigninProcessHandler)
	e.DELETE("/code/delete", handler.RequiresSignin(codes.Delete))
	return e
}

func signup(t *testing.T, e *echo.Echo, name string) *http.Cookie {
	t.Helper()
	body := `{"username":"` + name + `","email":"` + name + `@x.io","password":"correct-horse","passwordConfirmation":"correct-horse"}`
	if rec := handlertest.JSON(e, "POST", "/signup", body); rec.Code != http.StatusOK {
		t.Fatalf("signup: %d %s", rec.Code, rec.Body)
	}
	rec := handlertest.JSON(e, "POST", "/signin", `{"username":"`+name+`","password":"correct-horse"}`)
	return rec.Result().Cookies()[0]
}

func TestDeleteOnlyOwnCode(t *testing.T) {
	e := app(t)
	alice := signup(t, e, "alice")
	bob := signup(t, e, "bob")

	c := &repocodes.Codes{Language: "cpp", Code: "x", Owner: 1} // alice
	if err := c.Create(); err != nil {
		t.Fatal(err)
	}
	body := `{"id":"1"}`

	if rec := handlertest.JSON(e, "DELETE", "/code/delete", body, bob); rec.Code != http.StatusNotFound {
		t.Fatalf("bob deleting alice's code: got %d %s", rec.Code, rec.Body)
	}
	if got, _ := (&repocodes.Codes{ID: c.ID}).FindByID(); got == nil {
		t.Fatal("code was deleted by a non-owner")
	}
	if rec := handlertest.JSON(e, "DELETE", "/code/delete", body); rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("signed out: got %d", rec.Code)
	}
	if rec := handlertest.JSON(e, "DELETE", "/code/delete", body, alice); rec.Code != http.StatusOK {
		t.Fatalf("owner delete: got %d %s", rec.Code, rec.Body)
	}
	if got, _ := (&repocodes.Codes{ID: c.ID}).FindByID(); got != nil {
		t.Fatal("owner delete did not remove the code")
	}
}
