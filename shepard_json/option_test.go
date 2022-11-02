package shepard_json_test

import (
	"encoding/json"
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/shepard_json"
	"github.com/marlaone/shepard/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestOptionResponse struct {
	Data *shepard_json.Option[testutils.TestType] `json:"data"`
	Size *shepard_json.Option[uint]               `json:"size"`
}

func TestOption_MarshalJSON(t *testing.T) {
	some := shepard_json.ParseOption(shepard.Some[int](5))
	b, err := json.Marshal(some)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "5", string(b))

	someTest := shepard_json.ParseOption(shepard.Some[testutils.TestType](testutils.TestType{}.Default()))
	b, err = json.Marshal(someTest)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "{\"Val\":\"test\"}", string(b))
}

func TestOption_UnmarshalJSON(t *testing.T) {
	exampleData := []byte("{\"data\":{\"Val\":\"test\"}}")
	var resp TestOptionResponse
	if err := json.Unmarshal(exampleData, &resp); err != nil {
		t.Fatal(err)
	}
	assert.True(t, resp.Data.IntoOption().Equal(shepard.Some[testutils.TestType](testutils.TestType{}.Default())))

	exampleData2 := []byte("{\"data\":{\"Val\":\"\"}}")
	if err := json.Unmarshal(exampleData2, &resp); err != nil {
		t.Fatal(err)
	}
	assert.True(t, resp.Data.IntoOption().Equal(shepard.Some[testutils.TestType](testutils.TestType{Val: ""})))

	exampleData3 := []byte("{\"data\":{}}")
	if err := json.Unmarshal(exampleData3, &resp); err != nil {
		t.Fatal(err)
	}
	assert.True(t, resp.Data.IntoOption().Equal(shepard.None[testutils.TestType]()))
	assert.True(t, resp.Size.IntoOption().Equal(shepard.None[uint]()))

	exampleData4 := []byte("{\"data\":{},\"size\":0}")
	if err := json.Unmarshal(exampleData4, &resp); err != nil {
		t.Fatal(err)
	}
	assert.True(t, resp.Data.IntoOption().Equal(shepard.None[testutils.TestType]()))
	assert.True(t, resp.Size.IntoOption().Equal(shepard.Some[uint](0)))
}
