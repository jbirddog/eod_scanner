package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DataDir string `json:dataDir`
}

func ConfigFromFile(path string) *Config {
	if path == "" {
		log.Fatal("No config file specified")
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to open config file: %s\n", err)
	}

	var config *Config
	err = json.Unmarshal(contents, &config)
	if err != nil {
		log.Fatalf("Invalid config file: %s\n", err)
	}
	
	if config.DataDir == "" {
		log.Fatal("Field `dataDir` is missing from config")
	}

	return config
}
