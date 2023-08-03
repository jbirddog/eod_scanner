package main

import (
	"time"
)

// TODO: pass in via config
const riskPerTrade = 50000.0 * 0.005

func PrintReport(results []*ScanResult, currentDay time.Time, writer Writer) {
	writer.WriteHeader(currentDay)

	for _, result := range results {
		result.Sort()
		writer.WriteSectionHeader(result)

		for _, v := range result.Detected {
			p := PositionFromAnalyzedData(v, riskPerTrade)
			writer.WriteRecord(v, p, riskPerTrade)
		}

		writer.WriteSectionFooter(result)
	}

	writer.WriteFooter()
}
