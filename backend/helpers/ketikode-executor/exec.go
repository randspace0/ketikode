package executor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Code struct {
	Language string
	Code     string
	Stdin    string
	Stdout   string
	Stderr   string
}

type response struct {
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

var client = &http.Client{Timeout: 60 * time.Second}

func executorURL() string {
	if u := os.Getenv("EXECUTOR_URL"); u != "" {
		return u
	}
	return "http://localhost:8001"
}

// Exec runs the code on the executor service. Failures to reach the executor
// are reported through Stderr so the user sees them instead of a crash.
func (c *Code) Exec() {
	data := url.Values{
		"code":     {c.Code},
		"language": {c.Language},
		"stdin":    {c.Stdin},
	}

	resp, err := client.PostForm(executorURL()+"/api/exec", data)
	if err != nil {
		c.Stderr = fmt.Sprintf("executor unavailable: %v", err)
		return
	}
	defer resp.Body.Close()

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		c.Stderr = "executor returned an invalid response"
		return
	}

	c.Stdout = r.Stdout
	c.Stderr = r.Stderr
	if r.Message != "" && r.Stderr == "" {
		c.Stderr = r.Message
	}
}
