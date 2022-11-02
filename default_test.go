package shepard_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefault(t *testing.T) {
	var defaulter shepard.Default[testutils.TestType]
	defaulter = testutils.TestType{Val: ""}

	assert.Equal(t, "test", defaulter.Default().Val)
}
