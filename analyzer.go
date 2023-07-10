package main

import (
	"sort"
)

// TODO: Move SMA/EMA to indicators.go, add tests

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

type EMA struct {
	Periods int
	Value   float64
	_sma    SMA
}

func (e *EMA) Init(periods int) {
	e.Periods = periods
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

	weight := 2.0 / (1.0 + float64(e.Periods))
	e.Value = (new.Close * weight) + (e.Value * (1.0 - weight))
}

// TODO: Move the macd bools to flags
type AnalyzedData struct {
	Symbol         string
	DataPoints     int
	AvgVolume      int
	AvgClose       float64
	MACDSignalLine float64
	MACDLine       float64
	MACDWasNeg     bool
	MACDIsPos      bool
	EMA12          EMA
	EMA26          float64
	SMA20          SMA
	EODData        []*EODData

	EMA12_DEPRECATED float64
	SMA20_DEPRECATED float64
}

type AnalyzedDataBySymbol = map[string]*AnalyzedData

func Analyze(eodData [][]*EODData) AnalyzedDataBySymbol {
	analyzed := make(AnalyzedDataBySymbol)
	days := len(eodData)

	// TODO: move init in !found to New function or similar, inline the closure
	analyzedDataForSymbol := func(symbol string) *AnalyzedData {
		data, found := analyzed[symbol]
		if !found {
			data = &AnalyzedData{
				Symbol:  symbol,
				EODData: make([]*EODData, days),
			}
			data.EMA12.Init(12)
			data.SMA20.Periods = 20
			analyzed[symbol] = data
		}
		return data
	}

	sort.Slice(eodData, func(i, j int) bool {
		return eodData[i][0].Date.Before(eodData[j][0].Date)
	})

	for i, dailyData := range eodData {
		for _, record := range dailyData {
			data := analyzedDataForSymbol(record.Symbol)
			addEODData(data, record, i, days)
		}
	}

	return analyzed
}

func addEODData(data *AnalyzedData, record *EODData, day int, days int) {
	performConstantTimeCalculations(data, record, day, days)

	data.EODData[day] = record
	data.DataPoints += 1
}

func performConstantTimeCalculations(data *AnalyzedData, record *EODData, day int, days int) {
	dp := data.DataPoints
	dpF := float64(dp)
	data.AvgVolume = runningAvg(data.AvgVolume, dp, record.Volume)
	data.AvgClose = runningAvg(data.AvgClose, dpF, record.Close)
	daysRemaining := days - day

	// these values are imperfect but close enough for what we are trying to do.
	// we arn't bot trading here, just trying to trim down from 10Ks of symbols to many
	// dozen of symbols to manually look at charts

	if daysRemaining < 26 {
		data.EMA26 = ema(26, data.EMA26, record.Close)
	} else {
		data.EMA26 = data.AvgClose
	}

	if daysRemaining < 20 {
		data.SMA20_DEPRECATED = runningAvg(data.SMA20_DEPRECATED, dpF, record.Close)
	} else {
		data.SMA20_DEPRECATED = record.Close
	}

	data.EMA12.Add(record, data.EODData, day)
	data.SMA20.Add(record, data.EODData, day)

	if daysRemaining < 12 {
		data.EMA12_DEPRECATED = ema(12, data.EMA12_DEPRECATED, record.Close)
		data.MACDLine = data.EMA12_DEPRECATED - data.EMA26
	} else {
		data.EMA12_DEPRECATED = data.SMA20_DEPRECATED
	}

	if daysRemaining < 9 {
		data.MACDSignalLine = ema(9, data.MACDSignalLine, data.MACDLine)
		isNeg := data.MACDLine < 0 || data.MACDSignalLine < 0
		if isNeg {
			data.MACDWasNeg = true
		}
		data.MACDIsPos = !isNeg
	} else {
		data.MACDSignalLine = data.MACDLine
	}
}

func runningAvg[T int | float64](current T, n T, new T) T {
	return (current*n + new) / (n + 1)
}

func ema(days int, current float64, new float64) float64 {
	weight := 2.0 / (1.0 + float64(days))
	return (new * weight) + (current * (1.0 - weight))
}
