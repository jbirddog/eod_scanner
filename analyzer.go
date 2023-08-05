package main

import (
	"sort"
)

type AnalyzedData struct {
	Symbol     string
	DataPoints int
	EODData    []*EODData
	Indicators Indicators
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

	return ((close1 - close2) / close1) * 100.0
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

	// TODO: move init in !found to New function or similar, inline the closure
	analyzedDataForSymbol := func(symbol string) *AnalyzedData {
		data, found := analyzed[symbol]
		if !found {
			data = &AnalyzedData{
				Symbol:  symbol,
				EODData: make([]*EODData, days),
			}
			data.Indicators.Init()
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
	data.Indicators.Add(record, data.EODData, day, days)
	data.EODData[day] = record
	data.DataPoints += 1
}
