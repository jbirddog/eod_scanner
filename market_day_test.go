package main

import (
	"testing"
	"time"
)

func TestKnownMarketDays(t *testing.T) {
	days := []time.Time{
		day(2023, 01, 03),
		day(2023, 04, 18),
		day(2023, 06, 22),
		day(2023, 10, 31),
		day(2023, 11, 24),
	}

	for _, day := range days {
		if !IsMarketDay(day) {
			t.Fatalf("Expected %s to be a market day.", day)
		}

	}
}

func TestKnownNonMarketDays(t *testing.T) {
	days := []time.Time{
		day(2023, 01, 02),
		day(2023, 02, 18),
		day(2023, 04, 06),
		day(2023, 10, 31),
		day(2023, 11, 22),
	}

	for _, day := range days {
		if IsMarketDay(day) {
			t.Fatalf("Did not expected %s to be a market day.", day)
		}

	}
}
