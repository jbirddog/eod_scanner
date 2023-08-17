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

func (a *AnalyzedData) LastChange() float64 {
	days := len(a.EODData)

	if days < 2 {
		return 0.0
	}

	close1 := a.EODData[days-1].Close
	close2 := a.EODData[days-2].Close

	return percentage(close1, close2)
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
