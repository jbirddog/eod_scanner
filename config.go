package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"time"
)

type Config struct {
	DataDir        string    `json:dataDir`
	MarketDayCount int       `json:marketDayCount`
	CurrentDay     time.Time `json:currentDay`
	RiskPerTrade   float64   `json:riskPerTrade`
	StrategyNames  []string  `json:strategyNames`
	WriterName     string    `json:writerName`
	Writer         Writer
	Strategies     []Strategy
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
		return errors.New("Field `dataDir` is empty or missing")
	}

	if c.RiskPerTrade < math.SmallestNonzeroFloat64 {
		return errors.New("Field `riskPerTrade` must be a value > 0.0")
	}

	if len(c.StrategyNames) == 0 {
		return errors.New("Field `strategyNames` is empty or missing")
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
		return errors.New("Invalid writerName")
	}

	c.Strategies = make([]Strategy, len(c.StrategyNames))

	for i, name := range c.StrategyNames {
		switch name {
		case "monthClimb":
			c.Strategies[i] = &MonthClimb{}
		case "monthFall":
			c.Strategies[i] = &MonthFall{}
		default:
			return fmt.Errorf("Invalid strategyName '%s'", name)
		}
	}

	return nil
}
