package main

import (
	"flag"
	"log"
)

func main() {
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	config := ConfigFromFile(*configFile)

	dataDir := config.DataDir

	currentDay := Day(2023, 8, 17)
	marketDayCount := 65
	strategies := []Strategy{
		&MonthClimb{},
		&MonthFall{},
		//&MACDFuse{},
	}
	writer := NewMarkdownWriter()

	// TODO: pass in via config
	results, err := Scan(currentDay, marketDayCount, dataDir, strategies)
	if err != nil {
		log.Fatal(err)
	}

	PrintReport(results, currentDay, writer)
}
