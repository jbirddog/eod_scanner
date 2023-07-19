package main

import (
	"sort"
)

type AnalyzedData struct {
	Symbol     string
	DataPoints int
	AvgVolume  int
	AvgClose   float64
	MACD       MACD
	SMA20      SMA
	EODData    []*EODData
}

func (a *AnalyzedData) LastClose() float64 {
	return a.EODData[len(a.EODData)-1].Close
}

func (a *AnalyzedData) ClosedUp() bool {
	data := a.EODData[len(a.EODData)-1]
	return data.Close > data.Open
}

func (a *AnalyzedData) LastVolume() int {
	return a.EODData[len(a.EODData)-1].Volume
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

	data.MACD.Add(record, data.EODData, day, days)
	data.SMA20.Add(record, data.EODData, day)
}

func runningAvg[T int | float64](current T, n T, new T) T {
	return (current*n + new) / (n + 1)
}
