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
	DP      int
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
	//e.Value = (e.Weight * (new - e.Value)) + e.Value
	//e.Value = (new - e.Value) * e.Weight + e.Value
	e.DP += 1
}

//
// MACD
//

const (
	MACDCross_Neg = 1 << iota
	MACDCross_Pos
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

	m.Signal.AddPoint(m.Line)

	m.Trend += m.Line - prev
}

func (m *MACD) Gap() float64 {
	return m.Line - m.Signal.Value
}
