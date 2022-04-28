package simpletrace

import (
	"context"
	"testing"
)

func TestSpan_EnrichContext(t *testing.T) {
	originSpan := NewSpan(OptionName("test"))
	ctx := context.Background()
	newCtx := originSpan.EnrichContext(ctx, &Client{})

	newspan, err := SpanFromContextValues(newCtx)
	if err != nil {
		t.Error(err)
	}

	t.Run("matching TraceId", func(t *testing.T) {
		if newspan.TraceId != originSpan.TraceId {
			t.Errorf(
				"traceId not matching %+v != %+v",
				newspan.TraceId,
				originSpan.TraceId,
			)
		}
	})
	t.Run("matching SpanID", func(t *testing.T) {
		if newspan.ParentSpanId != originSpan.SpanId {
			t.Errorf(
				"spanId not matching %+v != %+v",
				newspan.ParentSpanId,
				originSpan.SpanId,
			)
		}
	})
}
