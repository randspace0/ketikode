package adminutil

import (
	"crypto/subtle"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const PageSize = 50

// Credentials returns the backoffice login from ADMIN_USER / ADMIN_PASSWORD.
// The backoffice stays disabled (routes are not registered) unless both are set.
func Credentials() (user, password string, ok bool) {
	user, password = os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD")
	return user, password, user != "" && password != ""
}

func Auth(user, password string) echo.MiddlewareFunc {
	return middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Realm: "ketide admin",
		Validator: func(u, p string, _ echo.Context) (bool, error) {
			userOK := subtle.ConstantTimeCompare([]byte(u), []byte(user)) == 1
			passOK := subtle.ConstantTimeCompare([]byte(p), []byte(password)) == 1
			return userOK && passOK, nil
		},
	})
}

// Pager describes one page of a paginated list.
type Pager struct {
	Page, Pages int
	Prev, Next  int // 0 when there is no such page
}

// NewPager clamps page into range for total items.
func NewPager(page int, total int64) Pager {
	pages := int((total + PageSize - 1) / PageSize)
	if pages < 1 {
		pages = 1
	}
	if page > pages {
		page = pages
	}
	if page < 1 {
		page = 1
	}

	p := Pager{Page: page, Pages: pages}
	if page > 1 {
		p.Prev = page - 1
	}
	if page < pages {
		p.Next = page + 1
	}
	return p
}
