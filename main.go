package simpletrace

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

func randomID(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Println(err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

type Span struct {
	Id            string            `json:"id"`
	TraceId       string            `json:"traceId"`
	ParentId      string            `json:"parentId,omitempty"`
	Timestamp     int64             `json:"timestamp"`
	Duration      int               `json:"duration"`
	Name          string            `json:"name,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	LocalEndpoint struct {
		ServiceName string `json:"serviceName"`
	} `json:"localEndpoint"`
	StartTime time.Time `json:"-"`
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

func (s *Span) Submit(c *Client) {
	c.Submit(*s)
}

type Client struct {
	URL    string
	Client http.Client
}

func NewSpan(name string) Span {
	span := Span{
		Id:      randomID(8),
		TraceId: randomID(8),
		Name:    name,
	}
	span.Tags = make(map[string]string)
	return span
}

func (s *Span) NewChildSpan(name string) Span {
	sub := NewSpan(name)
	sub.TraceId = s.TraceId
	sub.ParentId = s.Id
	return sub
}

func (s *Span) Start() {
	s.StartTime = time.Now()
	s.Timestamp = time.Now().UnixMicro()
}

func (s *Span) Finalize() *Span {
	s.Duration = int(time.Since(s.StartTime).Microseconds())
	return s
}

func NewClient(url string) Client {
	var c Client
	c.URL = url
	c.Client = http.Client{Timeout: 1 * time.Second}
	return c
}
