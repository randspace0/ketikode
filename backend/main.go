package main

import (
	"html/template"
	"log"
	"net/http"
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
	e.Use(session.Middleware(newSessionStore()))

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

const minSessionSecretLength = 32

// newSessionStore refuses to start with a weak secret, because the cookie only
// holds the user id and is trusted if its signature checks out.
func newSessionStore() *sessions.CookieStore {
	secret := os.Getenv("SESSION_SECRET")
	if len(secret) < minSessionSecretLength {
		log.Fatalf("SESSION_SECRET must be set to at least %d characters (try: openssl rand -hex 32)", minSessionSecretLength)
	}

	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   os.Getenv("COOKIE_SECURE") == "1",
	}
	return store
}
