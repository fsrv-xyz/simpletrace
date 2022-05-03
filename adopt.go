package simpletrace

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Header string

const (
	HeaderTraceId           Header = "X-B3-TraceId"
	HeaderParentSpanId      Header = "X-B3-ParentSpanId"
	ContextKeySpan          Header = "simpletrace/span-context"
	ContextSubmissionClient Header = "simpletrace/submission-context"
)

// EnrichContext - add required IDs/URLs to existing context
func (c *Client) EnrichContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, ContextSubmissionClient, c)
	return ctx
}

// ClientFromContext - generate the span submission client with values from ctx
func ClientFromContext(ctx context.Context) (*Client, error) {
	var client *Client
	switch ctx.Value(ContextSubmissionClient).(type) {
	case *Client:
		client = ctx.Value(ContextSubmissionClient).(*Client)
	default:
		return nil, fmt.Errorf("value of %+q not found in context", ContextSubmissionClient)
	}
	return client, nil
}

// EnrichContext - add required IDs/URLs to existing context
func (s *Span) EnrichContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, ContextKeySpan, s)
	return ctx
}

// SpanFromContext - generate the parent Span with values from ctx
func SpanFromContext(ctx context.Context) (*Span, error) {
	var span *Span
	switch ctx.Value(ContextKeySpan).(type) {
	case *Span:
		span = ctx.Value(ContextKeySpan).(*Span)
	default:
		return nil, fmt.Errorf("value of %+q not found in context", ContextKeySpan)
	}
	return span, nil
}

// SpanFromHttpHeader - generate the parent Span with parameters from request headers
func SpanFromHttpHeader(r *http.Request) (*Span, error) {
	spanId := r.Header.Get(string(HeaderParentSpanId))
	traceId := r.Header.Get(string(HeaderTraceId))
	if !validateSpanID(spanId) || !validateTraceID(traceId) {
		return nil, errors.New("one ore multiple header values not found/malformed")
	}
	span := NewSpan(OptionShared(), OptionSpanID(spanId), OptionFromParent(spanId), OptionTraceID(traceId))
	return span, nil
}
