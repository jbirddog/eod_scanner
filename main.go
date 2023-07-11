package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

const marketDayCount = 35

func main() {
	dataDir := os.Getenv("EOD_DATA_DIR")
	if dataDir == "" {
		log.Fatal("Must set environment variable EOD_DATA_DIR")
	}

	// TODO: move to driver, use channels
	dates := PreviousMarketDays(Day(2023, 7, 10), marketDayCount)
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
		if v.DataPoints != marketDayCount || v.AvgVolume < 500000 || v.AvgClose < 5.0 {
			continue
		}

		if !v.MACDWasNeg || !v.MACDIsPos {
			continue
		}

		symbols = append(symbols, v)
	}

	sort.Slice(symbols, func(i, j int) bool {
		// TODO: better sorting besides just by name
		return symbols[i].Symbol < symbols[j].Symbol
	})

	fmt.Printf("Found %d symbols\n\n", len(symbols))

	for _, v := range symbols {
		fmt.Printf("%s %d (%.2f) (%.2f -- %.2f %.2f)\n",
			v.Symbol,
			v.AvgVolume,
			v.SMA20.Value(),
			v.EMA12.Value,
			v.MACDLine,
			v.MACDSignalLine)
	}
}
