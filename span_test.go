package simpletrace

import (
	"fmt"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	// define the tracing client
	client := NewClient("http://127.0.0.1:9411/spans")

	// define the parent span
	parentSpan := NewSpan(
		OptionLocalEndpoint("testing", "127.0.0.1:1234"),
		OptionOfKind(KindServer),
		OptionTags(map[string]string{
			"test.a": "A",
			"test.b": "B",
		}),
	)
	parentSpan.Start()
	// simulate work load for getting a time difference
	time.Sleep(20 * time.Millisecond)

	// try to create a copied child span from `parentSpan`
	childSpan1 := parentSpan.NewCopiedChildSpan(
		OptionRemoteEndpoint("uffl1", "126.24.242.34"),
		OptionLocalEndpoint("bla1", "fe80::1"),
		OptionOfKind(KindClient),
		OptionName("testing_sub"),
	)
	// try to create a child span from `parentSpan`
	childSpan2 := parentSpan.NewChildSpan(
		OptionRemoteEndpoint("uffl2", "126.24.242.34"),
		OptionLocalEndpoint("bla2", "fe80::1"),
		OptionOfKind(KindServer),
		OptionName("testing_sub"),
	)
	// add time to child1
	childSpan1.Start()
	time.Sleep(40 * time.Millisecond)
	childSpan1.Finalize()

	// add time to child1
	childSpan2.Start()
	time.Sleep(30 * time.Millisecond)
	childSpan2.Finalize()

	parentSpan.Finalize()

	fmt.Println("TraceId: " + parentSpan.TraceId)
	client.Submit(parentSpan, childSpan1, childSpan2)
}
