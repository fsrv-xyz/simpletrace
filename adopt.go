package simpletrace

import (
	"context"
	"errors"
	"fmt"
)

type Header string

const (
	HeaderTraceId       Header = "X-B3-TraceId"
	HeaderParentSpanId  Header = "X-B3-ParentSpanId"
	HeaderTraceEndpoint Header = "X-B3-TraceEndpoint"
)

// EnrichContext - add required IDs/URLs to existing context
func (s *Span) EnrichContext(ctx context.Context, client *Client) context.Context {
	for key, value := range map[Header]string{
		HeaderTraceEndpoint: client.URL,
		HeaderTraceId:       s.TraceId,
		HeaderParentSpanId:  s.ParentSpanId,
	} {
		ctx = context.WithValue(ctx, key, value)
	}
	return ctx
}

// SpanFromContextValues - generate the parent Span with values from ctx
func SpanFromContextValues(ctx context.Context) (*Span, error) {
	spanId, parentIdFound := ctx.Value(HeaderParentSpanId).(string)
	traceId, traceIdFound := ctx.Value(HeaderTraceId).(string)
	if !parentIdFound || !traceIdFound {
		return nil, errors.New("one ore multiple Values not found")
	}
	span := NewSpan(OptionShared(), OptionSpanID(spanId), OptionFromParent(spanId), OptionTraceID(traceId))
	return span, nil
}

// ClientFromContextValues - generate a Client with url from ctx
func ClientFromContextValues(ctx context.Context) (*Client, error) {
	clientURL, found := ctx.Value(HeaderTraceEndpoint).(string)
	if !found {
		return nil, fmt.Errorf("%+q not found in context", HeaderTraceEndpoint)
	}
	client := NewClient(clientURL)
	return &client, nil
}
