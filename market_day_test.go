package main

import (
	"testing"
	"time"
)

func TestKnownMarketDays(t *testing.T) {
	days := []time.Time{
		Day(2023, 1, 3),
		Day(2023, 4, 18),
		Day(2023, 6, 22),
		Day(2023, 7, 3),
		Day(2023, 10, 31),
		Day(2023, 11, 24),
	}

	for _, day := range days {
		if !IsMarketDay(day) {
			t.Fatalf("Expected %s to be a market day.", day)
		}

	}
}

func TestKnownHalfMarketDays(t *testing.T) {
	days := []time.Time{
		Day(2023, 7, 3),
		Day(2023, 11, 24),
	}

	for _, day := range days {
		if IsFullMarketDay(day) {
			t.Fatalf("Did not expect %s to be a full market day.", day)
		}
	}
}

func TestKnownNonMarketDays(t *testing.T) {
	days := []time.Time{
		Day(2023, 1, 2),
		Day(2023, 2, 18),
		Day(2023, 4, 7),
		Day(2023, 10, 29),
		Day(2023, 11, 23),
	}

	for _, day := range days {
		if IsMarketDay(day) {
			t.Fatalf("Did not expect %s to be a market day.", day)
		}

	}
}

func TestPreviousMarketDay(t *testing.T) {
	cases := []struct {
		start    time.Time
		expected time.Time
	}{
		{Day(2023, 1, 3), Day(2022, 12, 30)},
		{Day(2023, 4, 18), Day(2023, 4, 17)},
		{Day(2023, 1, 16), Day(2023, 1, 13)},
		{Day(2023, 11, 24), Day(2023, 11, 22)},
		{Day(2023, 7, 12), Day(2023, 7, 11)},
	}

	for _, c := range cases {
		actual := PreviousMarketDay(c.start)
		if actual != c.expected {
			t.Fatalf("Expected previous market day %s, got %s", c.expected, actual)
		}
	}
}

func TestPreviousMarketDays(t *testing.T) {
	cases := []struct {
		start    time.Time
		expected []time.Time
	}{
		{Day(2023, 1, 3), []time.Time{Day(2022, 12, 30), Day(2022, 12, 29)}},
		{Day(2023, 4, 18), []time.Time{Day(2023, 4, 17)}},
		{Day(2023, 7, 1), []time.Time{
			Day(2023, 6, 30),
			Day(2023, 6, 29),
			Day(2023, 6, 28),
			Day(2023, 6, 27),
			Day(2023, 6, 26),
			Day(2023, 6, 23),
		}},
	}

	for _, c := range cases {
		actual := PreviousMarketDays(c.start, len(c.expected))
		aLen := len(actual)
		eLen := len(c.expected)

		if aLen != eLen {
			t.Fatalf("Expected %d previous market days, got %d", eLen, aLen)
		}

		for i, a := range actual {
			if a != c.expected[i] {
				t.Fatalf("Expected previous market day[%d] %s, got %s", i, c.expected, a)
			}
		}
	}
}
