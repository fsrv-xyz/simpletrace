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
