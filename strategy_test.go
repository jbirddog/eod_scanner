package main

import (
	"testing"
	"time"
)

// TODO: move these to own file once a few are built out
func sc1(sym string,
	d1 time.Time, o1, h1, l1, c1, v1 float64,
	d2 time.Time, o2, h2, l2, c2, v2 float64,
	ml, ms float64,
	r1, r2, r3, r4, r5 float64,
	s float64,
) *AnalyzedData {
	d := &AnalyzedData{Symbol: sym}

	d.EODData = []*EODData{
		&EODData{
			Symbol: sym,
			Date:   d1,
			Open:   o1,
			High:   h1,
			Low:    l1,
			Close:  c1,
			Volume: v1,
		},
		&EODData{
			Symbol: sym,
			Date:   d2,
			Open:   o2,
			High:   h2,
			Low:    l2,
			Close:  c2,
			Volume: v2,
		},
	}

	i := &d.Indicators
	i.Init()

	i.AvgVolume = 1000001.0
	i.AvgClose = 5.01

	i.MACD.Line = ml
	i.MACD.Signal.Value = ms

	setRSIs(i, r1, r2, r3, r4, r5)

	i.SMA20.Value = s

	return d
}

func CRDO_05152023() *AnalyzedData {
	return sc1("CRDO",
		Day(2023, 5, 12), 7.95, 8.06, 7.82, 7.91, 537400.0,
		Day(2023, 5, 15), 7.99, 8.64, 7.96, 8.61, 1463100.0,
		-0.2746, -0.3848,
		35.48, 41.89, 43.11, 42.62, 55.19,
		8.09)
}

func RIOT_01052023() *AnalyzedData {
	return sc1("RIOT",
		Day(2023, 1, 4), 3.44, 3.95, 3.38, 3.88, 12325300.0,
		Day(2023, 1, 5), 3.86, 4.28, 3.70, 4.22, 14097000.0,
		-0.2691, -0.3502,
		31.64, 31.17, 30.84, 46.35, 53.79,
		3.85)
}

func RIOT_06272023() *AnalyzedData {
	return sc1("RIOT",
		Day(2023, 6, 26), 11.49, 12.18, 10.72, 10.77, 22037414.0,
		Day(2023, 6, 27), 11.07, 11.71, 10.88, 11.65, 25685628.0,
		0.0572, -0.0364,
		56.95, 52.84, 55.66, 48.47, 55.09,
		11.01)
}

func GRPN_06282023() *AnalyzedData {
	return sc1("GRPN",
		Day(2023, 6, 27), 5.21, 5.73, 5.10, 5.56, 1205454.0,
		Day(2023, 6, 28), 5.60, 6.00, 5.49, 5.93, 1152836.0,
		0.2382, 0.2085,
		57.21, 57.61, 54.05, 57.63, 61.28,
		5.31)
}

func RIVN_06292023() *AnalyzedData {
	return sc1("RIVN",
		Day(2023, 6, 28), 13.90, 14.87, 13.82, 14.64, 32296426.0,
		Day(2023, 6, 29), 14.74, 16.01, 14.61, 16.01, 48833726.0,
		0.1716, 0.1408,
		43.35, 42.78, 47.36, 53.12, 61.91,
		14.49)
}

func TestMonthClimb(t *testing.T) {
	strategy := &MonthClimb{}
	cases := []*AnalyzedData{
		CRDO_05152023(),
		RIOT_01052023(),
		RIOT_06272023(),
		GRPN_06282023(),
		RIVN_06292023(),
	}

	for i, data := range cases {
		signaled := strategy.SignalDetected(data)

		if !signaled {
			t.Fatalf("Expected signal in case %d for %s on %s",
				i,
				data.Symbol,
				data.LastDate().Format("01/02/2006"))
		}
	}
}
