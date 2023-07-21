package main

import (
	"fmt"
	"sort"
	"time"
)

func PrintReport(symbols []*AnalyzedData, currentDay time.Time) {
	sort.Slice(symbols, func(i, j int) bool {
		return symbols[i].SortWeight() > symbols[j].SortWeight()
	})

	for _, v := range symbols {
		p := PositionFromAnalyzedData(v, riskPerTrade)

		// TODO: string builder or buffered writer?
		fmt.Printf("%s %d (%.2f) (%.2f %.2f -- %.2f %.2f) %.2f | %d @ %.2f ~ %.2f > %.2f\n",
			v.Symbol,
			v.AvgVolume,
			v.SMA20.Value,
			v.MACD._ema12.Value,
			v.MACD._ema26.Value,
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
