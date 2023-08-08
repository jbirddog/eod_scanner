package main

import (
	"container/ring"
)

// TODO: need to add some tests for these indicators

type Indicators struct {
	AvgVolume float64
	AvgClose  float64
	EMA8      EMA
	EMA12     EMA
	EMA26     EMA
	SMA20     SMA
	RSI       RSI
	MACD      MACD
}

func (i *Indicators) Init() {
	i.EMA8.Init(8)
	i.EMA12.Init(12)
	i.EMA26.Init(26)
	i.SMA20.Init(20)
	i.RSI.Periods = 14
	i.MACD.Init(&i.EMA12, &i.EMA26, 9)
}

func (i *Indicators) Add(new *EODData, previous []*EODData, period int, totalPeriods int) {
	i.AvgVolume = runningAvg(i.AvgVolume, period, new.Volume)

	close := new.Close
	i.AvgClose = runningAvg(i.AvgClose, period, close)

	i.EMA8.Add(close, period)
	i.EMA12.Add(close, period)
	i.EMA26.Add(close, period)
	i.SMA20.Add(close)
	i.RSI.Add(new, previous, period, totalPeriods)
	i.MACD.UpdateSignal(period)
}

//
// SMA
//

type SMA struct {
	Periods    int
	Value      float64
	cumulative float64
	ring       *ring.Ring
}

func (s *SMA) Init(periods int) {
	s.Periods = periods
	s.ring = ring.New(periods)
}

func (s *SMA) Add(new float64) {
	s.cumulative += new

	if s.ring.Value != nil {
		s.cumulative -= s.ring.Value.(float64)
	}

	s.ring.Value = new
	s.ring = s.ring.Next()

	s.Value = s.cumulative / float64(s.Periods)
}

//
// EMA
//

type EMA struct {
	Periods int
	Value   float64
	weight  float64
	sma     SMA
}

func (e *EMA) Init(periods int) {
	e.Periods = periods
	e.weight = 2.0 / (1.0 + float64(periods))
	e.sma.Init(periods)
}

func (e *EMA) Add(new float64, period int) {
	if period < e.Periods {
		e.sma.Add(new)
		e.Value = e.sma.Value
		return
	}

	e.Value = ((new - e.Value) * e.weight) + e.Value
}

//
// MACD
//

type MACD struct {
	Line   float64
	Signal EMA
	fast   *EMA
	slow   *EMA
}

func (m *MACD) Init(fast *EMA, slow *EMA, signalPeriods int) {
	m.fast = fast
	m.slow = slow
	m.Signal.Init(signalPeriods)
}

func (m *MACD) UpdateSignal(period int) {
	m.Line = m.fast.Value - m.slow.Value

	m.Signal.Add(m.Line, period)
}

func (m *MACD) Gap() float64 {
	return m.Line - m.Signal.Value
}

//
// RSI
//

type RSI struct {
	Periods  int
	AvgGain  float64
	AvgLoss  float64
	Value    float64
	Lookback U8LossyLookback
}

func (r *RSI) Add(new *EODData, previous []*EODData, period int, totalPeriods int) {
	if period == 0 {
		return
	}

	prevClose := 0.0

	if lookBack := previous[period-1]; lookBack != nil {
		prevClose = lookBack.Close
	}

	gain := 0.0
	loss := 0.0

	if prevClose < new.Close {
		gain = new.Close - prevClose
	} else {
		loss = prevClose - new.Close
	}

	if period <= r.Periods+1 {
		r.AvgGain = runningAvg(r.AvgGain, period-1, gain)
		r.AvgLoss = runningAvg(r.AvgLoss, period-1, loss)
		return
	}

	r.AvgGain = r.smooth(r.AvgGain, gain)
	r.AvgLoss = r.smooth(r.AvgLoss, loss)
	r.Value = 100.0 - (100.0 / (1.0 + (r.AvgGain / r.AvgLoss)))

	r.Lookback.Push(r.Value)
}

func (r *RSI) smooth(current float64, new float64) float64 {
	periods := float64(r.Periods)
	a := 1.0 / periods
	b := (periods - 1.0) / periods

	return a*new + b*current
}
