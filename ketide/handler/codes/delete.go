package codes

import (
	"net/http"
	"strconv"

	"gitlab.com/cikadev/ketide/handler"

	"github.com/labstack/echo/v4"

	"gitlab.com/cikadev/ketide/repository/codes"
)

type deleteForm struct {
	ID string `validate:"required"`
}

// TODO implement AJAX

func Delete(c echo.Context) error {
	data := handler.AppendSessionData(c, map[string]interface{}{})

	sessionData := data["session"].(map[interface{}]interface{})
	userID := sessionData["id"]

	df := new(deleteForm)

	if err := c.Bind(df); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	if err := c.Validate(df); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	id, err := strconv.ParseUint(df.ID, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	code := codes.Codes{
		ID:    id,
		Owner: userID.(uint64),
	}

	if err := code.Delete(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "",
		"success": true,
	})
}
