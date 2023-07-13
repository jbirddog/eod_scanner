package main

const (
	PositionTypeLong = iota
	PositionTypeShort
)

type Position struct {
	Type     int
	Shares   int
	Entry    float64
	Capitol  float64
	StopLoss float64
}

func PositionFromAnalyzedData(data *AnalyzedData, risk float64) *Position {
	entry := data.LastClose()
	stopLoss := data.SMA20.Value()
	// TODO: revist when adding support for shorts
	shares := int(risk / (entry - stopLoss))
	capitol := float64(shares) * entry

	return &Position{
		Shares:   shares,
		Entry:    entry,
		Capitol:  capitol,
		StopLoss: stopLoss,
	}
}
