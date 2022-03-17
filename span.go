package simpletrace

import (
	"time"
)

type Kind string

const (
	KindClient   Kind = "CLIENT"
	KindServer   Kind = "SERVER"
	KindProducer Kind = "PRODUCER"
	KindConsumer Kind = "CONSUMER"
)

type Span struct {
	Id             string            `json:"id"`
	TraceId        string            `json:"traceId"`
	ParentId       string            `json:"parentId,omitempty"`
	Kind           Kind              `json:"kind,omitempty"`
	Timestamp      int64             `json:"timestamp"`
	Duration       int               `json:"duration"`
	Name           string            `json:"name,omitempty"`
	Tags           map[string]string `json:"tags,omitempty"`
	Shared         bool              `json:"shared"`
	LocalEndpoint  Service           `json:"localEndpoint"`
	RemoteEndpoint Service           `json:"remoteEndpoint"`
	Annotations    []Annotation      `json:"annotations"`

	StartTime time.Time `json:"-"`
}

type Service struct {
	ServiceName string `json:"serviceName"`
	IPv4        string `json:"ipv4,omitempty"`
	IPv6        string `json:"ipv6,omitempty"`
	Port        int    `json:"port,omitempty"`
}

type Annotation struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}

func (s *Span) Submit(c *Client) error {
	return c.Submit(*s)
}

func (s *Span) AddAnnotation(timestamp time.Time, value string) {
	s.Annotations = append(s.Annotations, Annotation{
		Timestamp: timestamp.UnixMicro(),
		Value:     value,
	})
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
