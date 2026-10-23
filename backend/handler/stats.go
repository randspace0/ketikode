package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/randspace0/ketikode/backend/repository/codes"
	"github.com/randspace0/ketikode/backend/repository/programminglang"
)

func StatsHandler(c echo.Context) error {
	data := AppendSessionData(c, map[string]interface{}{})

	languages := map[string]int64{}
	for _, language := range programminglang.AllList() {
		c := codes.Codes{
			Language: language.LanguageID,
		}
		total, _ := c.Total()
		languages[language.Name] = total
	}

	data["languages"] = languages

	return c.Render(http.StatusOK, "stats.html", data)
}
