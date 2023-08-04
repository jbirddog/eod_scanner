package main

import (
	"sort"
	"time"
)

type ScanResult struct {
	Strategy Strategy
	Detected []*AnalyzedData
}

func newResults(strategies []Strategy, detectedCap int) []*ScanResult {
	results := make([]*ScanResult, len(strategies))

	for i, strategy := range strategies {
		results[i] = &ScanResult{
			Strategy: strategy,
			Detected: make([]*AnalyzedData, 0, detectedCap),
		}
	}

	return results
}

func (s *ScanResult) Sort() {
	sort.Slice(s.Detected, func(i, j int) bool {
		return s.Strategy.SortWeight(s.Detected[i]) > s.Strategy.SortWeight(s.Detected[j])
	})

}

func Scan(currentDay time.Time, marketDayCount int, dataDir string, strategies []Strategy) ([]*ScanResult, error) {
	// TODO: use channels?
	dates := PreviousMarketDays(currentDay, marketDayCount)
	// TODO: AMEX, NYSE
	exchange := "NASDAQ"
	eodData := make([][]*EODData, marketDayCount)

	for i := len(dates) - 1; i >= 0; i-- {
		date := dates[i]
		rawData, err := LoadEODFile(dataDir, exchange, date)
		if err != nil {
			return nil, err
		}

		data, err := ParseEODFile(rawData)
		if err != nil {
			return nil, err
		}

		eodData[i] = data
	}

	analyzedDataBySymbol := Analyze(eodData)
	results := newResults(strategies, len(analyzedDataBySymbol))

	for _, v := range analyzedDataBySymbol {
		if v.DataPoints != marketDayCount {
			continue
		}

		for i, strategy := range strategies {
			if strategy.SignalDetected(v) {
				results[i].Detected = append(results[i].Detected, v)
			}
		}
	}

	return results, nil
}
