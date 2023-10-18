package main

import (
	"sort"
	"time"
)

type AnalyzedData struct {
	Symbol     string
	DataPoints int
	EODData    []*EODData
	Indicators Indicators
}

func NewAnalyzedData(symbol string, days int) *AnalyzedData {
	data := &AnalyzedData{
		Symbol:  symbol,
		EODData: make([]*EODData, days),
	}
	data.Indicators.Init()
	return data
}

func (a *AnalyzedData) LastClose() float64 {
	return a.EODData[len(a.EODData)-1].Close
}

func (a *AnalyzedData) LastHigh() float64 {
	return a.EODData[len(a.EODData)-1].High
}

func (a *AnalyzedData) LastLow() float64 {
	return a.EODData[len(a.EODData)-1].Low
}

func (a *AnalyzedData) PreviousClose() float64 {
	return a.EODData[len(a.EODData)-2].Close
}

func (a *AnalyzedData) MaxOfNCloses(n int) float64 {
	maxClose := 0.0

	for i := 0; i < n; i++ {
		if close := a.EODData[len(a.EODData)-1-i].Close; close > maxClose {
			maxClose = close
		}
	}

	return maxClose
}

func (a *AnalyzedData) LastChange() float64 {
	days := len(a.EODData)

	if days < 2 {
		return 0.0
	}

	return percentage(a.LastClose(), a.PreviousClose())
}

func (a *AnalyzedData) LastDate() time.Time {
	return a.EODData[len(a.EODData)-1].Date
}

func (a *AnalyzedData) LastVolume() float64 {
	return a.EODData[len(a.EODData)-1].Volume
}

func (a *AnalyzedData) LastVolumeMultiplier() float64 {
	return a.LastVolume() / a.Indicators.AvgVolume
}

type AnalyzedDataBySymbol = map[string]*AnalyzedData

func Analyze(eodData [][]*EODData) AnalyzedDataBySymbol {
	analyzed := make(AnalyzedDataBySymbol)
	days := len(eodData)

	sort.Slice(eodData, func(i, j int) bool {
		return eodData[i][0].Date.Before(eodData[j][0].Date)
	})

	for i, dailyData := range eodData {
		for _, record := range dailyData {
			symbol := record.Symbol
			data, found := analyzed[symbol]
			if !found {
				data = NewAnalyzedData(symbol, days)
				analyzed[symbol] = data
			}
			addEODData(data, record, i)
		}
	}

	return analyzed
}

func addEODData(data *AnalyzedData, record *EODData, day int) {
	data.Indicators.Add(record, day)
	data.EODData[day] = record
	data.DataPoints += 1
}
