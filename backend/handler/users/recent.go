package users

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/repository/codes"
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
