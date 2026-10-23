package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			t, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

func main() {
	e := echo.New()

	e.Static("/static", "static")

	// renderer := &TemplateRenderer{
	// 	 templates: template.Must(template.ParseGlob("templates/*.html")),
	// }
	renderer := &TemplateRenderer{
		templates: template.Must(findAndParseTemplates("templates", template.FuncMap{})),
	}
	// renderer, err := findAndParseTemplates("templates", a)
	e.Renderer = renderer

	// Main

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", map[string]interface{}{})
	})

	e.GET("/help", func(c echo.Context) error {
		return c.Render(http.StatusOK, "help.html", map[string]interface{}{})
	})

	// User

	e.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, "user/signup.html", map[string]interface{}{})
	})

	e.GET("/signin", func(c echo.Context) error {
		return c.Render(http.StatusOK, "user/signin.html", map[string]interface{}{})
	})

	e.GET("/settings", func(c echo.Context) error {
		return c.Render(http.StatusOK, "user/settings.html", map[string]interface{}{})
	})

	e.GET("/user/:username", func(c echo.Context) error {
		return c.Render(http.StatusOK, "user/profile.html", map[string]interface{}{})
	})

	e.GET("/problems/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "problems.html", map[string]interface{}{})
	})

	e.GET("/problems/:id", func(c echo.Context) error {
		return c.Render(http.StatusOK, "single-problems.html", map[string]interface{}{})
	})

	e.GET("/playground/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "playground.html", map[string]interface{}{})
	})

	e.Logger.Debug(e.Start(":1323"))
}
