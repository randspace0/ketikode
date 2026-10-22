package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type fakeRunner struct{ got Request }

func (f *fakeRunner) Run(_ context.Context, r Request) (*Response, error) {
	f.got = r
	return &Response{Phase: "run", Stdout: "hi", Success: true, Status: "OK"}, nil
}

func post(h http.Handler, form url.Values) *httptest.ResponseRecorder {
	r := httptest.NewRequest("POST", "/api/exec", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestExecClampsLimits(t *testing.T) {
	f := &fakeRunner{}
	w := post(newServer(f), url.Values{"language": {"cpp"}, "code": {"x"}, "timeout": {"999"}, "memory": {"-1"}})
	if w.Code != 200 || f.got.TimeSec != maxTimeSec || f.got.MemoryMB != defaultMemoryMB {
		t.Fatalf("code=%d got=%+v", w.Code, f.got)
	}
}

func TestExecRejectsBadInput(t *testing.T) {
	h := newServer(&fakeRunner{})
	for _, form := range []url.Values{
		{"language": {"cobol"}, "code": {"x"}},
		{"language": {"cpp"}},
	} {
		if w := post(h, form); w.Code != 400 {
			t.Errorf("%v: code=%d", form, w.Code)
		}
	}
}

func TestParseMeta(t *testing.T) {
	m := parseMeta(strings.NewReader("time:0.012\nstatus:TO\nmessage:Time limit exceeded\nmax-rss:3000\nexit-code:0\n"))
	r := responseFromMeta(m)
	if r.Status != "TO" || r.Success || r.TimeSec != 0.012 || r.MemoryKB != 3000 {
		t.Fatalf("%+v", r)
	}
	if !responseFromMeta(map[string]string{}).Success {
		t.Fatal("empty status should be OK")
	}
}

func TestRunArgs(t *testing.T) {
	s := &Isolate{Cgroups: true}
	a := s.runArgs(3, limits{timeSec: 2, memMB: 64, procs: 1, fsizeKB: 1024}, "m", "stdin.txt", []string{"./main"})
	for _, want := range []string{"--box-id=3", "--cg", "--cg-mem=65536", "--time=2", "--wall-time=5", "--stdin=stdin.txt"} {
		if !slices.Contains(a, want) {
			t.Errorf("missing %s in %v", want, a)
		}
	}
	if tail := a[len(a)-2:]; !reflect.DeepEqual(tail, []string{"--", "./main"}) {
		t.Errorf("bad tail %v", tail)
	}
	if slices.Contains(a, "--share-net") {
		t.Error("network must stay disabled")
	}
	s.Cgroups = false
	if a = s.runArgs(0, limits{timeSec: 1, memMB: 8, procs: 1, fsizeKB: 1, asLimit: true}, "m", "", nil); !slices.Contains(a, "--mem=8192") {
		t.Errorf("want address-space limit: %v", a)
	}
}
