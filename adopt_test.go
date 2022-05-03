package simpletrace

import (
	"context"
	"testing"
)

func TestClient_EnrichContext(t *testing.T) {
	originClient1 := NewClient("https://example.com")
	originClient2 := NewClient("https://example.com")
	ctx := context.Background()
	newCtx := originClient1.EnrichContext(ctx)

	newClient, err := ClientFromContext(newCtx)
	if err != nil {
		t.Error(err)
	}
	t.Run("compare memory addresses of clients", func(t *testing.T) {
		// positive test
		if originClient1 != newClient {
			t.Errorf("client pointer mismatch: %p and %p", originClient1, newClient)
		}
		// negative test
		if originClient1 == originClient2 {
			t.Errorf("unexpected client pointer match: %p and %p", originClient1, originClient2)
		}
	})
}

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
}
