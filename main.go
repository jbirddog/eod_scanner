package main

import (
	"log"
	"os"
)

const marketDayCount = 52
const riskPerTrade = 50000.0 * 0.005

func main() {
	dataDir := os.Getenv("EOD_DATA_DIR")
	if dataDir == "" {
		log.Fatal("Must set environment variable EOD_DATA_DIR")
	}

	// TODO: move to driver, use channels
	currentDay := Day(2023, 7, 22)
	dates := PreviousMarketDays(currentDay, marketDayCount)
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

		if v.AvgVolume < 1000000 || v.AvgClose < 5.0 {
			continue
		}

		if v.LastVolume() < v.AvgVolume {
			continue
		}

		// TODO: break out into buy vs sell signals
		if !v.ClosedUp() {
			continue
		}

		if v.MACD.Gap() < 0 {
			continue
		}

		if v.LastClose() < v.SMA20.Value {
			continue
		}

		symbols = append(symbols, v)
	}

	//PrintReport(symbols, currentDay)
}
