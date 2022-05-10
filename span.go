package simpletrace

import (
	"encoding/json"
	"encoding/xml"
	"net"
	"sync"
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
	SpanId         string            `json:"id"`
	TraceId        string            `json:"traceId"`
	ParentSpanId   string            `json:"parentId,omitempty"`
	Kind           Kind              `json:"kind,omitempty"`
	Timestamp      int64             `json:"timestamp"`
	Duration       int               `json:"duration"`
	Name           string            `json:"name,omitempty"`
	Tags           map[string]string `json:"tags,omitempty"`
	Shared         bool              `json:"shared"`
	LocalEndpoint  Service           `json:"localEndpoint,omitempty"`
	RemoteEndpoint Service           `json:"remoteEndpoint,omitempty"`
	Annotations    []Annotation      `json:"annotations,omitempty"`

	startTime time.Time
	mutex     sync.Mutex
}

type Service struct {
	ServiceName string `json:"serviceName,omitempty"`
	IPv4        net.IP `json:"ipv4,omitempty"`
	IPv6        net.IP `json:"ipv6,omitempty"`
	Port        int    `json:"port,omitempty"`
}

type Annotation struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}

func (s *Span) Submit(c *Client) error {
	return c.Submit(s)
}

func (s *Span) AddAnnotation(timestamp time.Time, value string) {
	s.Annotations = append(s.Annotations, Annotation{
		Timestamp: timestamp.UnixMicro(),
		Value:     value,
	})
}

func (s *Span) AddJsonAnnotation(timestamp time.Time, value interface{}) {
	data, _ := json.Marshal(value)
	s.AddAnnotation(timestamp, string(data))
}

func (s *Span) AddXMLAnnotation(timestamp time.Time, value interface{}) {
	data, _ := xml.Marshal(value)
	s.AddAnnotation(timestamp, string(data))
}

// Tag - assign a tag to the span
func (s *Span) Tag(key, value string) {
	s.Lock()
	s.Tags[key] = value
	s.Unlock()
}

// NewSpan - create a new span; assign default values; generate random IDs
func NewSpan(options ...SpanOption) *Span {
	// create basic span
	span := &Span{
		SpanId:  randomID(8),
		TraceId: randomID(16),
		mutex:   sync.Mutex{},
		Tags:    make(map[string]string),
	}
	// prepend parenting operation
	options = append([]SpanOption{OptionFromParent(span.SpanId)}, options...)

	// apply span options
	span.Use(options...)
	return span
}

// NewChildSpan - Create a child Span of the Span s. Rewrite the TraceId and ParentSpanId
func (s *Span) NewChildSpan(options ...SpanOption) *Span {
	sub := NewSpan(options...)
	sub.Use(
		OptionTraceID(s.TraceId),
		OptionFromParent(s.SpanId),
	)
	return sub
}

// NewCopiedChildSpan - Create a child Span of the Span s. Copy all other parameters of Span s
func (s *Span) NewCopiedChildSpan(options ...SpanOption) *Span {
	// create a copy of Span s
	sub := *s

	// add cloning options to child span options
	options = append(options,
		OptionTraceID(s.TraceId),
		OptionFromParent(s.SpanId),
		OptionSpanID(randomID(8)), // need to randomize the spanId to not overwrite parent span
	)

	// apply options to child span
	sub.Use(options...)
	return &sub
}

func (s *Span) Start() *Span {
	s.startTime = time.Now()
	s.Timestamp = time.Now().UnixMicro()
	return s
}

func (s *Span) Finalize() *Span {
	s.Duration = int(time.Since(s.startTime).Microseconds())
	return s
}

func (s *Span) Lock() {
	s.mutex.Lock()
}

func (s *Span) Unlock() {
	s.mutex.Unlock()
}
