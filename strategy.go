package main

import (
	"math"
)

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
	i := a.Indicators

	return i.AvgVolume < 1000000 || i.AvgClose < 5.0
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

	// TODO: for the strategies, avgVolume makes it hard to test historic scenarios
	// since there is no way to tell the avg volume at that time
	// shift to comparison to previous volume(s)?

	if hasLowVolumeOrPrice(a) || a.LastVolume() < i.AvgVolume {
		return false
	}

	if a.LastChange() < 0.0 || a.LastClose() < i.SMA20.Value {
		return false
	}

	if i.RSI.Value < 50 {// || i.RSI.LookbackMin()+15.0 > i.RSI.Value {
		return false
	}

	if i.MACD.Gap() <= 0.0 {
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

	if hasLowVolumeOrPrice(a) || a.LastVolume() < i.AvgVolume {
		return false
	}

	if a.LastChange() > 0.0 || a.LastClose() > i.SMA20.Value {
		return false
	}

	if i.RSI.Value > 50 || i.RSI.LookbackMax()-15.0 < i.RSI.Value {
		return false
	}

	if i.MACD.Gap() >= 0.0 {
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

//
// MACD Fuse
//

type MACDFuse struct{}

func (s *MACDFuse) Name() string {
	return "MACD Fuse"
}

func (s *MACDFuse) SignalDetected(a *AnalyzedData) bool {
	if hasLowVolumeOrPrice(a) {
		return false
	}

	macd := a.Indicators.MACD
	rsi := a.Indicators.RSI

	gap := math.Abs(macd.Gap())
	gapSMA5 := math.Abs(macd.GapSMA5.Value)

	if gapSMA5 > 0.1 || gap < gapSMA5*5.0 {
		return false
	}

	if rsi.Value < 50.0 || !rsi.Rising() || percentage(rsi.Value, rsi.LookbackMax()) > 5.0 {
		return false
	}

	return true
}

func (s *MACDFuse) SignalType() SignalType {
	return Buy
}

func (s *MACDFuse) SortWeight(a *AnalyzedData) float64 {
	return -a.Indicators.MACD.Signal.Value
}
