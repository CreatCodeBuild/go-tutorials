package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Map(data []int, mapper func(int) int) []int {
	results := make([]int, len(data))
	for i, ele := range data {
		results[i] = mapper(ele)
	}
	return results
}

func TestMap(t *testing.T) {
	// Map
	results := Map([]int{1, 2, 3}, func(x int) int { return x + 1 })
	require.Equal(t, []int{2, 3, 4}, results)
}
