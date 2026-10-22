package main

import (
	"html/template"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo-contrib/session"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/sessions"

	_ "github.com/lib/pq"

	"gitlab.com/cikadev/ketide/handler"
	"gitlab.com/cikadev/ketide/repository"
)

func main() {
	// Database
	MigrateTable(repository.DB)
	MigrateLanguage(repository.DB)

	e := echo.New()

	// Logger
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: `[${time_rfc3339}]  ${status}  ${method} ${host}${path} ${latency_human}` + "\n",
	// }))

	// Session
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Validator = &handler.CustomValidator{Validator: validator.New()}

	// Validator
	e.Static("/static", "static")

	// Template renderer
	renderer := &TemplateRenderer{
		templates: template.Must(findAndParseTemplates("templates", template.FuncMap{})),
	}
	e.Renderer = renderer

	// Router
	route(e)

	e.Logger.Debug(e.Start(":1324"))
}
