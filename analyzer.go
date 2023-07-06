package main

import (
	"sort"
)

type AnalyzedData struct {
	DataPoints int
	AvgVolume  int
	AvgClose   float64
	MACDSignalLine       float64
	MACDLine float64
	EMA12      float64
	EMA26      float64
	SMA20      float64
	LastClose  float64
	EODData    []*EODData
}

type AnalyzedDataBySymbol = map[string]*AnalyzedData

func Analyze(eodData [][]*EODData) AnalyzedDataBySymbol {
	analyzed := make(AnalyzedDataBySymbol)
	days := len(eodData)

	analyzedDataForSymbol := func(symbol string) *AnalyzedData {
		data, found := analyzed[symbol]
		if !found {
			data = &AnalyzedData{
				EODData: make([]*EODData, days),
			}
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
	data.LastClose = record.Close
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
		data.SMA20 = runningAvg(data.SMA20, dpF, record.Close)
	} else {
		data.SMA20 = record.Close
	}

	if daysRemaining < 12 {
		data.EMA12 = ema(12, data.EMA12, record.Close)
		data.MACDLine = data.EMA12 - data.EMA26
	} else {
		data.EMA12 = data.SMA20
	}

	if daysRemaining < 9 {
		data.MACDSignalLine = ema(9, data.MACDSignalLine, data.MACDLine)
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
