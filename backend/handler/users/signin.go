package users

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/cikadev/ketide/handler"
	"gitlab.com/cikadev/ketide/repository/users"
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

	u := &users.Users{
		Username: sf.Username,
		Password: sf.Password,
	}

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

	handler.Signin(c, user.ID)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "",
		"success": true,
	})
}
