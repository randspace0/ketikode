package admin

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"gitlab.com/cikadev/ketide/handler/adminutil"
	"gitlab.com/cikadev/ketide/repository/codes"
	"gitlab.com/cikadev/ketide/repository/programminglang"
	"gitlab.com/cikadev/ketide/repository/users"
)

type row struct {
	codes.Codes
	Username string
}

func Submissions(c echo.Context) error {
	language := c.QueryParam("language")
	requested, _ := strconv.Atoi(c.QueryParam("page"))

	total, err := (&codes.Codes{Language: language}).Total()
	if err != nil {
		return fail(c, err)
	}
	p := adminutil.NewPager(requested, total)

	found, err := codes.List(language, adminutil.PageSize, (p.Page-1)*adminutil.PageSize)
	if err != nil {
		return fail(c, err)
	}

	ids := make([]uint64, 0, len(found))
	for _, code := range found {
		if code.Owner != 0 {
			ids = append(ids, code.Owner)
		}
	}
	names, err := users.UsernamesByID(ids)
	if err != nil {
		return fail(c, err)
	}

	rows := make([]row, len(found))
	for i, code := range found {
		rows[i] = row{Codes: code, Username: "guest"}
		if name, ok := names[code.Owner]; ok {
			rows[i].Username = name
		}
	}

	err = c.Render(http.StatusOK, "admin/submissions.html", map[string]interface{}{
		"title":     "Admin",
		"rows":      rows,
		"total":     total,
		"pager":     p,
		"language":  language,
		"languages": programminglang.AllList(),
	})
	if err != nil {
		return fail(c, err)
	}
	return nil
}

func fail(c echo.Context, err error) error {
	c.Logger().Error(err)
	return c.String(http.StatusInternalServerError, "Failed to load submissions")
}
