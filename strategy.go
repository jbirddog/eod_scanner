package main

type Strategy struct {
	Name           string
	SignalDetected func(*AnalyzedData) bool
	SortWeight     func(*AnalyzedData) float64
}

//
// Month Climb
//

var MonthClimb = Strategy{
	Name:           "Month Climb",
	SignalDetected: mc_SignalDetected,
	SortWeight:     mc_SortWeight,
}

func mc_SignalDetected(a *AnalyzedData) bool {
	if a.AvgVolume < 1000000 || a.AvgClose < 5.0 ||
		a.LastVolume() < a.AvgVolume ||
		a.LastChange() < 0.0 ||
		a.MACD.Gap() < 0.0 ||
		a.RSI.Value < 50.0 ||
		a.LastClose() < a.SMA20.Value {
		return false
	}

	return true
}

func mc_SortWeight(a *AnalyzedData) float64 {
	macdWeight := a.MACD.Line * a.MACD.Gap()
	volumeWeight := a.LastVolumeMultiplier()
	weight := macdWeight * volumeWeight

	return weight
}

//
// Month Fall
//

var MonthFall = Strategy{
	Name:           "Month Fall",
	SignalDetected: mf_SignalDetected,
	SortWeight:     mf_SortWeight,
}

func mf_SignalDetected(a *AnalyzedData) bool {
	if a.AvgVolume < 1000000 || a.AvgClose < 5.0 ||
		a.LastVolume() < a.AvgVolume ||
		a.LastChange() > 0.0 ||
		a.MACD.Gap() > 0.0 ||
		a.RSI.Value > 50.0 ||
		a.LastClose() > a.SMA20.Value {
		return false
	}

	return true
}

func mf_SortWeight(a *AnalyzedData) float64 {
	macdWeight := a.MACD.Signal.Value * a.MACD.Gap()
	volumeWeight := a.LastVolumeMultiplier()
	weight := macdWeight * volumeWeight

	return weight
}
