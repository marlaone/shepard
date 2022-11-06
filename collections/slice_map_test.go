package collections_test

import (
	"github.com/marlaone/shepard/collections"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	s := collections.Init(1, 2, 3)

	assert.Equal(t, collections.Init("1", "2", "3"), collections.CollectSlice(collections.MapSlice[int, string](s, func(i int) string { return strconv.Itoa(i) })))
}
