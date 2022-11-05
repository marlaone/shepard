package slice_test

import (
	"github.com/marlaone/shepard/slice"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	s := slice.Init(1, 2, 3)

	assert.Equal(t, slice.Init("1", "2", "3"), slice.Collect(slice.Map[int, string](s, func(i int) string { return strconv.Itoa(i) })))
}
