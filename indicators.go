package main

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
	i.SMA20.Periods = 20
	i.RSI.Periods = 14
	i.MACD.Init(&i.EMA12, &i.EMA26, 9)
}

func (i *Indicators) Add(new *EODData, previous []*EODData, period int, totalPeriods int) {
	i.AvgVolume = runningAvg(i.AvgVolume, period, new.Volume)
	i.AvgClose = runningAvg(i.AvgClose, period, new.Close)

	i.EMA8.Add(new, previous, period, totalPeriods)
	i.EMA12.Add(new, previous, period, totalPeriods)
	i.EMA26.Add(new, previous, period, totalPeriods)
	i.SMA20.Add(new, previous, period)
	i.RSI.Add(new, previous, period, totalPeriods)
	i.MACD.Add(new, previous, period, totalPeriods)
}

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
	_fast  *EMA
	_slow  *EMA
}

func (m *MACD) Init(fast *EMA, slow *EMA, signalPeriods int) {
	m._fast = fast
	m._slow = slow
	m.Signal.Init(signalPeriods)
}

func (m *MACD) Add(new *EODData, previous []*EODData, period int, totalPeriods int) {
	daysLeft := totalPeriods - period
	m.Line = m._fast.Value - m._slow.Value

	if daysLeft > m.Signal.Periods {
		m.Signal.Value = m.Line
		return
	}

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
}

func (r *RSI) smooth(current float64, new float64) float64 {
	periods := float64(r.Periods)
	a := 1.0 / periods
	b := (periods - 1.0) / periods

	return a*new + b*current
}

//
// utils
//

func runningAvg(current float64, n int, new float64) float64 {
	n64 := float64(n)
	return (current*n64 + new) / (n64 + 1.0)
}
