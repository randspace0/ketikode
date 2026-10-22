package users

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/cikadev/ketide/handler"
	"gitlab.com/cikadev/ketide/repository/codes"
)

func RecentCodesHandler(c echo.Context) error {
	data := handler.AppendSessionData(c, map[string]interface{}{})

	sessData := data["session"].(map[interface{}]interface{})
	u := codes.Codes{
		Owner: sessData["id"].(uint64),
	}

	codes, err := u.FindAllOwnedCodesByUserID()

	if err == nil {
		data["codes"] = codes
	}

	return c.Render(http.StatusOK, "code.html", data)
}
