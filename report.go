package main

import (
	"fmt"
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

			// TODO: string builder or buffered writer?
			fmt.Printf("%s %.2f %.0f (%.2f %.2f) (%.2f %.2f) %.2f | %d @ %.2f ~ %.2f > %.2f\n",
				v.Symbol,
				v.RSI.Value,
				v.AvgVolume,
				v.LastClose(),
				v.LastChange(),
				v.MACD.Line,
				v.MACD.Signal.Value,
				v.MACD.Trend,
				p.Shares,
				p.Entry,
				p.Capitol,
				p.StopLoss)
		}

		writer.WriteSectionFooter(result)
	}

	writer.WriteFooter()
}
