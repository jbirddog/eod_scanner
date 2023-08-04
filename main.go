package main

import (
	"log"
	"os"
)

func main() {
	dataDir := os.Getenv("EOD_DATA_DIR")
	if dataDir == "" {
		log.Fatal("Must set environment variable EOD_DATA_DIR")
	}

	currentDay := Day(2023, 8, 4)
	marketDayCount := 52
	strategies := []Strategy{
		&MonthClimb{},
		&MonthFall{},
	}
	writer := &MarkdownWriter{}

	// TODO: pass in via config
	results, err := Scan(currentDay, marketDayCount, dataDir, strategies)
	if err != nil {
		log.Fatal(err)
	}

	PrintReport(results, currentDay, writer)
}
