package main

import (
	"testing"

	assert "gotest.tools/assert"
)

// template for running test
func TestAdd(t *testing.T) {
	sum := Add(1, 2)
	expected := 3

	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum) 
	}

	assert.Equal(t, sum, expected)
}

func runCommand(cmd string) ([]byte, error) {
	return nil, nil 
}

func TestQuizProgram(t *testing.T) {
	return 
}
