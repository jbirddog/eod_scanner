package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	configFile := flag.String("config", "", "Path to config file")
	cpuProfile := flag.String("cpuprofile", "", "Write cpu profile to file")
	memProfile := flag.String("memprofile", "", "Write memory profile to file")
	
	flag.Parse()

	if *cpuProfile != "" {
		file, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatalf("Failed to create cpu profile file: %s\n", err)
		}

		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	config, err := ConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s\n", err)
	}

	results, err := Scan(config.CurrentDay,
		config.MarketDayCount,
		config.DataDir,
		config.Exchange,
		config.Strategies)

	if err != nil {
		log.Fatal(err)
	}

	PrintReport(results,
		config.CurrentDay,
		config.RiskPerTrade,
		config.Writer)
		
	if *memProfile != "" {
		file, err := os.Create(*memProfile)
		if err != nil {
			log.Fatalf("Failed to create memory profile file: %s\n", err)
		}

		pprof.WriteHeapProfile(file)
	}
}
