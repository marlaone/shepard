package shepard_json_test

import (
	"encoding/json"
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/shepard_json"
	"github.com/marlaone/shepard/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestResultResponse struct {
	Data *shepard_json.Result[testutils.TestType, string] `json:"data"`
}

func TestResult_MarshalJSON(t *testing.T) {
	some := shepard_json.ParseResult(shepard.Ok[int, int](5))
	b, err := json.Marshal(some)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "5", string(b))

	someTest := shepard_json.ParseResult(shepard.Ok[testutils.TestType, string](testutils.TestType{}.Default()))
	b, err = json.Marshal(someTest)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "{\"Val\":\"test\"}", string(b))
}

func TestResult_UnmarshalJSON(t *testing.T) {
	exampleData := []byte("{\"data\":{\"Val\":\"test\"}}")
	var resp TestResultResponse
	if err := json.Unmarshal(exampleData, &resp); err != nil {
		t.Fatal(err)
	}
	assert.True(t, resp.Data.IntoResult().Equal(shepard.Ok[testutils.TestType, string](testutils.TestType{}.Default())))
}
