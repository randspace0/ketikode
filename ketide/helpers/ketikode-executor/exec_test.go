package executor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExec(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/exec" || r.PostFormValue("language") != "cpp" || r.PostFormValue("stdin") != "in" {
			t.Errorf("unexpected request %s %v", r.URL.Path, r.PostForm)
		}
		w.Write([]byte(`{"stdout":"out","stderr":"err","success":true}`))
	}))
	defer srv.Close()
	t.Setenv("EXECUTOR_URL", srv.URL)

	c := Code{Language: "cpp", Code: "x", Stdin: "in"}
	c.Exec()
	if c.Stdout != "out" || c.Stderr != "err" {
		t.Fatalf("%+v", c)
	}
}

func TestExecMessageFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success":false,"message":"Invalid language"}`))
	}))
	defer srv.Close()
	t.Setenv("EXECUTOR_URL", srv.URL)

	c := Code{Language: "x"}
	c.Exec()
	if c.Stderr != "Invalid language" {
		t.Fatalf("%+v", c)
	}
}

func TestExecUnreachable(t *testing.T) {
	t.Setenv("EXECUTOR_URL", "http://127.0.0.1:1")

	c := Code{Language: "cpp"}
	c.Exec() // must not panic
	if c.Stderr == "" {
		t.Fatal("want an error in Stderr")
	}
}
