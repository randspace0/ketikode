package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every /static/... file a template or script points at must exist, so a typo or
// a missing vendored file shows up here instead of as a blank page.
func TestStaticReferencesExist(t *testing.T) {
	ref := regexp.MustCompile(`/static/[A-Za-z0-9_./-]+\.(?:js|css|png|ico)`)
	checked := 0

	filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}
		b, _ := os.ReadFile(path)
		for _, m := range ref.FindAllString(string(b), -1) {
			checked++
			if _, err := os.Stat(strings.TrimPrefix(m, "/")); err != nil {
				t.Errorf("%s references %s, which does not exist", path, m)
			}
		}
		return nil
	})

	if checked == 0 {
		t.Fatal("found no static references; the check is not looking at anything")
	}
}
