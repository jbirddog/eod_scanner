package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type EODData struct {
	Symbol string
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

const header = "Symbol,Date,Open,High,Low,Close,Volume"

func ParseEODFile(rawData []string) ([]*EODData, error) {
	if len(rawData) == 0 || rawData[0] != header {
		return nil, errors.New("Expected header as first line")
	}

	rawData = rawData[1:]

	if len(rawData) == 0 {
		return nil, errors.New("Expected records to parse")
	}

	data := make([]*EODData, len(rawData))

	for i, raw := range rawData {
		eodData, err := parseRow(raw)

		if err != nil {
			return nil, err
		}

		data[i] = eodData
	}

	return data, nil
}

func parseRow(row string) (*EODData, error) {
	parts := strings.Split(row, ",")

	if len(parts) != 7 {
		return nil, errors.New("Expected record to have 7 fields")
	}

	symbol := parts[0]

	date, err := parseDate(parts[1])
	if err != nil {
		return nil, err
	}

	var prices [4]float64
	if err := parsePrices(parts[2:6], &prices); err != nil {
		return nil, err
	}

	volume, err := strconv.ParseFloat(parts[6], 64)
	if err != nil {
		return nil, err
	}

	data := &EODData{
		Symbol: symbol,
		Date:   date,
		Open:   prices[0],
		High:   prices[1],
		Low:    prices[2],
		Close:  prices[3],
		Volume: volume,
	}

	return data, nil
}

func parseDate(field string) (time.Time, error) {
	date, err := time.Parse("02-Jan-2006", field)
	if err == nil {
		date = Day(date.Year(), date.Month(), date.Day())
	}

	return date, err
}

func parsePrices(fields []string, prices *[4]float64) error {
	if len(fields) != 4 {
		return errors.New("Expected 4 price fields")
	}

	for i, f := range fields {
		price, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return err
		}
		prices[i] = price
	}

	return nil
}
