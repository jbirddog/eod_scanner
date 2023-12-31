package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jbirddog/marketday"
)

type Writer interface {
	Name() string
	WriteHeader(currentDay time.Time)
	WriteSectionHeader(r *ScanResult)
	WriteRecord(a *AnalyzedData, p *Position, risk float64)
	WriteSectionFooter(r *ScanResult)
	WriteFooter()
}

func WriterNamed(name string) (Writer, error) {
	if name == "markdown" {
		return NewMarkdownWriter(), nil
	}

	return nil, fmt.Errorf("Unknown writer name: '%s'", name)
}

const headerText = "EOD Report for"
const sectionHeaderText = "Strategy"
const sectionFooterText = "symbols"

func headerDateString(currentDay time.Time) string {
	return marketday.PreviousMarketDay(currentDay).Format("01/02/2006")
}

var columns = []struct {
	Lbl string
	Fmt string
}{
	{"Symbol", "%s"},
	{"Vol X", "%.2f"},
	{"%", "%.2f%%"},
	{"Close", "%.2f (%.2f %.2f %.2f)"},
	{"RSI", "%.2f"},
	{"MACD", "%.2f %.2f"},
	{"Position", "%s ~%d, %.2f, ~%.2f"},
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
	tableHeader string
	recordFmt   string
}

func NewMarkdownWriter() *MarkdownWriter {
	writer := &MarkdownWriter{}
	lbls, fmts := columnFields()

	writer.setTableHeader(lbls)
	writer.setRecordFmt(fmts)

	return writer
}

func (m *MarkdownWriter) Name() string {
	return "Markdown"
}

func (m *MarkdownWriter) setTableHeader(lbls []string) {
	m.tableHeader = fmt.Sprintf("| %s |\n|%s\n",
		strings.Join(lbls, " | "), strings.Repeat("----|", len(lbls)))
}

func (m *MarkdownWriter) setRecordFmt(fmts []string) {
	m.recordFmt = fmt.Sprintf("| %s |\n", strings.Join(fmts, " | "))
}

func (m *MarkdownWriter) WriteHeader(currentDay time.Time) {
	fmt.Printf("# %s %s\n\n", headerText, headerDateString(currentDay))
}

func (m *MarkdownWriter) WriteSectionHeader(r *ScanResult) {
	fmt.Printf("## %s '%s'\n\n%s", sectionHeaderText, r.Strategy.Name(), m.tableHeader)
}

func (m *MarkdownWriter) WriteRecord(a *AnalyzedData, p *Position, risk float64) {
	i := a.Indicators
	lastClose := a.LastClose()

	fmt.Printf(m.recordFmt,
		a.Symbol,
		a.LastVolumeMultiplier(),
		a.LastChange(),
		lastClose,
		i.BB.Upper,
		i.BB.Middle,
		i.BB.Lower,
		//percentage(lastClose, i.EMA8.Value),
		//percentage(lastClose, i.SMA20.Value),
		i.RSI.Value,
		i.MACD.Line,
		i.MACD.Signal.Value,
		p.Type.String(),
		p.Shares,
		p.StopLoss,
		p.Capital)
}

func (m *MarkdownWriter) WriteSectionFooter(r *ScanResult) {
	fmt.Printf("\n%d %s.\n\n", len(r.Detected), sectionFooterText)
}

func (m *MarkdownWriter) WriteFooter() {
}
