package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

const marketDayCount = 36

func main() {
	dataDir := os.Getenv("EOD_DATA_DIR")
	if dataDir == "" {
		log.Fatal("Must set environment variable EOD_DATA_DIR")
	}

	// TODO: move to driver, use channels
	dates := PreviousMarketDays(Day(2023, 7, 12), marketDayCount)
	// TODO: AMEX, NYSE
	exchange := "NASDAQ"
	eodData := make([][]*EODData, marketDayCount)

	for i, date := range dates {
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
		// TODO: break out into signals
		if v.DataPoints != marketDayCount {
			continue
		}

		if v.AvgVolume < 500000 || v.AvgClose < 5.0 {
			continue
		}

		// TODO: buy vs sell signals
		if v.MACD.Gap() < 0 || v.MACD.Trend < 0 {
			continue
		}

		if v.LastClose() < v.SMA20.Value() {
			continue
		}

		if v.LastVolume() < v.AvgVolume {
			continue
		}

		symbols = append(symbols, v)
	}

	sort.Slice(symbols, func(i, j int) bool {
		// TODO: better sorting
		a := symbols[i].MACD.Trend + symbols[i].MACD.Gap()
		b := symbols[j].MACD.Trend + symbols[j].MACD.Gap()
		
		return a < b
	})

	for _, v := range symbols {
		fmt.Printf("%s %d (%.2f) (%.2f %.2f) %.2f\n",
			v.Symbol,
			v.AvgVolume,
			v.SMA20.Value(),
			v.MACD.Line,
			v.MACD.Signal.Value,
			v.MACD.Trend)
	}

	fmt.Printf("Found %d symbols\n\n", len(symbols))
}
