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

//
// Month Climb
//

type MonthClimb struct{}

func (s *MonthClimb) Name() string {
	return "Month Climb"
}

func (s *MonthClimb) SignalDetected(a *AnalyzedData) bool {
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

func (s *MonthClimb) SignalType() SignalType {
	return Buy
}

func (s *MonthClimb) SortWeight(a *AnalyzedData) float64 {
	macdWeight := a.MACD.Line * a.MACD.Gap()
	volumeWeight := a.LastVolumeMultiplier()
	weight := macdWeight * volumeWeight

	return weight
}

//
// Month Fall
//

type MonthFall struct{}

func (s *MonthFall) Name() string {
	return "Month Fall"
}

func (s *MonthFall) SignalDetected(a *AnalyzedData) bool {
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

func (s *MonthFall) SignalType() SignalType {
	return Sell
}

func (s *MonthFall) SortWeight(a *AnalyzedData) float64 {
	macdWeight := a.MACD.Signal.Value * a.MACD.Gap()
	volumeWeight := a.LastVolumeMultiplier()
	weight := macdWeight * volumeWeight

	return weight
}
