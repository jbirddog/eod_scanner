package main

import (
	"sort"
)

// TODO: Move the macd bools to flags
type AnalyzedData struct {
	Symbol     string
	DataPoints int
	AvgVolume  int
	AvgClose   float64
	MACDWasNeg bool
	MACDIsPos  bool
	MACD       MACD
	SMA20      SMA
	EODData    []*EODData

	EMA12_DEPRECATED          float64
	EMA26_DEPRECATED          float64
	SMA20_DEPRECATED          float64
	MACDSignalLine_DEPRECATED float64
	MACDLine_DEPRECATED       float64
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
			data.MACD.Init()
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

	// these are the new values that appear much more correct. at times they are pennies off
	// once indicators are factored out into their own file and tested I'm sure they will
	// tighten again
	data.MACD.Add(record, data.EODData, day)
	data.SMA20.Add(record, data.EODData, day)

	// OLD: these values are imperfect but close enough for what we are trying to do.
	// we arn't bot trading here, just trying to trim down from 10Ks of symbols to many
	// dozen of symbols to manually look at charts
	// - truth they are a little further off than I thought, being corrected above.

	daysRemaining := days - day

	if daysRemaining < 26 {
		data.EMA26_DEPRECATED = ema_deprecated(26, data.EMA26_DEPRECATED, record.Close)
	} else {
		data.EMA26_DEPRECATED = data.AvgClose
	}

	if daysRemaining < 20 {
		data.SMA20_DEPRECATED = runningAvg(data.SMA20_DEPRECATED, dpF, record.Close)
	} else {
		data.SMA20_DEPRECATED = record.Close
	}

	if daysRemaining < 12 {
		data.EMA12_DEPRECATED = ema_deprecated(12, data.EMA12_DEPRECATED, record.Close)
		data.MACDLine_DEPRECATED = data.EMA12_DEPRECATED - data.EMA26_DEPRECATED
	} else {
		data.EMA12_DEPRECATED = data.SMA20_DEPRECATED
	}

	if daysRemaining < 9 {
		data.MACDSignalLine_DEPRECATED = ema_deprecated(9, data.MACDSignalLine_DEPRECATED, data.MACDLine_DEPRECATED)
		isNeg := data.MACDLine_DEPRECATED < 0 || data.MACDSignalLine_DEPRECATED < 0
		if isNeg {
			data.MACDWasNeg = true
		}
		data.MACDIsPos = !isNeg
	} else {
		data.MACDSignalLine_DEPRECATED = data.MACDLine_DEPRECATED
	}
}

func runningAvg[T int | float64](current T, n T, new T) T {
	return (current*n + new) / (n + 1)
}

func ema_deprecated(days int, current float64, new float64) float64 {
	weight := 2.0 / (1.0 + float64(days))
	return (new * weight) + (current * (1.0 - weight))
}
