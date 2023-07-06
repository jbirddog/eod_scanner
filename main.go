package main

import (
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
	dates := PreviousMarketDays(Day(2023, 07, 02), marketDayCount)
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

	c := 0
	d := 0
	e := 0

	for _, v := range analyzedDataBySymbol {
		if v.DataPoints == marketDayCount {
			c += 1

			if v.AvgVolume >= 500000 {
				d += 1

				if v.SMA20 >= v.LastClose {
					e += 1
				}
			}
		}
	}

	log.Fatalf("%d , %d , %d of %d (%f) %d", c, d, e, len(analyzedDataBySymbol), analyzedDataBySymbol["AMD"].SMA20, analyzedDataBySymbol["AMD"].AvgVolume)
}
