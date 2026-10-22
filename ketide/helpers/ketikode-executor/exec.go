package executor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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
	Success bool   `json:"success"`
}

func (c *Code) Exec() {
	data := url.Values{
		"code":     {c.Code},
		"language": {c.Language},
		"stdin":    {c.Stdin},
	}

	// TODO handle error
	resp, _ := http.PostForm("http://localhost:8001/api/exec", data)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	responseResult := response{}
	json.Unmarshal([]byte(body), &responseResult)

	c.Stdout = responseResult.Stdout
	c.Stderr = responseResult.Stderr
}
