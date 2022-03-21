package simpletrace

import (
	"testing"
)

func TestValidateId(t *testing.T) {
	for _, test := range []struct {
		Name           string
		Id             string
		ExpectedResult bool
	}{
		{
			Name:           "valid spanId",
			Id:             "71605d1bec5be2a1",
			ExpectedResult: true,
		},
		{
			Name:           "valid traceId",
			Id:             "0c7c4df52b1414245e181c63f8b8476a",
			ExpectedResult: true,
		},
		{
			Name:           "invalid traceId",
			Id:             "0c7c4dXf52b1411c63f8b8476a",
			ExpectedResult: false,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			if validateId(test.Id, "16,32") != test.ExpectedResult {
				t.Error("unexpected result happened")
			}
		})
	}
}
