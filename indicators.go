package main

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
	Weight float64
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

	e.Value = (new.Close * e.Weight) + (e.Value * (1.0 - e.Weight))
}
