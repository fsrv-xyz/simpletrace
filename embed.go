package simpletrace

import "context"

type EmbeddedTracing struct {
	Client *Client
	Span   *Span
}

// ReadCtxValues - read context values for Client and Span
func (e *EmbeddedTracing) ReadCtxValues(ctx context.Context) bool {
	var clientParseError, spanParseError error
	e.Client, clientParseError = ClientFromContext(ctx)
	e.Span, spanParseError = SpanFromContext(ctx)
	if clientParseError != nil || spanParseError != nil {
		return false
	}
	return true
}

// CopyFromCtxValues - read context values and create a child span with NewCopiedChildSpan; pass options to child span
func (e *EmbeddedTracing) CopyFromCtxValues(ctx context.Context, options ...SpanOption) bool {
	if e.ReadCtxValues(ctx) {
		e.Span = e.Span.NewCopiedChildSpan(options...)
		return true
	}
	return false
}
