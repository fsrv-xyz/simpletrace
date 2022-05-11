package simpletrace

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	submitWorkerInput := make(chan *Span)
	submitWorkerDone := make(chan bool)
	ctx, cancelSubmitWorker := context.WithCancel(context.Background())

	// define the tracing client
	client := NewClient("http://127.0.0.1:9411/spans")
	client.Logger = log.Default()
	go client.SubmitAsyncWorker(submitWorkerInput, ctx, submitWorkerDone)

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
		OptionOfKind(KindClient),
		OptionName("testing_sub1"),
	)
	// try to create a child span from `parentSpan`
	childSpan2 := parentSpan.NewCopiedChildSpan(
		OptionRemoteEndpoint("uffl2", "126.24.242.34"),
		OptionOfKind(KindServer),
		OptionName("testing_sub2"),
	)
	// add time to child1
	childSpan1.Start()
	time.Sleep(40 * time.Millisecond)
	submitWorkerInput <- childSpan1.Finalize()

	// add time to child1
	childSpan2.Start()
	time.Sleep(30 * time.Millisecond)
	submitWorkerInput <- childSpan2.Finalize()

	submitWorkerInput <- parentSpan.Finalize()

	time.Sleep(1 * time.Second)
	cancelSubmitWorker()

	<-submitWorkerDone
	fmt.Println("worker going down")

	fmt.Println("TraceId: " + parentSpan.TraceId)
}

func TestSpan_NewCopiedChildSpan(t *testing.T) {
	parent := NewSpan(OptionName("test"))

	child1 := parent.NewCopiedChildSpan()
	child1.Tag("test", "testvalue")

	t.Run("check if tags are not the same in child and parent span", func(t *testing.T) {
		if len(parent.Tags) == len(child1.Tags) {
			t.Error("parent and child tags are the same")
		}
	})
}

func TestSpan_NewChildSpan(t *testing.T) {
	parent := NewSpan(OptionName("test"))

	child1 := parent.NewChildSpan()
	child1.Tag("test", "testvalue")

	t.Run("check if tags are not the same in child and parent span", func(t *testing.T) {
		if len(parent.Tags) == len(child1.Tags) {
			t.Error("parent and child tags are the same")
		}
	})
}

func TestSpan_Valid(t *testing.T) {
	for _, tt := range []struct {
		Name           string
		TestData       *Span
		ExpectedResult bool
	}{
		{
			Name:           "valid span",
			TestData:       NewSpan(),
			ExpectedResult: true,
		},
		{
			Name:           "invalid span",
			TestData:       &Span{},
			ExpectedResult: false,
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.ExpectedResult != tt.TestData.Valid() {
				t.Error("testdata result is not as expected")
			}
		})
	}
}
