package main

import (
	"fmt"
	"time"
)

type Writer interface {
	WriteHeader(currentDay time.Time)
	WriteSectionHeader(r *ScanResult)
	WriteRecord(a *AnalyzedData, p *Position, risk float64)
	WriteSectionFooter(r *ScanResult)
	WriteFooter()
}

//
// Markdown Writer
//

type MarkdownWriter struct{}

func (m *MarkdownWriter) WriteHeader(currentDay time.Time) {
	fmt.Printf("# EOD Report for %s\n\n", PreviousMarketDay(currentDay).Format("01/02/2006"))
}

func (m *MarkdownWriter) WriteSectionHeader(r *ScanResult) {
	fmt.Printf(`## Strategy '%s'

| Symbol | RSI | AvgVolume | Close | Change | MACD | MACD Signal | MACD Trend | Shares | Entry | Capitol | Stop Loss |
|----|----|----|----|----|----|----|----|----|----|----|----|
`, r.Strategy.Name)
}

func (m *MarkdownWriter) WriteRecord(a *AnalyzedData, p *Position, risk float64) {
	fmt.Print("RECORD")
	/*
	
			// TODO: string builder or buffered writer?
			fmt.Printf("%s %.2f %.0f (%.2f %.2f) (%.2f %.2f) %.2f | %d @ %.2f ~ %.2f > %.2f\n",
				v.Symbol,
				v.RSI.Value,
				v.AvgVolume,
				v.LastClose(),
				v.LastChange(),
				v.MACD.Line,
				v.MACD.Signal.Value,
				v.MACD.Trend,
				p.Shares,
				p.Entry,
				p.Capitol,
				p.StopLoss)
		}

	*/
}

func (m *MarkdownWriter) WriteSectionFooter(r *ScanResult) {
	fmt.Printf("Strategy '%s' found %d symbols.\n\n",
		r.Strategy.Name,
		len(r.Detected))
}

func (m *MarkdownWriter) WriteFooter() {
}
