package simpletrace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var submitBackoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

type Client struct {
	URL    string
	Client http.Client
	Logger *log.Logger
}

// Submit - send the spans to the tracing endpoint synchronously
func (c *Client) Submit(spans ...*Span) error {
	var err error
	var response *http.Response

	body, err := json.Marshal(spans)
	if err != nil {
		return err
	}

	for _, backoff := range submitBackoffSchedule {
		response, err = c.Client.Post(c.URL, "application/json", bytes.NewBuffer(body))
		if err == nil && response.StatusCode == http.StatusAccepted {
			break
		}
		time.Sleep(backoff)
	}

	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("got unexpected status code %+q", response.Status)
	}
	return nil
}

// SubmitAsync - send the spans to the tracing endpoint asynchronous
func (c *Client) SubmitAsync(errBack chan error, spans ...*Span) {
	go func(errBack chan error) {
		errBack <- c.Submit(spans...)
	}(errBack)
}

// NewClient - instantiate a new client with given url
func NewClient(url string) Client {
	var c Client
	c.URL = url
	c.Client = http.Client{Timeout: 1 * time.Second}
	return c
}
