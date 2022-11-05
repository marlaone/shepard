package num_test

import (
	"github.com/marlaone/shepard/iter"
	"github.com/marlaone/shepard/num"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSum(t *testing.T) {
	assert.Equal(t, 10, num.Sum(iter.New([]int{1, 2, 3, 4})))
	assert.Equal(t, 8, num.Sum(iter.New([]int{-1, 2, 3, 4})))
}
