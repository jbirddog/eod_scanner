package main

import (
	"fmt"
	"sort"
	"time"
)

func PrintReport(symbols []*AnalyzedData, currentDay time.Time) {
	sort.Slice(symbols, func(i, j int) bool {
		return MonthClimb.SortWeight(symbols[i]) > MonthClimb.SortWeight(symbols[j])
	})

	for _, v := range symbols {
		p := PositionFromAnalyzedData(v, riskPerTrade)

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

	fmt.Printf("Found %d symbols. %.2f risk per trade\n\n", len(symbols), riskPerTrade)
}
