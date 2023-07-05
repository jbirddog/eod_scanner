package main

import (
"sort"
)

type AnalyzedData struct {
	DataPoints int
	AvgVolume  int
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
			addEODData(data, record, i)
		}
	}

	return analyzed
}

func addEODData(data *AnalyzedData, record *EODData, slot int) {
	data.AvgVolume = (data.AvgVolume*data.DataPoints + record.Volume) / (data.DataPoints + 1)

	data.EODData[slot] = record
	data.DataPoints += 1
}
