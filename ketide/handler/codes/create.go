package codes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/cikadev/ketide/handler"
	executor "gitlab.com/cikadev/ketide/helpers/ketikode-executor"
	"gitlab.com/cikadev/ketide/repository/codes"
)

type codeForm struct {
	Code     string
	Stdin    string
	Language string `validate:"required"`
}

func getOutput(codesFromUser codes.Codes) codes.Codes {
	code := executor.Code{
		Code:     codesFromUser.Code,
		Language: codesFromUser.Language,
		Stdin:    codesFromUser.Stdin,
	}
	code.Exec()

	codesFromUser.Stdout = code.Stdout
	codesFromUser.Stderr = code.Stderr

	return codesFromUser
}

func Create(c echo.Context) error {
	cf := new(codeForm)
	if err := c.Bind(cf); err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	if err := c.Validate(cf); err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	sessUserId := handler.GetSessionData(c)["id"]
	var userId uint64 = 0
	if sessUserId != nil {
		userId = sessUserId.(uint64)
	}

	code := getOutput(codes.Codes{
		Code:     cf.Code,
		Stdin:    cf.Stdin,
		Language: cf.Language,
		Owner:    userId,
	})

	if err := code.Create(); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/%d", code.ID))
}
