package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainsInt(t *testing.T) {
	haystack := []int{1, 2, 3, 4, 5, 6}
	needles := []int{4, 10}
	wantedResults := []bool{true, false}

	for i, needle := range needles {
		result := ContainsInt(haystack, needle)
		assert.Equal(t, result, wantedResults[i])
	}
}

func TestContainsString(t *testing.T) {
	haystack := []string{"aa", "bb", "ccc"}
	needles := []string{"aa", "cc"}
	wantedResults := []bool{true, false}

	for i, needle := range needles {
		result := ContainsString(haystack, needle)
		assert.Equal(t, result, wantedResults[i])
	}
}
