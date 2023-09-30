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

type parsedEODData struct {
	data [][]*EODData
	err  error
}

func parse(parser Parser, dates []time.Time, c chan parsedEODData) {
	data, err := parser.Parse(dates)
	c <- parsedEODData{data: data, err: err}
}

func Scan(
	currentDay time.Time,
	marketDayCount int,
	dataDir string,
	exchange string,
	strategies []Strategy,
) ([]*ScanResult, error) {
	dates := PreviousMarketDays(currentDay, marketDayCount)
	eodData := make([][]*EODData, 0, marketDayCount)
	dateBatches := batch(dates)
	parseChan := make(chan parsedEODData, len(dateBatches))

	parser := &EODExchangeStdCSVParser{
		DataDir:  dataDir,
		Exchange: exchange,
	}

	for _, dateBatch := range dateBatches {
		go parse(parser, dateBatch, parseChan)
	}

	for i := 0; i < len(dateBatches); i++ {
		result := <-parseChan
		if result.err != nil {
			return nil, result.err
		}
		eodData = append(eodData, result.data...)
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
