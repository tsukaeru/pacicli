package lib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Status     string
	StatusCode int
	Body       []byte
}

func (r *Response) String() string {
	if r == nil {
		return "(nil)"
	}
	return fmt.Sprintf("Status: %s, Body:\n%s\n", r.Status, r.Body)
}

type Client struct {
	baseURL  string
	username string
	password string
}

func NewClient(baseURL, username, password string) *Client {
	return &Client{
		baseURL:  baseURL,
		username: username,
		password: password,
	}
}

func (c *Client) SendRequest(method, path string, data io.Reader) (*Response, error) {
	debug("Request Method: " + method + ", Path: " + c.baseURL + path)
	if data != nil {
		if b, ok := data.(*bytes.Buffer); ok {
			debug("Raw request:\n" + b.String())
		}
	}

	req, err := http.NewRequest(method, c.baseURL+path, data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/xml")
	req.SetBasicAuth(c.username, c.password)

	cli := http.DefaultClient
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       body,
	}
	debug("Raw response: " + r.String())
	return r, nil
}
