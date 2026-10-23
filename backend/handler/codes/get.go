package codes

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"gitlab.com/cikadev/ketide/handler"
	"gitlab.com/cikadev/ketide/repository/codes"
	"gitlab.com/cikadev/ketide/repository/users"
)

func Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return c.String(http.StatusBadRequest, "Err")
	}

	codeToFind := codes.Codes{
		ID: id,
	}

	foundCode, err := codeToFind.FindByID()

	if foundCode == nil || err != nil {
		return c.Render(http.StatusNotFound, "notfound.html", map[string]interface{}{})
	}

	var foundUser *users.Users
	if foundCode.Owner != 0 {
		userToFind := users.Users{
			ID: foundCode.Owner,
		}

		foundUser, err = userToFind.FindByID()
	} else {
		foundUser = &users.Users{
			ID:       0,
			Username: "guest",
		}
	}

	data := map[string]interface{}{
		"code":  foundCode,
		"owner": foundUser,
		"id":    id,
		"title": id,
	}

	return c.Render(http.StatusOK, "single-code.html", handler.AppendSessionData(c, data))
}

func GetRAW(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return c.String(http.StatusBadRequest, "Err")
	}

	code := codes.Codes{
		ID: id,
	}

	foundCode, err := code.FindByID()

	if foundCode == nil || err != nil {
		return c.String(http.StatusNotFound, "This is a bruh moment. You don't find anything here.")
	}

	return c.String(http.StatusOK, foundCode.Code)
}
