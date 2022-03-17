package simpletrace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

type Client struct {
	URL    string
	Client http.Client
}

func (c *Client) Submit(spans ...Span) error {
	var err error
	var response *http.Response

	body, _ := json.Marshal(spans)
	fmt.Println(string(body))

	for _, backoff := range backoffSchedule {
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

func NewClient(url string) Client {
	var c Client
	c.URL = url
	c.Client = http.Client{Timeout: 1 * time.Second}
	return c
}
