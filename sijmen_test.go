package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	main()
	b.ResetTimer()
}

func TestCanReachEndVisitingAllStars(t *testing.T) {
	assert.True(t, canReachEndVisitingAllStars(0, []bool{true, false, false, false, false}, 1))
	assert.True(t, canReachEndVisitingAllStars(0, []bool{true, true, false, false, false}, 2))
	assert.True(t, canReachEndVisitingAllStars(0, []bool{true, true, true, false, false}, 3))

	assert.False(t, canReachEndVisitingAllStars(0, []bool{true, true, true, true, false}, 4))
	assert.True(t, canReachEndVisitingAllStars(0, []bool{true, false, true, true, false}, 3))

	assert.True(t, canReachEndVisitingAllStars(2, []bool{false, false, true, false, false}, 1))
	assert.True(t, canReachEndVisitingAllStars(1, []bool{false, true, false, false, false}, 1))
	assert.True(t, canReachEndVisitingAllStars(3, []bool{false, false, false, true, false}, 1))

	assert.False(t, canReachEndVisitingAllStars(7, []bool{
		true, false, false, false, false, false, false, true, false, false, false, false, true, true, false,
	}, 1))

}