package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const maxOutput = 64 << 10

func itoa(i int) string { return strconv.Itoa(i) }

type Request struct {
	Language string
	Code     string
	Stdin    string
	TimeSec  int
	MemoryMB int
}

// Response.Phase is "compile" when compilation failed, otherwise "run".
type Response struct {
	Phase    string  `json:"phase"`
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	Success  bool    `json:"success"`
	Status   string  `json:"status"` // OK, RE (non-zero exit), SG (signal), TO (timeout)
	ExitCode int     `json:"exit_code"`
	TimeSec  float64 `json:"time"`
	MemoryKB int     `json:"memory"`
	Message  string  `json:"message,omitempty"`
}

type Runner interface {
	Run(ctx context.Context, req Request) (*Response, error)
}

type Isolate struct {
	Bin     string
	Cgroups bool
	boxes   chan int
}

func NewIsolate(bin string, cgroups bool, boxes int) *Isolate {
	s := &Isolate{Bin: bin, Cgroups: cgroups, boxes: make(chan int, boxes)}
	for i := 0; i < boxes; i++ {
		s.boxes <- i
	}
	return s
}

type limits struct {
	timeSec, memMB, procs, fsizeKB int
	asLimit                        bool
}

func (s *Isolate) base(box int) []string {
	a := []string{"--box-id=" + itoa(box)}
	if s.Cgroups {
		a = append(a, "--cg")
	}
	return a
}

// runArgs builds the `isolate --run` argument list. Network stays disabled
// because --share-net is never passed.
func (s *Isolate) runArgs(box int, l limits, meta, stdin string, cmd []string) []string {
	a := append(s.base(box), "--run", "--silent", "--meta="+meta,
		"--time="+itoa(l.timeSec), "--wall-time="+itoa(l.timeSec*2+1), "--extra-time=1",
		"--fsize="+itoa(l.fsizeKB), "--processes="+itoa(l.procs),
		"--env=PATH=/usr/local/bin:/usr/bin:/bin", "--env=HOME=/box",
		"--stdout=stdout.txt", "--stderr=stderr.txt")
	if stdin != "" {
		a = append(a, "--stdin="+stdin)
	}
	switch {
	case s.Cgroups:
		a = append(a, "--cg-mem="+itoa(l.memMB<<10))
	case l.asLimit:
		a = append(a, "--mem="+itoa(l.memMB<<10))
	}
	return append(append(a, "--"), cmd...)
}

func (s *Isolate) Run(ctx context.Context, req Request) (*Response, error) {
	lang, ok := languages[req.Language]
	if !ok {
		return nil, errors.New("invalid language")
	}

	var box int
	select {
	case box = <-s.boxes:
		defer func() { s.boxes <- box }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Clear anything left behind by a crashed run before reusing the box.
	s.isolate(ctx, append(s.base(box), "--cleanup"))
	out, err := s.isolate(ctx, append(s.base(box), "--init"))
	if err != nil {
		return nil, fmt.Errorf("isolate --init: %w: %s", err, out)
	}
	defer s.isolate(context.WithoutCancel(ctx), append(s.base(box), "--cleanup"))
	dir := filepath.Join(strings.TrimSpace(string(out)), "box")

	if err := os.WriteFile(filepath.Join(dir, lang.Source), []byte(req.Code), 0o644); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "stdin.txt"), []byte(req.Stdin), 0o644); err != nil {
		return nil, err
	}

	if lang.Compile != nil {
		res, err := s.step(ctx, box, dir, limits{timeSec: 15, memMB: 1024, procs: 64, fsizeKB: 256 << 10}, "", lang.Compile)
		if err != nil {
			return nil, err
		}
		if res.Status != "OK" {
			res.Phase = "compile"
			res.Stderr = res.Stdout + res.Stderr
			res.Stdout = ""
			return res, nil
		}
	}

	procs := 1
	if lang.VM {
		procs = 64
	}
	res, err := s.step(ctx, box, dir, limits{timeSec: req.TimeSec, memMB: req.MemoryMB, procs: procs, fsizeKB: 1024, asLimit: !lang.VM}, "stdin.txt", lang.runCommand(req.MemoryMB))
	if err != nil {
		return nil, err
	}
	res.Phase = "run"
	return res, nil
}

func (s *Isolate) step(ctx context.Context, box int, dir string, l limits, stdin string, cmd []string) (*Response, error) {
	meta := filepath.Join(dir, "..", "meta.txt")
	os.Remove(meta)
	out, err := s.isolate(ctx, s.runArgs(box, l, meta, stdin, cmd))
	m, merr := readMeta(meta)
	if merr != nil {
		return nil, fmt.Errorf("isolate --run: %w: %s", err, out)
	}
	if m["status"] == "XX" {
		return nil, fmt.Errorf("isolate internal error: %s", m["message"])
	}

	res := responseFromMeta(m)
	res.Stdout = readCapped(filepath.Join(dir, "stdout.txt"))
	res.Stderr = readCapped(filepath.Join(dir, "stderr.txt"))
	os.Remove(filepath.Join(dir, "stdout.txt"))
	os.Remove(filepath.Join(dir, "stderr.txt"))
	return res, nil
}

func (s *Isolate) isolate(ctx context.Context, args []string) ([]byte, error) {
	return exec.CommandContext(ctx, s.Bin, args...).CombinedOutput()
}

func readMeta(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseMeta(f), nil
}

func parseMeta(r io.Reader) map[string]string {
	m := map[string]string{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if k, v, ok := strings.Cut(sc.Text(), ":"); ok {
			m[k] = v
		}
	}
	return m
}

func responseFromMeta(m map[string]string) *Response {
	r := &Response{Status: m["status"], Message: m["message"]}
	if r.Status == "" {
		r.Status = "OK"
	}
	r.ExitCode, _ = strconv.Atoi(m["exit-code"])
	r.TimeSec, _ = strconv.ParseFloat(m["time"], 64)
	if v, ok := m["cg-mem"]; ok {
		r.MemoryKB, _ = strconv.Atoi(v)
	} else {
		r.MemoryKB, _ = strconv.Atoi(m["max-rss"])
	}
	r.Success = r.Status == "OK"
	return r
}

func readCapped(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	b, _ := io.ReadAll(io.LimitReader(f, maxOutput))
	return strings.ToValidUTF8(string(b), "�")
}
