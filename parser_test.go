package main

import (
	"math"
	"testing"
)

func TestParseMinimalEODFile(t *testing.T) {
	rawData := []string{
		"Symbol,Date,Open,High,Low,Close,Volume",
		"AACG,30-May-2023,1.5,1.5745,1.48,1.4906,16900",
	}

	expected := []*EODData{
		&EODData{
			Symbol: "AACG",
			Date:   Day(2023, 5, 30),
			Open:   1.5,
			High:   1.5745,
			Low:    1.48,
			Close:  1.4906,
			Volume: 16900,
		},
	}

	actual, err := ParseEODFile(rawData)

	if err != nil {
		t.Fatalf("Got error: %v", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf("Expected len of %d, got %d", len(expected), len(actual))
	}

	for i, a := range actual {
		if !sameData(a, expected[i]) {
			t.Fatalf("Expected %v, got %v", expected[i], a)
		}
	}
}

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
