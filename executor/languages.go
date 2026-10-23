package main

import "strings"

// Language describes how to compile and run one language inside an isolate box.
// "{mem}" in Run is replaced with the memory limit in MB.
type Language struct {
	Source  string
	Compile []string // nil for interpreted languages
	Run     []string
	// VM languages cannot run under an address-space limit, so they size
	// their own heap (via {mem}) instead.
	VM bool
	// CompileAS caps the compiler's address space when cgroups are off. Only
	// for compilers that stay within it (gcc); javac and rustc reserve far more.
	CompileAS bool
}

var languages = map[string]Language{
	"c":       {Source: "main.c", Compile: []string{"/usr/bin/gcc", "-O2", "-o", "main", "main.c", "-lm"}, Run: []string{"./main"}, CompileAS: true},
	"cpp":     {Source: "main.cpp", Compile: []string{"/usr/bin/g++", "-O2", "-o", "main", "main.cpp"}, Run: []string{"./main"}, CompileAS: true},
	"python3": {Source: "main.py", Run: []string{"/usr/bin/python3", "main.py"}},
	"ruby":    {Source: "main.rb", Run: []string{"/usr/bin/ruby", "main.rb"}},
	"rust":    {Source: "main.rs", Compile: []string{"/usr/bin/rustc", "-O", "-o", "main", "main.rs"}, Run: []string{"./main"}},
	"java": {
		Source:  "Submission.java",
		Compile: []string{"/usr/bin/javac", "Submission.java"},
		Run:     []string{"/usr/bin/java", "-Xmx{mem}m", "-Xss64m", "-cp", ".", "Submission"},
		VM:      true,
	},
}

func (l Language) runCommand(memMB int) []string {
	cmd := make([]string, len(l.Run))
	for i, a := range l.Run {
		cmd[i] = strings.ReplaceAll(a, "{mem}", itoa(memMB))
	}
	return cmd
}
