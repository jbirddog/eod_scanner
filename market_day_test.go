package main

import (
	"testing"
	"time"
)

func TestKnownMarketDays(t *testing.T) {
	days := []time.Time{
		Day(2023, 01, 03),
		Day(2023, 04, 18),
		Day(2023, 06, 22),
		Day(2023, 10, 31),
		Day(2023, 11, 24),
	}

	for _, day := range days {
		if !IsMarketDay(day) {
			t.Fatalf("Expected %s to be a market day.", day)
		}

	}
}

func TestKnownNonMarketDays(t *testing.T) {
	days := []time.Time{
		Day(2023, 01, 02),
		Day(2023, 02, 18),
		Day(2023, 04, 07),
		Day(2023, 10, 29),
		Day(2023, 11, 23),
	}

	for _, day := range days {
		if IsMarketDay(day) {
			t.Fatalf("Did not expected %s to be a market day.", day)
		}

	}
}
