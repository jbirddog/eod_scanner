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
	"fallLevelFall":  &FallLevelFall{},
	"monthClimb":     &MonthClimb{},
	"monthFall":      &MonthFall{},
	"threeUps":       &ThreeUps{},
	"threeDowns":     &ThreeDowns{},
	"failedBBTop":    &FailedBBTop{},
	"failedBBBottom": &FailedBBBottom{},
}

func StrategyNamed(name string) (Strategy, error) {
	if strategy, found := strategies[name]; found {
		return strategy, nil
	}

	return nil, fmt.Errorf("Unknown strategy name: '%s'", name)
}

func hasLowVolumeOrPrice(a *AnalyzedData) bool {
	i := a.Indicators

	return i.AvgVolume < 1_000_000 || i.AvgClose < 5.0 || a.LastClose() < 5.0
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

	if i.MACD.Gap() >= 0.0 {
		return false
	}

	if s.lastFall(a) < 3.0 {
		return false
	}

	return true
}

func (s *FallLevelFall) SignalType() SignalType {
	return Sell
}

func (s *FallLevelFall) SortWeight(a *AnalyzedData) float64 {
	return s.lastFall(a)
}

func (s *FallLevelFall) lastFall(a *AnalyzedData) float64 {
	minClose := min(a.LastClose(), a.PreviousClose())
	maxClose := a.MaxOfNCloses(5)
	change := percentage(maxClose, minClose)

	return change
}

//
// Three Ups
//

type ThreeUps struct{}

func (s *ThreeUps) Name() string {
	return "Three Ups"
}

func (s *ThreeUps) SignalDetected(a *AnalyzedData) bool {
	if hasLowVolumeOrPrice(a) {
		return false
	}

	i := a.Indicators

	if i.RSI.Value < 50 || i.RSI.LastChange() < 0.0 {
		return false
	}

	lastThree := a.EODData[len(a.EODData)-3:]

	for i := 2; i > 0; i-- {
		cur, prev := lastThree[i], lastThree[i-1]

		if cur.Low < prev.Low || cur.Close < prev.Close {
			return false
		}
	}

	return true
}

func (s *ThreeUps) SignalType() SignalType {
	return Buy
}

func (s *ThreeUps) SortWeight(a *AnalyzedData) float64 {

	lastThree := a.EODData[len(a.EODData)-3:]

	return percentage(lastThree[2].Close, lastThree[0].Close)
}

//
// Three Downs
//

type ThreeDowns struct{}

func (s *ThreeDowns) Name() string {
	return "Three Downs"
}

func (s *ThreeDowns) SignalDetected(a *AnalyzedData) bool {
	if hasLowVolumeOrPrice(a) {
		return false
	}

	i := a.Indicators

	if i.RSI.Value > 50 || i.RSI.LastChange() > 0.0 {
		return false
	}

	lastThree := a.EODData[len(a.EODData)-3:]

	for i := 2; i > 0; i-- {
		cur, prev := lastThree[i], lastThree[i-1]

		if cur.High > prev.High || cur.Close > prev.Close {
			return false
		}
	}

	return true
}

func (s *ThreeDowns) SignalType() SignalType {
	return Sell
}

func (s *ThreeDowns) SortWeight(a *AnalyzedData) float64 {

	lastThree := a.EODData[len(a.EODData)-3:]

	return percentage(lastThree[0].Close, lastThree[2].Close)
}

//
// Failed BB Top
//

type FailedBBTop struct{}

func (s *FailedBBTop) Name() string {
	return "Failed BB Top"
}

func (s *FailedBBTop) SignalDetected(a *AnalyzedData) bool {
	if hasLowVolumeOrPrice(a) {
		return false
	}

	bb := a.Indicators.BB

	return a.PreviousClose() > bb.Upper && a.LastClose() < bb.Upper
}

func (s *FailedBBTop) SignalType() SignalType {
	return Sell
}

func (s *FailedBBTop) SortWeight(a *AnalyzedData) float64 {
	return a.PreviousClose() - a.LastClose()
}

//
// Failed BB Bottom
//

type FailedBBBottom struct{}

func (s *FailedBBBottom) Name() string {
	return "Failed BB Bottom"
}

func (s *FailedBBBottom) SignalDetected(a *AnalyzedData) bool {
	if hasLowVolumeOrPrice(a) {
		return false
	}

	bb := a.Indicators.BB

	return a.PreviousClose() < bb.Lower && a.LastClose() > bb.Lower
}

func (s *FailedBBBottom) SignalType() SignalType {
	return Buy
}

func (s *FailedBBBottom) SortWeight(a *AnalyzedData) float64 {
	return a.LastClose() - a.PreviousClose()
}
