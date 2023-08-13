package main

import (
	"testing"
)

// TODO: move these to own file once a few are built out
func setRSIs(i *Indicators, a, b, c, d, e float64) {
	r := i.RSI.ring

	r.Value = a
	r = r.Next()

	r.Value = b
	r = r.Next()

	r.Value = c
	r = r.Next()

	r.Value = d
	r = r.Next()

	r.Value = e
	r = r.Next()

	i.RSI.Value = e
}

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

func Test_CRDO_05152023(t *testing.T) {
	strategy := &MonthClimb{}
	data := CRDO_05152023()
	signaled := strategy.SignalDetected(data)

	t.Logf("MACD gap: %f\n", data.Indicators.MACD.Gap())

	if !signaled {
		t.Fatalf("Expected '%s' signal for CRDO on 05/15/2023", strategy.Name())
	}
}
