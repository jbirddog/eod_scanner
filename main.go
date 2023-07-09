package main

import (
	"fmt"
	"log"
	"os"
)

const marketDayCount = 30

func main() {
	dataDir := os.Getenv("EOD_DATA_DIR")
	if dataDir == "" {
		log.Fatal("Must set environment variable EOD_DATA_DIR")
	}

	// TODO: move to driver, use channels
	dates := PreviousMarketDays(Day(2023, 07, 06), marketDayCount)
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

	for k, v := range analyzedDataBySymbol {
		if v.DataPoints != marketDayCount || v.AvgVolume < 500000 || v.AvgClose < 5.0 {
			continue
		}

		if !v.MACDWasNeg || !v.MACDIsPos {
			continue
		}

		fmt.Printf("%s: (%f %f) %d (%f %f %f)\n",
			k,
			v.SMA20,
			v.SMA20_X.Value,
			v.AvgVolume,
			v.EMA26,
			v.MACDLine,
			v.MACDSignalLine)
	}
}
