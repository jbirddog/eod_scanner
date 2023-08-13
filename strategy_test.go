package main

import (
	"testing"
)

// TODO: move these to own file once a few are built out
func CRDO_05152023() *AnalyzedData {
	d := &AnalyzedData{Symbol: "CRDO"}

	d.EODData = []*EODData{
		&EODData{
			Symbol: "CRDO",
			Date:   Day(2023, 5, 12),
			Open:   7.95,
			High:   8.06,
			Low:    7.82,
			Close:  7.91,
			Volume: 537400.0,
		},
		&EODData{
			Symbol: "CRDO",
			Date:   Day(2023, 5, 15),
			Open:   7.99,
			High:   8.64,
			Low:    7.96,
			Close:  8.61,
			Volume: 1463100.0,
		},
	}

	i := &d.Indicators
	i.Init()

	i.AvgVolume = 1000001.0
	i.AvgClose = 5.01

	i.MACD.Line = -0.2746
	i.MACD.Signal.Value = -0.3848

	setRSIs(i, 35.48, 41.89, 43.11, 42.62, 55.19)

	i.SMA20.Value = 8.09

	return d
}

func RIOT_01052023() *AnalyzedData {
	d := &AnalyzedData{Symbol: "RIOT"}

	d.EODData = []*EODData{
		&EODData{
			Symbol: "RIOT",
			Date:   Day(2023, 1, 4),
			Open:   3.44,
			High:   3.95,
			Low:    3.38,
			Close:  3.88,
			Volume: 12325300.0,
		},
		&EODData{
			Symbol: "RIOT",
			Date:   Day(2023, 1, 5),
			Open:   3.86,
			High:   4.28,
			Low:    3.70,
			Close:  4.22,
			Volume: 14097000.0,
		},
	}

	i := &d.Indicators
	i.Init()

	i.AvgVolume = 1000001.0
	i.AvgClose = 5.01

	i.MACD.Line = -0.2691
	i.MACD.Signal.Value = -0.3502

	setRSIs(i, 31.64, 31.17, 30.84, 46.35, 53.79)

	i.SMA20.Value = 3.85

	return d
}

func RIOT_06272023() *AnalyzedData {
	d := &AnalyzedData{Symbol: "RIOT"}

	d.EODData = []*EODData{
		&EODData{
			Symbol: "RIOT",
			Date:   Day(2023, 6, 26),
			Open:   11.49,
			High:   12.18,
			Low:    10.72,
			Close:  10.77,
			Volume: 22037414.0,
		},
		&EODData{
			Symbol: "RIOT",
			Date:   Day(2023, 6, 27),
			Open:   11.07,
			High:   11.71,
			Low:    10.88,
			Close:  11.65,
			Volume: 25685628.0,
		},
	}

	i := &d.Indicators
	i.Init()

	i.AvgVolume = 1000001.0
	i.AvgClose = 5.01

	i.MACD.Line = 0.0572
	i.MACD.Signal.Value = -0.0364

	setRSIs(i, 56.95, 52.84, 55.66, 48.47, 55.09)

	i.SMA20.Value = 11.01

	return d
}

func GRPN_06282023() *AnalyzedData {
	d := &AnalyzedData{Symbol: "GRPN"}

	d.EODData = []*EODData{
		&EODData{
			Symbol: "GRPN",
			Date:   Day(2023, 6, 27),
			Open:   5.21,
			High:   5.73,
			Low:    5.10,
			Close:  5.56,
			Volume: 1205454.0,
		},
		&EODData{
			Symbol: "GRPN",
			Date:   Day(2023, 6, 28),
			Open:   5.60,
			High:   6.00,
			Low:    5.49,
			Close:  5.93,
			Volume: 1152836.0,
		},
	}

	i := &d.Indicators
	i.Init()

	i.AvgVolume = 1000001.0
	i.AvgClose = 5.01

	i.MACD.Line = 0.2382
	i.MACD.Signal.Value = 0.2085

	setRSIs(i, 57.21, 57.61, 54.05, 57.63, 61.28)

	i.SMA20.Value = 5.31

	return d
}

func TestMonthClimb(t *testing.T) {
	strategy := &MonthClimb{}
	cases := []*AnalyzedData{
		CRDO_05152023(),
		RIOT_01052023(),
		RIOT_06272023(),
		GRPN_06282023(),
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
