package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Config struct {
	DataDir        string    `json:dataDir`
	MarketDayCount int       `json:marketDayCount`
	CurrentDay     time.Time `json:currentDay`
	Writer         Writer
	WriterName     string `json:writerName`
}

func ConfigFromFile(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("No config file specified")
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config *Config
	if err = json.Unmarshal(contents, &config); err != nil {
		return nil, err
	}

	if err = config.validate(); err != nil {
		return nil, err
	}

	config.setDefaultValues()

	if err = config.instantiateObjects(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.DataDir == "" {
		return errors.New("Field `dataDir` is missing from config")
	}

	return nil
}

func (c *Config) setDefaultValues() {
	if c.MarketDayCount == 0 {
		c.MarketDayCount = 65
	}

	if c.CurrentDay.IsZero() {
		now := time.Now()
		c.CurrentDay = Day(now.Year(), now.Month(), now.Day())
	}

	if c.WriterName == "" {
		c.WriterName = "markdown"
	}
}

func (c *Config) instantiateObjects() error {
	switch c.WriterName {
	case "markdown":
		c.Writer = NewMarkdownWriter()
	default:
		return errors.New("Invalid writerName in config")
	}
	
	return nil
}
