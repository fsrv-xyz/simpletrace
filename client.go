package simpletrace

import (
	"bytes"
	"context"
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
	URL       string
	Client    http.Client
	Logger    *log.Logger
	BasicAuth ClientAuth
}

type ClientAuth struct {
	Enable   bool
	Username string
	Password string
}

// Submit - send the spans to the tracing endpoint synchronously
func (c *Client) Submit(spans ...*Span) error {
	var err error
	var response *http.Response

	body, err := json.Marshal(spans)
	if err != nil {
		return err
	}
	// build request
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	request, requestBuilderError := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewBuffer(body))
	if requestBuilderError != nil {
		return err
	}

	request.Header.Set("content-type", "application/json")

	// set basic auth headers if enabled
	if c.BasicAuth.Enable {
		request.SetBasicAuth(c.BasicAuth.Username, c.BasicAuth.Password)
	}

	for _, backoff := range submitBackoffSchedule {
		// execute http request; stop if request is accepted on tracing server
		response, err = c.Client.Do(request)
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
	c.Client = http.Client{}
	return c
}
