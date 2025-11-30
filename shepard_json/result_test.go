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

func TestResult_UnmarshalJSON_EmptyObject(t *testing.T) {
	exampleData := []byte("{\"data\":{}}")
	var resp TestResultResponse
	if err := json.Unmarshal(exampleData, &resp); err != nil {
		t.Fatal(err)
	}
	// Empty object should result in zero value
	result := resp.Data.IntoResult()
	assert.True(t, result.IsOk())
}

func TestResult_MarshalJSON_Err(t *testing.T) {
	errResult := shepard_json.ParseResult(shepard.Err[int, string]("error message"))
	b, err := json.Marshal(errResult)
	// Err case returns nil, nil which causes a MarshalerError
	// This is expected behavior - Err results don't serialize cleanly
	assert.NotNil(t, err)
	assert.Empty(t, b)
}

func TestResult_IntoResult_Nil(t *testing.T) {
	var nilResult *shepard_json.Result[int, string]
	result := nilResult.IntoResult()
	// nil receiver should return Ok with zero value
	assert.True(t, result.IsOk())
	assert.Equal(t, 0, result.Unwrap())
}

func TestParseResult(t *testing.T) {
	// Test with Ok result
	okResult := shepard.Ok[string, error]("hello")
	parsed := shepard_json.ParseResult(okResult)
	assert.NotNil(t, parsed)
	assert.True(t, parsed.IntoResult().IsOk())
	assert.Equal(t, "hello", parsed.IntoResult().Unwrap())

	// Test with Err result
	errResult := shepard.Err[string, string]("error")
	parsedErr := shepard_json.ParseResult(errResult)
	assert.NotNil(t, parsedErr)
	assert.True(t, parsedErr.IntoResult().IsErr())
}
