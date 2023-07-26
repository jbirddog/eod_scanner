package main

// TODO: need to add some tests for these indicators

//
// SMA
//

type SMA struct {
	Periods    int
	Cumulative float64
	Value      float64
}

func (s *SMA) Add(new *EODData, previous []*EODData, period int) {
	s.Cumulative += new.Close

	if period < s.Periods {
		s.Value = s.Cumulative / float64(period+1)
		return
	}

	if lookBack := previous[period-s.Periods]; lookBack != nil {
		s.Cumulative -= lookBack.Close
	}

	s.Value = s.Cumulative / float64(s.Periods)
}

//
// EMA
//

type EMA struct {
	Periods int
	Weight  float64
	Value   float64
	_sma    SMA
}

func (e *EMA) Init(periods int) {
	e.Periods = periods
	e.Weight = 2.0 / (1.0 + float64(periods))
	e._sma.Periods = periods
}

func (e *EMA) Add(new *EODData, previous []*EODData, period int, totalPeriods int) {
	daysLeft := totalPeriods - period
	if daysLeft > e.Periods {
		e._sma.Add(new, previous, period)
		e.Value = e._sma.Value
		return
	}

	e.AddPoint(new.Close)
}

func (e *EMA) AddPoint(new float64) {
	e.Value = (new * e.Weight) + (e.Value * (1.0 - e.Weight))
}

//
// MACD
//

const (
	MACDWas_Pos = 1 << iota
	MACDWas_Neg
)

type MACD struct {
	Line   float64
	Signal EMA
	Flags  int
	Trend  float64
	_ema12 EMA
	_ema26 EMA
}

func (m *MACD) Init() {
	m._ema12.Init(12)
	m._ema26.Init(26)
	m.Signal.Init(9)
}

func (m *MACD) Add(new *EODData, previous []*EODData, period int, totalPeriods int) {
	m._ema12.Add(new, previous, period, totalPeriods)
	m._ema26.Add(new, previous, period, totalPeriods)

	daysLeft := totalPeriods - period
	prev := m.Line
	m.Line = m._ema12.Value - m._ema26.Value

	if daysLeft > m.Signal.Periods {
		m.Signal.Value = m.Line
		return
	}

	m.Trend += m.Line - prev
	m.Signal.AddPoint(m.Line)

	if m.Line > m.Signal.Value {
		m.Flags |= MACDWas_Pos
	} else {
		m.Flags |= MACDWas_Neg
	}
}

func (m *MACD) Gap() float64 {
	return m.Line - m.Signal.Value
}

//
// RSI
//

type RSI struct {
	Periods int
	AvgGain float64
	AvgLoss float64
	Value   float64
}

const RSI_SMOOTHER = 13.0

func (r *RSI) Init() {
	r.Periods = 14
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

	if period == 1 {
		r.AvgGain = gain
		r.AvgLoss = loss
		return
	}

	if period <= r.Periods+1 {
		r.AvgGain = runningAvg(r.AvgGain, period-1, gain)
		r.AvgLoss = runningAvg(r.AvgLoss, period-1, loss)
		return
	}

	r.AvgGain = r.smooth(r.AvgGain, gain)
	r.AvgLoss = r.smooth(r.AvgLoss, loss)
	r.Value = 100.0 - (100.0 / (1.0 + (r.AvgGain / r.AvgLoss)))
}

func (r *RSI) smooth(current float64, new float64) float64 {
	periods := float64(r.Periods)
	a := 1.0 / periods
	b := (periods - 1) / periods

	return a*new + b*current
}

//
// utils
//

func runningAvg(current float64, n int, new float64) float64 {
	n64 := float64(n)
	return (current*n64 + new) / (n64 + 1.0)
}
