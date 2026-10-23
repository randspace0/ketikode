package users

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/repository/users"
)

func SignupHandler(c echo.Context) error {
	data := map[string]interface{}{
		"title": "Signup",
	}
	return c.Render(http.StatusOK, "user/signup.html", handler.AppendSessionData(c, data))
}

type signupForm struct {
	Username             string `json:"username" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,max=72"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
}

func SignupProcessHandler(c echo.Context) error {
	usf := new(signupForm)
	if err := c.Bind(usf); err != nil {
		return handler.JSONError(c, http.StatusBadRequest, "Invalid request", err)
	}

	if err := c.Validate(usf); err != nil {
		return handler.JSONError(c, http.StatusBadRequest, "Invalid input: username, a valid email, and a password of 8 to 72 characters (matching the confirmation) are required", nil)
	}

	hash, err := users.HashPassword(usf.Password)
	if err != nil {
		return handler.JSONError(c, http.StatusInternalServerError, "Could not create the account", err)
	}

	u := users.Users{
		Username: usf.Username,
		Email:    usf.Email,
		Password: hash,
	}

	if err := u.Create(); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return handler.JSONError(c, http.StatusConflict, "Username or email is already taken", nil)
		}
		return handler.JSONError(c, http.StatusInternalServerError, "Could not create the account", err)
	}

	return handler.JSONOK(c)
}
