package simpletrace

import "context"

type Embed struct {
	client *Client
	span   *Span
}

// ReadCtxValues - read context values for Client and Span
func (e *Embed) ReadCtxValues(ctx context.Context) bool {
	var clientParseError, spanParseError error
	e.client, clientParseError = ClientFromContext(ctx)
	e.span, spanParseError = SpanFromContext(ctx)
	if clientParseError != nil || spanParseError != nil {
		return false
	}
	return true
}

// CopyFromCtxValues - read context values and create a child span with NewCopiedChildSpan; pass options to child span
func (e *Embed) CopyFromCtxValues(ctx context.Context, options ...SpanOption) bool {
	if e.ReadCtxValues(ctx) {
		e.span = e.span.NewCopiedChildSpan(options...)
		return true
	}
	return false
}
