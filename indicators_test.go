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

func TestBollingerBandsStdDev(t *testing.T) {
	bb := &BollingerBands{}
	bb.Init(nil, 2)

	bb.stdDev(32.0)
	bb.stdDev(47.0)
	bb.stdDev(42.0)
	bb.stdDev(45.0)
	bb.stdDev(80.0)

	actual := bb.stdDev(90.0)

	if int(actual*100) != 2326 {
		t.Fatalf("Expected 23.26, got %.2f\n", actual)
	}

	actual = bb.stdDev(52.0)

	if int(actual*100) != 2129 {
		t.Fatalf("Expected 21.29, got %.2f\n", actual)
	}
}
