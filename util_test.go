package main

import (
	"math"
	"testing"
)

func sameFloat(a float64, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

func sameData(a *EODData, b *EODData) bool {
	return a.Symbol == b.Symbol &&
		a.Date == b.Date &&
		sameFloat(a.Open, b.Open) &&
		sameFloat(a.High, b.High) &&
		sameFloat(a.Low, b.Low) &&
		sameFloat(a.Close, b.Close) &&
		a.Volume == b.Volume
}

func TestU8LossyLookback(t *testing.T) {
	l := &U8LossyLookback{}
	values := []float64{124.3, 1.3, 3.3, 0.3, 55.3, 78.3, 81.3, 11.3, 44.4}

	for _, v := range values {
		l.Push(v)
	}

	for i := 0; i < 8; i++ {
		actual := l.LossyValue(i)
		expected := float64(uint8(values[len(values)-i-1]))

		if !sameFloat(actual, expected) {
			t.Fatalf("Expected %f at %d, got %f\n", expected, i, actual)
		}
	}
}
