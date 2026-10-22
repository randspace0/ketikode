package users

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/cikadev/ketide/handler"
	"gitlab.com/cikadev/ketide/repository/users"
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
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
}

func SignupProcessHandler(c echo.Context) error {
	usf := new(signupForm)
	if err := c.Bind(usf); err != nil {
		return c.JSON(http.StatusBadGateway, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	if err := c.Validate(usf); err != nil {
		return c.JSON(http.StatusBadGateway, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	u := users.Users{
		Username: usf.Username,
		Email:    usf.Email,
		Password: usf.Password,
	}

	if err := u.Create(); err != nil {
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
