package main

type SignalType int

const (
	Buy SignalType = iota
	Sell
)

type Strategy interface {
	Name() string
	SignalDetected(a *AnalyzedData) bool
	SignalType() SignalType
	SortWeight(a *AnalyzedData) float64
}

func hasLowVolumeOrPrice(a *AnalyzedData) bool {
	avgClose := a.Indicators.AvgClose
	avgVol := a.Indicators.AvgVolume

	return avgVol < 1000000 || a.LastVolume() < avgVol || avgClose < 5.0
}

//
// Month Climb
//

type MonthClimb struct{}

func (s *MonthClimb) Name() string {
	return "Month Climb"
}

func (s *MonthClimb) SignalDetected(a *AnalyzedData) bool {
	i := a.Indicators

	if hasLowVolumeOrPrice(a) {
		return false
	}

	if a.LastChange() < 0.0 || a.LastClose() < i.SMA20.Value {
		return false
	}

	if i.RSI.Value < 50 || i.RSI.Lookback.LossyValue(5)+15.0 > i.RSI.Value {
		return false
	}

	if i.MACD.Gap() < 0.0 {
		return false
	}

	return true
}

func (s *MonthClimb) SignalType() SignalType {
	return Buy
}

func (s *MonthClimb) SortWeight(a *AnalyzedData) float64 {
	return -a.Indicators.MACD.Line
}

//
// Month Fall
//

type MonthFall struct{}

func (s *MonthFall) Name() string {
	return "Month Fall"
}

func (s *MonthFall) SignalDetected(a *AnalyzedData) bool {
	i := a.Indicators

	if hasLowVolumeOrPrice(a) {
		return false
	}

	if a.LastChange() > 0.0 || a.LastClose() > i.SMA20.Value {
		return false
	}

	if i.RSI.Value > 50 || i.RSI.Lookback.LossyValue(5)-15.0 < i.RSI.Value {
		return false
	}

	if i.MACD.Gap() > 0.0 {
		return false
	}

	return true
}

func (s *MonthFall) SignalType() SignalType {
	return Sell
}

func (s *MonthFall) SortWeight(a *AnalyzedData) float64 {
	return a.Indicators.MACD.Line
}
