package main

import (
	"testing"
)

func TestEODFilePath(t *testing.T) {
	actual := EODFilePath("/data", "NASDAQ", Day(2023, 7, 4))
	expected := "/data/NASDAQ_20230704.csv"

	if actual != expected {
		t.Fatalf("Expected eod file name %s, got %s", expected, actual)
	}
}
