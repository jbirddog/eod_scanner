package main

import (
	"math"
	"testing"
)

func TestParseMinimalEODFile(t *testing.T) {
	cases := []struct {
		rawData  []string
		expected []*EODData
	}{
		// single record
		{
			[]string{
				"Symbol,Date,Open,High,Low,Close,Volume",
				"AACG,30-May-2023,1.5,1.5745,1.48,1.4906,16900",
			},
			[]*EODData{
				&EODData{
					Symbol: "AACG",
					Date:   Day(2023, 5, 30),
					Open:   1.5,
					High:   1.5745,
					Low:    1.48,
					Close:  1.4906,
					Volume: 16900,
				},
			},
		},
		// multiple records
		{
			[]string{
				"Symbol,Date,Open,High,Low,Close,Volume",
				"AACG,30-May-2023,1.5,1.5745,1.48,1.4906,16900",
				"AAL,30-May-2023,14.44,14.75,14.42,14.62,20424600",
			},
			[]*EODData{
				&EODData{
					Symbol: "AACG",
					Date:   Day(2023, 5, 30),
					Open:   1.5,
					High:   1.5745,
					Low:    1.48,
					Close:  1.4906,
					Volume: 16900,
				},
				&EODData{
					Symbol: "AAL",
					Date:   Day(2023, 5, 30),
					Open:   14.44,
					High:   14.75,
					Low:    14.42,
					Close:  14.62,
					Volume: 20424600,
				},
			},
		},
		// tests price with no decimal place
		{
			[]string{
				"Symbol,Date,Open,High,Low,Close,Volume",
				"VIA,30-Jun-2023,7,7.1846,6.69,6.96,61000",
			},
			[]*EODData{
				&EODData{
					Symbol: "VIA",
					Date:   Day(2023, 6, 30),
					Open:   7.0,
					High:   7.1846,
					Low:    6.69,
					Close:  6.96,
					Volume: 61000,
				},
			},
		},
		// invalid input - no header
		{
			[]string{
				"VIA,30-Jun-2023,7,7.1846,6.69,6.96,61000",
			},
			nil,
		},
		// invalid input - no data
		{
			[]string{
				"Symbol,Date,Open,High,Low,Close,Volume",
			},
			nil,
		},
	}

	for i, c := range cases {
		actual, err := ParseEODFile(c.rawData)

		if actual == nil {
			if c.expected != nil {
				t.Fatalf("[%d] Got nil response, expected %v", i, c.expected)
			}

			if err == nil {
				t.Fatalf("[%d] Got nil response, expected error", i)
			}

			continue
		}

		if len(actual) != len(c.expected) {
			t.Fatalf("[%d] Expected len of %d, got %d", i, len(c.expected), len(actual))
		}

		for j, a := range actual {
			expected := c.expected[j]

			if !sameData(a, expected) {
				t.Fatalf("[%d:%d] Expected %v, got %v", i, j, expected, a)
			}
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
