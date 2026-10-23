package codes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/repository/codes"
)

type deleteForm struct {
	ID string `validate:"required"`
}

// TODO implement AJAX

func Delete(c echo.Context) error {
	df := new(deleteForm)

	if err := c.Bind(df); err != nil {
		return handler.JSONError(c, http.StatusBadRequest, "Invalid request", err)
	}

	if err := c.Validate(df); err != nil {
		return handler.JSONError(c, http.StatusBadRequest, "An id is required", nil)
	}

	id, err := strconv.ParseUint(df.ID, 10, 64)
	if err != nil {
		return handler.JSONError(c, http.StatusBadRequest, "Invalid id", nil)
	}

	// RequiresSignin has already checked the session user exists.
	code := codes.Codes{
		ID:    id,
		Owner: handler.GetSessionData(c)["id"].(uint64),
	}

	if err := code.Delete(); err != nil {
		if errors.Is(err, codes.ErrNotFound) {
			return handler.JSONError(c, http.StatusNotFound, "Code not found", nil)
		}
		return handler.JSONError(c, http.StatusInternalServerError, "Could not delete the code", err)
	}

	return handler.JSONOK(c)
}
