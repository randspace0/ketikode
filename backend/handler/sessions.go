package handler

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/randspace0/ketikode/backend/repository/users"
)

// TODO refactor session part

// IsLoggedIn is true only when the session user still exists. A cookie for a
// removed user is cleared.
func IsLoggedIn(c echo.Context) bool {
	sess, err := session.Get("user", c)
	if err != nil {
		return false
	}
	id, ok := sess.Values["id"].(uint64)
	if !ok {
		return false
	}

	user, err := (&users.Users{ID: id}).FindByID()
	if err != nil || user == nil {
		if err == nil {
			sess.Values = map[interface{}]interface{}{}
			sess.Save(c.Request(), c.Response())
		}
		return false
	}
	sess.Values["username"] = user.Username
	sess.Save(c.Request(), c.Response())
	return true
}

func RequiresSignin(handler func(c echo.Context) error) func(echo.Context) error {
	return func(c echo.Context) error {
		if !IsLoggedIn(c) {
			return c.Redirect(http.StatusTemporaryRedirect, "/signin")
		}

		return handler(c)
	}
}

func RequiresSignout(handler func(c echo.Context) error) func(echo.Context) error {
	return func(c echo.Context) error {
		if IsLoggedIn(c) {
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		return handler(c)
	}
}

func Signin(c echo.Context, id uint64) bool {
	sess, err := session.Get("user", c)
	if err != nil || sess.Values["id"] != nil {
		return false
	}
	sess.Values["id"] = id
	sess.Save(c.Request(), c.Response())
	return true
}

func Signout(c echo.Context) bool {
	sess, err := session.Get("user", c)
	if err != nil || sess.Values["id"] == nil {
		return false
	}
	sess.Values = map[interface{}]interface{}{}
	sess.Save(c.Request(), c.Response())
	return true
}

func AppendSessionData(c echo.Context, data map[string]interface{}) map[string]interface{} {
	if IsLoggedIn(c) {
		sess, _ := session.Get("user", c)
		data["session"] = sess.Values
	}
	return data
}

func GetSessionData(c echo.Context) map[interface{}]interface{} {
	sess, err := session.Get("user", c)
	if err != nil {
		return map[interface{}]interface{}{}
	}
	return sess.Values
}
