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

| Symbol | Vol X | Change | RSI | Volume | Close | MACD Signal | MACD Gap | Shares | Entry | Capitol | Stop Loss |
|----|----|----|----|----|----|----|----|----|----|----|----|
`, r.Strategy.Name())
}

func (m *MarkdownWriter) WriteRecord(a *AnalyzedData, p *Position, risk float64) {
	fmt.Printf("| %s | %.2f | %.2f | %.2f | %.0f | %.2f | %.2f | %.2f | %d | %.2f | %.2f | %.2f |\n",
		a.Symbol,
		a.LastVolumeMultiplier(),
		a.LastChange(),
		a.RSI.Value,
		a.LastVolume(),
		a.LastClose(),
		a.MACD.Signal.Value,
		a.MACD.Gap(),
		p.Shares,
		p.Entry,
		p.Capitol,
		p.StopLoss)

}

func (m *MarkdownWriter) WriteSectionFooter(r *ScanResult) {
	fmt.Printf("\n%d symbols.\n\n", len(r.Detected))
}

func (m *MarkdownWriter) WriteFooter() {
}
