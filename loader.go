package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func EODFilePath(dir string, exchange string, date time.Time) string {
	file_name := fmt.Sprintf("%s_%d%02d%02d.csv", exchange, date.Year(), date.Month(), date.Day())
	return path.Join(dir, file_name)
}

func LoadEODFile(dir string, exchange string, date time.Time) []string {
	path := EODFilePath(dir, exchange, date)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0, 4096)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}
