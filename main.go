package main

import (
	"flag"
	"log"
)

func main() {
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	config, err := ConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s\n", err)
	}

	strategies := []Strategy{
		&MonthClimb{},
		&MonthFall{},
		//&MACDFuse{},
	}

	// TODO: pass in via config
	results, err := Scan(config.CurrentDay, config.MarketDayCount, config.DataDir, strategies)
	if err != nil {
		log.Fatal(err)
	}

	PrintReport(results, config.CurrentDay, config.Writer)
}
