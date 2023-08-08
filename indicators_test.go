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

func TestSMA(t *testing.T) {
	sma := &SMA{}
	sma.Init(5)

	sma.Add(1.1)
	sma.Add(3.3)
	sma.Add(5.5)
	sma.Add(7.7)
	sma.Add(9.9)

	if !sameFloat(sma.Value, 5.5) {
		t.Fatalf("Expected 5.5, got %f\n", sma.Value)
	}

	sma.Add(11.11)

	if !sameFloat(sma.Value, 7.502) {
		t.Fatalf("Expected 7.502, got %f\n", sma.Value)
	}
}
