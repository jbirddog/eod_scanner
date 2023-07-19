package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

const marketDayCount = 52
const riskPerTrade = 1000.0

func main() {
	dataDir := os.Getenv("EOD_DATA_DIR")
	if dataDir == "" {
		log.Fatal("Must set environment variable EOD_DATA_DIR")
	}

	// TODO: move to driver, use channels
	dates := PreviousMarketDays(Day(2023, 7, 19), marketDayCount)
	// TODO: AMEX, NYSE
	exchange := "NASDAQ"
	eodData := make([][]*EODData, marketDayCount)

	for i := len(dates) - 1; i >= 0; i-- {
		date := dates[i]
		rawData, err := LoadEODFile(dataDir, exchange, date)
		if err != nil {
			log.Fatal(err)
		}

		data, err := ParseEODFile(rawData)
		if err != nil {
			log.Fatal(err)
		}

		eodData[i] = data
	}

	analyzedDataBySymbol := Analyze(eodData)
	symbols := make([]*AnalyzedData, 0, len(analyzedDataBySymbol))

	for _, v := range analyzedDataBySymbol {
		if v.DataPoints != marketDayCount {
			continue
		}

		if v.AvgVolume < 500000 || v.AvgClose < 5.0 {
			continue
		}

		if !v.ClosedUp() {
			continue
		}

		// TODO: break out into buy vs sell signals
		if v.MACD.Gap() < 0 {
			continue
		}

		if v.LastClose() < v.SMA20.Value {
			continue
		}

		if v.LastVolume() < v.AvgVolume {
			continue
		}

		symbols = append(symbols, v)
	}

	sort.Slice(symbols, func(i, j int) bool {
		// TODO: define the real sort criteria
		symbol_a := symbols[i]
		symbol_b := symbols[j]
		a := float64(symbol_a.LastVolume()-symbol_a.AvgVolume) / symbol_a.MACD.Gap()
		b := float64(symbol_b.LastVolume()-symbol_b.AvgVolume) / symbol_b.MACD.Gap()

		return a < b
	})

	for _, v := range symbols {
		p := PositionFromAnalyzedData(v, riskPerTrade)

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

	fmt.Printf("Found %d symbols\n\n", len(symbols))
}
