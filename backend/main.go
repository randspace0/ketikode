package main

import (
	"html/template"
	"log"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo-contrib/session"

	"github.com/go-playground/validator/v10"

	"github.com/gorilla/sessions"

	"github.com/randspace0/ketikode/backend/handler"
	"github.com/randspace0/ketikode/backend/repository"
	"github.com/randspace0/ketikode/backend/repository/migrate"
)

func main() {
	// Database
	if err := repository.Open(); err != nil {
		log.Fatal(err)
	}
	if err := migrate.Run(repository.DB); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	e := echo.New()

	// Logger
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: `[${time_rfc3339}]  ${status}  ${method} ${host}${path} ${latency_human}` + "\n",
	// }))

	// Session
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret()))))

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

func sessionSecret() string {
	if s := os.Getenv("SESSION_SECRET"); s != "" {
		return s
	}
	return "secret"
}
