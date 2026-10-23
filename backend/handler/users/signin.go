package users

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/repository/users"
)

func SigninHandler(c echo.Context) error {
	data := map[string]interface{}{
		"title": "Signin",
	}
	return c.Render(http.StatusOK, "user/signin.html", handler.AppendSessionData(c, data))
}

type signinForm struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func SigninProcessHandler(c echo.Context) error {
	sf := new(signinForm)
	if err := c.Bind(sf); err != nil {
		return handler.JSONError(c, http.StatusBadRequest, "Invalid request", err)
	}

	if err := c.Validate(sf); err != nil {
		return handler.JSONError(c, http.StatusBadRequest, "Username and password are required", nil)
	}

	user, err := users.FindByUsername(sf.Username)
	if err != nil {
		return handler.JSONError(c, http.StatusInternalServerError, "Could not sign in", err)
	}

	// Same answer for an unknown user and a wrong password.
	if user == nil || !user.CheckPassword(sf.Password) {
		return handler.JSONError(c, http.StatusUnauthorized, "Wrong username or password", nil)
	}

	handler.Signin(c, user.ID)

	return handler.JSONOK(c)
}
