package main

import (
	"fmt"
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

var strategies = map[string]Strategy{
	"fallLevelFall": &FallLevelFall{},
	"monthClimb":    &MonthClimb{},
	"monthFall":     &MonthFall{},
}

func StrategyNamed(name string) (Strategy, error) {
	if strategy, found := strategies[name]; found {
		return strategy, nil
	}

	return nil, fmt.Errorf("Unknown strategy name: '%s'", name)
}

func hasLowVolumeOrPrice(a *AnalyzedData) bool {
	i := a.Indicators

	return i.AvgVolume < 1_000_000 || i.AvgClose < 5.0
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

	if hasLowVolumeOrPrice(a) || a.LastChange() < 0.0 {
		return false
	}

	if i.Flags&0x1 == 0 || i.Flags&0x6 == 0x6 {
		return false
	}

	if i.RSI.Value < 50 || i.RSI.LastChange() < 0.0 {
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

	if hasLowVolumeOrPrice(a) || a.LastChange() > 0.0 {
		return false
	}

	if i.Flags&0x6 == 0 {
		return false
	}

	if i.RSI.Value > 50 || i.RSI.LastChange() > 0.0 {
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
	return "MACD Fuse (WIP)"
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

//
// Fall Level Fall
//

type FallLevelFall struct{}

func (s *FallLevelFall) Name() string {
	return "Fall Level Fall"
}

func (s *FallLevelFall) SignalDetected(a *AnalyzedData) bool {
	i := a.Indicators

	if hasLowVolumeOrPrice(a) || math.Abs(a.LastChange()) > 0.25 {
		return false
	}

	if i.RSI.Value > 50 || i.MACD.Gap() >= 0.0 {
		return false
	}

	if s.lastFall(a) > -3.0 {
		return false
	}

	return true
}

func (s *FallLevelFall) SignalType() SignalType {
	return Sell
}

func (s *FallLevelFall) SortWeight(a *AnalyzedData) float64 {
	return -s.lastFall(a)
}

func (s *FallLevelFall) lastFall(a *AnalyzedData) float64 {
	minClose := min(a.LastClose(), a.PreviousClose())
	maxClose := a.MaxOfNCloses(5)
	change := percentage(minClose, maxClose)

	return change
}
