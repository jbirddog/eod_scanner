package main

import (
	"sort"
)

type AnalyzedData struct {
	DataPoints int
	AvgVolume  int
	SMA20      float64
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
}

func performConstantTimeCalculations(data *AnalyzedData, record *EODData, day int, days int) {
	data.AvgVolume = runningAvg(data.AvgVolume, data.DataPoints, record.Volume)

	if days-day < 20 {
		data.SMA20 = runningAvg(data.SMA20, float64(data.DataPoints), record.Close)
	} else {
		data.SMA20 = record.Close
	}
}

func runningAvg[T int | float64](current T, n T, new T) T {
	return (current*n + new) / (n + 1)
}
