package users

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/repository/users"
)

func SettingsHandler(c echo.Context) error {
	data := map[string]interface{}{
		"title": "Settings",
	}
	return c.Render(http.StatusOK, "user/settings.html", handler.AppendSessionData(c, data))
}

type settingsForm struct {
	PasswordOld string `json:"passwordOld" validate:"required"`
	Password    string `json:"password" validate:"required,min=8,max=72,nefield=PasswordOld"`
}

func SettingsProcessHandler(c echo.Context) error {
	sf := new(settingsForm)
	if err := c.Bind(sf); err != nil {
		return handler.JSONError(c, http.StatusBadRequest, "Invalid request", err)
	}

	if err := c.Validate(sf); err != nil {
		return handler.JSONError(c, http.StatusBadRequest, "The new password must be 8 to 72 characters and differ from the old one", nil)
	}

	// RequiresSignin has already checked the session user exists.
	id := handler.GetSessionData(c)["id"].(uint64)
	user, err := (&users.Users{ID: id}).FindByID()
	if err != nil || user == nil {
		return handler.JSONError(c, http.StatusInternalServerError, "Could not change the password", err)
	}

	if !user.CheckPassword(sf.PasswordOld) {
		return handler.JSONError(c, http.StatusUnauthorized, "Wrong password", nil)
	}

	hash, err := users.HashPassword(sf.Password)
	if err != nil {
		return handler.JSONError(c, http.StatusInternalServerError, "Could not change the password", err)
	}
	if err := users.UpdatePassword(id, hash); err != nil {
		return handler.JSONError(c, http.StatusInternalServerError, "Could not change the password", err)
	}

	return handler.JSONOK(c)
}
