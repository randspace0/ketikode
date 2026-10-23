package main

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/randspace0/ketikode/backend/repository/codes"
	"github.com/randspace0/ketikode/backend/repository/programminglang"
	"github.com/randspace0/ketikode/backend/repository/users"
)

func TestEveryPageRenders(t *testing.T) {
	tpl, err := findAndParseTemplates("templates", template.FuncMap{})
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	code := &codes.Codes{ID: 7, Language: "cpp", Code: "<b>x</b>", Stdin: "in", Stdout: "out", Stderr: "err", Owner: 1, CreatedAt: now, UpdatedAt: now}
	owner := &users.Users{ID: 1, Username: "alice"}
	session := map[interface{}]interface{}{"id": uint64(1), "username": "alice"}

	pages := map[string]map[string]interface{}{
		"home.html":              {},
		"help.html":              {},
		"notfound.html":          {},
		"stats.html":             {"languages": map[string]int64{"C++": 3, "Python 3": 1}},
		"code.html":              {"codes": []codes.Codes{*code}},
		"single-code.html":       {"code": code, "owner": owner, "id": uint64(7)},
		"user/signin.html":       {},
		"user/signup.html":       {},
		"user/settings.html":     {},
		"admin/submissions.html": {"rows": []struct{}{}, "total": int64(0), "pager": struct{ Page, Pages, Prev, Next int }{1, 1, 0, 0}, "language": "", "languages": []programminglang.ProgrammingLang{{LanguageID: "cpp", Name: "C++"}}},
	}

	for name, data := range pages {
		for _, signedIn := range []bool{false, true} {
			d := map[string]interface{}{"reverse": nil}
			for k, v := range data {
				d[k] = v
			}
			if signedIn {
				d["session"] = session
			}
			var out bytes.Buffer
			if err := tpl.ExecuteTemplate(&out, name, d); err != nil {
				t.Errorf("%s (signed in: %v): %v", name, signedIn, err)
				continue
			}
			html := out.String()
			for _, want := range []string{`<html lang="en">`, `id="main"`} {
				if !strings.Contains(html, want) {
					t.Errorf("%s: missing %s", name, want)
				}
			}
		}
	}
}

func TestCodeIsEscaped(t *testing.T) {
	tpl, err := findAndParseTemplates("templates", template.FuncMap{})
	if err != nil {
		t.Fatal(err)
	}
	code := &codes.Codes{ID: 1, Code: "<script>alert(1)</script>", Stdout: "<img src=x onerror=alert(2)>", UpdatedAt: time.Now()}
	var out bytes.Buffer
	data := map[string]interface{}{"code": code, "owner": &users.Users{Username: "a"}, "id": uint64(1)}
	if err := tpl.ExecuteTemplate(&out, "single-code.html", data); err != nil {
		t.Fatal(err)
	}
	if html := out.String(); strings.Contains(html, "<script>alert(1)") || strings.Contains(html, "<img src=x") {
		t.Fatal("user code must be HTML-escaped on the code page")
	}
}
