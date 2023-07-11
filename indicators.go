package main

// TODO: need to add some tests for these indicators

//
// SMA
//

type SMA struct {
	Periods    int
	Cumulative float64
}

func (s *SMA) Add(new *EODData, previous []*EODData, period int) {
	s.Cumulative += new.Close

	if period < s.Periods {
		return
	}

	if lookBack := previous[period-s.Periods]; lookBack != nil {
		s.Cumulative -= lookBack.Close

	}
}

func (s *SMA) Value() float64 {
	return s.Cumulative / float64(s.Periods)
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

func (e *EMA) Add(new *EODData, previous []*EODData, period int) {
	if period < e.Periods {
		e._sma.Add(new, previous, period)

		if period == e.Periods-1 {
			e.Value = e._sma.Value()
		}

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

type MACD struct {
	Line   float64
	Signal EMA
	_ema12 EMA
	_ema26 EMA
}

func (m *MACD) Init() {
	m._ema12.Init(12)
	m._ema26.Init(26)
	m.Signal.Init(9)
}

func (m *MACD) Add(new *EODData, previous []*EODData, period int) {
	m._ema12.Add(new, previous, period)
	m._ema26.Add(new, previous, period)
	m.Line = m._ema12.Value - m._ema26.Value

	if period < m._ema26.Periods {
		m.Signal.Value = m.Line
		return
	}

	m.Signal.AddPoint(m.Line)
}
