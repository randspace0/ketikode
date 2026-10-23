package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/handler/admin"
	"github.com/randspace0/ketikode/backend/handler/adminutil"
	"github.com/randspace0/ketikode/backend/handler/codes"
	"github.com/randspace0/ketikode/backend/handler/users"
)

func route(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", handler.AppendSessionData(c, map[string]interface{}{}))
	})

	e.GET("/help", func(c echo.Context) error {
		return c.Render(http.StatusOK, "help.html", handler.AppendSessionData(c, map[string]interface{}{}))
	})

	e.GET("/stats", handler.StatsHandler)

	e.GET("/myrecent", handler.RequiresSignin(users.RecentCodesHandler))

	e.GET("/settings", handler.RequiresSignin(users.SettingsHandler))
	e.POST("/settings", handler.RequiresSignin(users.SettingsProcessHandler))

	e.GET("/signin", handler.RequiresSignout(users.SigninHandler))
	e.POST("/signin", handler.RequiresSignout(users.SigninProcessHandler), handler.RateLimit(0.2, 5))

	e.GET("/signup", handler.RequiresSignout(users.SignupHandler))
	e.POST("/signup", handler.RequiresSignout(users.SignupProcessHandler), handler.RateLimit(0.2, 5))

	e.GET("/signout", handler.RequiresSignin(func(c echo.Context) error {
		handler.Signout(c)
		return c.Redirect(http.StatusTemporaryRedirect, "/signin")
	}))

	// Each run holds an executor sandbox for up to several seconds.
	e.POST("/code/create", codes.Create, handler.RateLimit(0.5, 10))
	e.DELETE("/code/delete", handler.RequiresSignin(codes.Delete))

	if user, password, ok := adminutil.Credentials(); ok {
		e.GET("/admin", admin.Submissions, adminutil.Auth(user, password))
	}

	e.GET("/:id", codes.Get)
	e.GET("/raw/:id", codes.GetRAW)
}
