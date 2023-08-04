package main

import (
	"testing"
)

func TestRunningAvgFloat64(t *testing.T) {
	inputs := []float64{10, 3, 20, 5, 30, 7, 40, 9, 50, 11}
	var actual float64
	expected := 18.5

	for i, input := range inputs {
		actual = runningAvg(actual, i, input)
	}

	if !sameFloat(actual, expected) {
		t.Fatalf("Expected %f, got %f\n", expected, actual)
	}
}
