package simpletrace

import (
	"context"
	"testing"
)

func TestSpan_EnrichContext(t *testing.T) {
	originSpan1 := NewSpan(OptionName("test"))
	originSpan2 := NewSpan(OptionName("test"))
	ctx := context.Background()
	newCtx := originSpan1.EnrichContext(ctx)

	newspan, err := SpanFromContext(newCtx)
	if err != nil {
		t.Error(err)
	}
	t.Run("compare memory addresses of spans", func(t *testing.T) {
		// positive test
		if originSpan1 != newspan {
			t.Errorf("span pointer mismatch: %p and %p", originSpan1, newspan)
		}
		// negative test
		if originSpan1 == originSpan2 {
			t.Errorf("unexpected span pointer match: %p and %p", originSpan1, originSpan2)
		}
	})

	t.Run("matching TraceId", func(t *testing.T) {
		if newspan.TraceId != originSpan1.TraceId {
			t.Errorf(
				"traceId not matching %+v != %+v",
				newspan.TraceId,
				originSpan1.TraceId,
			)
		}
	})
	t.Run("matching SpanID", func(t *testing.T) {
		if newspan.ParentSpanId != originSpan1.SpanId {
			t.Errorf(
				"spanId not matching %+v != %+v",
				newspan.ParentSpanId,
				originSpan1.SpanId,
			)
		}
	})
}
