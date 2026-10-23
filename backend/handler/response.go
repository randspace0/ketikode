package handler

import (
	"github.com/labstack/echo/v4"
)

// JSONError sends a fixed message to the client and logs the real cause.
func JSONError(c echo.Context, status int, message string, cause error) error {
	if cause != nil {
		c.Logger().Error(cause)
	}
	return c.JSON(status, map[string]interface{}{
		"message": message,
		"success": false,
	})
}

func JSONOK(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "",
		"success": true,
	})
}
