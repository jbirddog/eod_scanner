package main

import (
	"fmt"
	"strings"
	"time"
)

type Writer interface {
	WriteHeader(currentDay time.Time)
	WriteSectionHeader(r *ScanResult)
	WriteRecord(a *AnalyzedData, p *Position, risk float64)
	WriteSectionFooter(r *ScanResult)
	WriteFooter()
}

const headerText = "EOD Report for"
const sectionHeaderText = "Strategy"
const sectionFooterText = "symbols"

func headerDateString(currentDay time.Time) string {
	return PreviousMarketDay(currentDay).Format("01/02/2006")
}

const intFmt = "%d"
const floatFmt = "%.2f"
const stringFmt = "%s"

var columns = []struct {
	Lbl string
	Fmt string
}{
	{"Symbol", stringFmt},
	{"Vol X", floatFmt},
	{"Change", floatFmt},
	{"RSI", floatFmt},
	{"Close", floatFmt},
	{"MACD Signal", floatFmt},
	{"MACD Gap", floatFmt},
	{"Position", stringFmt},
	{"Shares", intFmt},
	{"Entry", floatFmt},
	{"Capitol", floatFmt},
	{"Stop Loss", floatFmt},
}

func columnFields() ([]string, []string) {
	lbls := make([]string, len(columns))
	fmts := make([]string, len(columns))

	for i, c := range columns {
		lbls[i] = c.Lbl
		fmts[i] = c.Fmt
	}

	return lbls, fmts
}

//
// Markdown Writer
//

type MarkdownWriter struct {
	_tableHeader string
	_recordFmt   string
}

func NewMarkdownWriter() *MarkdownWriter {
	writer := &MarkdownWriter{}
	lbls, fmts := columnFields()

	writer.setTableHeader(lbls)
	writer.setRecordFmt(fmts)

	return writer
}

func (m *MarkdownWriter) setTableHeader(lbls []string) {
	m._tableHeader = fmt.Sprintf("| %s |\n|%s\n",
		strings.Join(lbls, " | "), strings.Repeat("----|", len(lbls)))
}

func (m *MarkdownWriter) setRecordFmt(fmts []string) {
	m._recordFmt = fmt.Sprintf("| %s |\n", strings.Join(fmts, " | "))
}

func (m *MarkdownWriter) WriteHeader(currentDay time.Time) {
	fmt.Printf("# %s %s\n\n", headerText, headerDateString(currentDay))
}

func (m *MarkdownWriter) WriteSectionHeader(r *ScanResult) {
	fmt.Printf("## %s '%s'\n\n%s", sectionHeaderText, r.Strategy.Name(), m._tableHeader)
}

func (m *MarkdownWriter) WriteRecord(a *AnalyzedData, p *Position, risk float64) {
	fmt.Printf(m._recordFmt,
		a.Symbol,
		a.LastVolumeMultiplier(),
		a.LastChange(),
		a.Indicators.RSI.Value,
		a.LastClose(),
		a.Indicators.MACD.Signal.Value,
		a.Indicators.MACD.Gap(),
		p.Type.String(),
		p.Shares,
		p.Entry,
		p.Capitol,
		p.StopLoss)

}

func (m *MarkdownWriter) WriteSectionFooter(r *ScanResult) {
	fmt.Printf("\n%d %s.\n\n", len(r.Detected), sectionFooterText)
}

func (m *MarkdownWriter) WriteFooter() {
}
