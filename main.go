package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	data_dir := os.Getenv("EOD_DATA_DIR")
	if data_dir == "" {
		log.Fatal("Must set environment variable EOD_DATA_DIR")
	}
	
	date := Day(2023, 5, 18)
	data := LoadEODFile(data_dir, "NASDAQ", date)
	fmt.Printf("scanned %d...\n", len(data))
}
