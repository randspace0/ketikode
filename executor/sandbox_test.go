package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain lets the test binary stand in for the isolate binary: Run is
// pointed at os.Args[0] with FAKE_ISOLATE set, so the whole Run flow is
// exercised without isolate installed.
func TestMain(m *testing.M) {
	if os.Getenv("FAKE_ISOLATE") != "" {
		os.Exit(fakeIsolate(os.Args[1:]))
	}
	os.Exit(m.Run())
}

func fakeIsolate(args []string) int {
	root := os.Getenv("FAKE_ROOT")
	var meta string
	var cmd []string
	for i, a := range args {
		if a == "--" {
			cmd = args[i+1:]
			break
		}
		if v, ok := strings.CutPrefix(a, "--meta="); ok {
			meta = v
		}
	}
	box := filepath.Join(root, "box")
	switch {
	case slicesContain(args, "--init"):
		os.MkdirAll(box, 0o755)
		os.Stdout.WriteString(root + "\n")
	case slicesContain(args, "--run"):
		isCompile := strings.HasSuffix(cmd[0], "g++")
		m, out, errOut := "time:0.01\nmax-rss:100\n", "hello", ""
		switch os.Getenv("FAKE_SCENARIO") {
		case "compile_error":
			if isCompile {
				m, out, errOut = "status:RE\nexit-code:1\n", "", "boom"
			}
		case "timeout":
			if !isCompile {
				m = "status:TO\nmessage:Time limit exceeded\n"
			}
		case "internal":
			m = "status:XX\nmessage:box exploded\n"
		}
		os.WriteFile(meta, []byte(m), 0o644)
		os.WriteFile(filepath.Join(box, "stdout.txt"), []byte(out), 0o644)
		os.WriteFile(filepath.Join(box, "stderr.txt"), []byte(errOut), 0o644)
	}
	return 0
}

func slicesContain(a []string, s string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

func fakeRun(t *testing.T, scenario string, req Request) (*Response, error, string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("FAKE_ISOLATE", "1")
	t.Setenv("FAKE_ROOT", root)
	t.Setenv("FAKE_SCENARIO", scenario)
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	res, err := NewIsolate(exe, false, 1).Run(context.Background(), req)
	return res, err, root
}

func TestRunOK(t *testing.T) {
	res, err, root := fakeRun(t, "ok", Request{Language: "cpp", Code: "int main(){}", Stdin: "in", TimeSec: 1, MemoryMB: 64})
	if err != nil {
		t.Fatal(err)
	}
	if res.Phase != "run" || !res.Success || res.Stdout != "hello" || res.MemoryKB != 100 || res.TimeSec != 0.01 {
		t.Fatalf("%+v", res)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "box", "main.cpp")); string(b) != "int main(){}" {
		t.Errorf("source not written to box: %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "box", "stdin.txt")); string(b) != "in" {
		t.Errorf("stdin not written to box: %q", b)
	}
}

func TestRunCompileError(t *testing.T) {
	res, err, _ := fakeRun(t, "compile_error", Request{Language: "cpp", Code: "x", TimeSec: 1, MemoryMB: 64})
	if err != nil {
		t.Fatal(err)
	}
	if res.Phase != "compile" || res.Success || res.Stderr != "boom" || res.Stdout != "" {
		t.Fatalf("%+v", res)
	}
}

func TestRunTimeout(t *testing.T) {
	res, err, _ := fakeRun(t, "timeout", Request{Language: "python3", Code: "x", TimeSec: 1, MemoryMB: 64})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "TO" || res.Success {
		t.Fatalf("%+v", res)
	}
}

func TestRunInternalError(t *testing.T) {
	if _, err, _ := fakeRun(t, "internal", Request{Language: "python3", Code: "x", TimeSec: 1, MemoryMB: 64}); err == nil {
		t.Fatal("want error for status XX")
	}
}

func TestRunUnknownLanguage(t *testing.T) {
	if _, err := NewIsolate("isolate", false, 1).Run(context.Background(), Request{Language: "cobol"}); err == nil {
		t.Fatal("want error")
	}
}

func TestRunCommandMemory(t *testing.T) {
	got := strings.Join(languages["java"].runCommand(128), " ")
	if !strings.Contains(got, "-Xmx128m") {
		t.Fatal(got)
	}
}

// Real-sandbox smoke test; needs isolate installed and root (see Dockerfile).
func TestIsolateIntegration(t *testing.T) {
	if os.Getenv("EXECUTOR_INTEGRATION") == "" {
		t.Skip("set EXECUTOR_INTEGRATION=1 to run against a real isolate")
	}
	if _, err := exec.LookPath("isolate"); err != nil {
		t.Skip("isolate not installed")
	}
	s := NewIsolate("isolate", os.Getenv("ISOLATE_CGROUPS") == "1", 1)
	run := func(lang, code, stdin string) *Response {
		t.Helper()
		res, err := s.Run(context.Background(), Request{Language: lang, Code: code, Stdin: stdin, TimeSec: 2, MemoryMB: 128})
		if err != nil {
			t.Fatal(err)
		}
		return res
	}

	if r := run("python3", "print(input()[::-1])", "abc"); !r.Success || strings.TrimSpace(r.Stdout) != "cba" {
		t.Errorf("python: %+v", r)
	}
	if r := run("cpp", "int main(){ not c++ }", ""); r.Phase != "compile" || r.Success || r.Stderr == "" {
		t.Errorf("compile error: %+v", r)
	}
	if r := run("python3", "while True: pass", ""); r.Status != "TO" {
		t.Errorf("timeout: %+v", r)
	}
	if r := run("python3", "import urllib.request as u; u.urlopen('http://1.1.1.1', timeout=2)", ""); r.Success {
		t.Errorf("network must be blocked: %+v", r)
	}
}
