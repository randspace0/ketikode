package users

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/cikadev/ketide/handler"
	"gitlab.com/cikadev/ketide/repository/users"
)

func SettingsHandler(c echo.Context) error {
	data := map[string]interface{}{
		"title": "Settings",
	}
	return c.Render(http.StatusOK, "user/settings.html", handler.AppendSessionData(c, data))
}

type settingsForm struct {
	PasswordOld string `json:"passwordOld" validate:"required"`
	Password    string `json:"password" validate:"required,nefield=PasswordOld"`
}

func SettingsProcessHandler(c echo.Context) error {
	data := handler.AppendSessionData(c, map[string]interface{}{})

	sf := new(settingsForm)

	if err := c.Bind(sf); err != nil {
		return c.JSON(http.StatusBadGateway, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	if err := c.Validate(sf); err != nil {
		return c.JSON(http.StatusBadGateway, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	sessData := data["session"].(map[interface{}]interface{})

	u := &users.Users{
		Username: sessData["username"].(string),
		Password: sf.PasswordOld,
	}

	fmt.Println(u)

	user, err := u.GetUserByUsernamePassword()

	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	if user == nil {
		return c.JSON(http.StatusBadGateway, map[string]interface{}{
			"message": "Wrong username or password",
			"success": false,
		})
	}

	userToBeChanged := users.Users{
		Password: sf.Password,
	}

	if err := userToBeChanged.UpdateWhereID(sessData["id"].(uint64)); err != nil {
		return c.JSON(http.StatusBadGateway, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "",
		"success": true,
	})
}
