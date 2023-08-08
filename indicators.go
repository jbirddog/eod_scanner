package main

import (
	"container/ring"
	"math"
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
	i.RSI.Init(14, 5)
	i.MACD.Init(&i.EMA12, &i.EMA26, 9)
}

func (i *Indicators) Add(new *EODData, period int) {
	i.AvgVolume = runningAvg(i.AvgVolume, period, new.Volume)

	close := new.Close
	i.AvgClose = runningAvg(i.AvgClose, period, close)

	i.EMA8.Add(close, period)
	i.EMA12.Add(close, period)
	i.EMA26.Add(close, period)
	i.SMA20.Add(close)
	i.RSI.Add(close, period)
	i.MACD.Update(period)
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

func (m *MACD) Update(period int) {
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
	Periods int
	Value   float64
	avgGain float64
	avgLoss float64
	prev    float64
	ring    *ring.Ring
}

func (r *RSI) Init(periods int, lookback int) {
	r.Periods = periods
	r.ring = ring.New(lookback)
}

func (r *RSI) Add(new float64, period int) {
	if period == 0 {
		r.prev = new
		return
	}

	gain := 0.0
	loss := 0.0

	if r.prev < new {
		gain = new - r.prev
	} else {
		loss = r.prev - new
	}

	r.prev = new

	if period <= r.Periods+1 {
		r.avgGain = runningAvg(r.avgGain, period-1, gain)
		r.avgLoss = runningAvg(r.avgLoss, period-1, loss)
		return
	}

	r.avgGain = r.smooth(r.avgGain, gain)
	r.avgLoss = r.smooth(r.avgLoss, loss)
	r.Value = 100.0 - (100.0 / (1.0 + (r.avgGain / r.avgLoss)))

	r.ring.Value = r.Value
	r.ring = r.ring.Next()
}

func (r *RSI) smooth(current float64, new float64) float64 {
	periods := float64(r.Periods)
	a := 1.0 / periods
	b := (periods - 1.0) / periods

	return a*new + b*current
}

func (r *RSI) LookbackMax() float64 {
	max := math.Inf(-1)

	r.ring.Do(func(val any) {
		if val != nil {
			max = math.Max(max, val.(float64))
		}
	})

	return max
}

func (r *RSI) LookbackMin() float64 {
	min := math.Inf(1)

	r.ring.Do(func(val any) {
		if val != nil {
			min = math.Min(min, val.(float64))
		}
	})

	return min
}
