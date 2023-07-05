package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	dataDir := os.Getenv("EOD_DATA_DIR")
	if dataDir == "" {
		log.Fatal("Must set environment variable EOD_DATA_DIR")
	}

	date := Day(2023, 5, 18)

	rawData, err := LoadEODFile(dataDir, "NASDAQ", date)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("scanned %d...\n", len(rawData))
}
