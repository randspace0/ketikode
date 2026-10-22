package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/cikadev/ketide/handler"
	"gitlab.com/cikadev/ketide/handler/codes"
	"gitlab.com/cikadev/ketide/handler/users"
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
	e.POST("/signin", handler.RequiresSignout(users.SigninProcessHandler))

	e.GET("/signup", handler.RequiresSignout(users.SignupHandler))
	e.POST("/signup", handler.RequiresSignout(users.SignupProcessHandler))

	e.GET("/signout", handler.RequiresSignin(func(c echo.Context) error {
		handler.Signout(c)
		return c.Redirect(http.StatusTemporaryRedirect, "/signin")
	}))

	e.POST("/code/create", codes.Create)
	e.DELETE("/code/delete", handler.RequiresSignin(codes.Delete))

	e.GET("/:id", codes.Get)
	e.GET("/raw/:id", codes.GetRAW)
}
