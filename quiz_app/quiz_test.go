package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestAdd(t *testing.T) {
	sum := Add(1, 2)
	expected := 3

	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum)
	}

	assert.Equal(t, sum, expected)
}
